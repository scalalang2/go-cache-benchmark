package main

import (
	"runtime"
	"sync"
	"time"

	"go-cache-benchmark/cache"
)

const workloadMultiplier = 15

type NewCacheFunc func(size int) cache.Cache

func main() {
	zipfAlphas := []float64{0.99}
	items := []int{1e5 * 5}
	concurrencies := []int{1, 2, 4, 8, 16}
	cacheSizeMultiplier := []float64{0.0006, 0.001, 0.01, 0.1}
	caches := []NewCacheFunc{
		cache.NewLRU,
		cache.NewSieve,
		cache.NewS3FIFO,
		cache.NewTwoQueue,
		cache.NewLRUGroupCache,
		cache.NewTinyLFU,
		cache.NewSLRU,
		cache.NewS4LRU,
		cache.NewClock,
		cache.NewOtter,
	}

	for _, itemSize := range items {
		for _, multiplier := range cacheSizeMultiplier {
			for _, curr := range concurrencies {
				for _, alpha := range zipfAlphas {
					runBenchmark(itemSize, multiplier, alpha, caches, curr)
				}
			}
		}
	}
}

func runBenchmark(itemSize int, cacheMultiplier float64, zipfAlpha float64, caches []NewCacheFunc, concurrency int) {
	b := &Benchmark{
		ItemSize:            itemSize,
		CacheSizeMultiplier: cacheMultiplier,
		ZipfAlpha:           zipfAlpha,
		Concurrency:         concurrency,
		Results:             make([]*BenchmarkResult, 0),
	}

	for _, newCache := range caches {
		b.Results = append(b.Results, run(newCache, itemSize, cacheMultiplier, zipfAlpha, concurrency))
	}

	b.WriteToConsole()
}

func run(newCache NewCacheFunc, itemSize int, cacheSizeMultiplier float64, zipfAlpha float64, concurrency int) *BenchmarkResult {
	gen := NewZipfGenerator(uint64(itemSize), zipfAlpha)

	total := itemSize * workloadMultiplier
	each := total / concurrency

	alloc1 := memAlloc()

	// create keys in advance to not taint the QPS
	keys := make([][]string, concurrency)
	for i := 0; i < concurrency; i++ {
		keys[i] = make([]string, 0, each)
		for j := 0; j < each; j++ {
			keys[i] = append(keys[i], gen.Next())
		}
	}

	cacheSize := int(float64(itemSize) * cacheSizeMultiplier)
	c := newCache(cacheSize)
	defer c.Close()

	start := time.Now()
	bench := func(c cache.Cache, gen *ZipfGenerator) (int64, int64) {
		var wg sync.WaitGroup
		hits := make([]int64, concurrency)
		misses := make([]int64, concurrency)

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func(k int) {
				for j := 0; j < each; j++ {
					key := keys[k][j]
					if c.Get(key) {
						hits[k]++
					} else {
						misses[k]++
						c.Set(key)
					}
				}
				wg.Done()
			}(i)
		}

		wg.Wait()
		var totalHits, totalMisses int64
		for i := 0; i < concurrency; i++ {
			totalHits += hits[i]
			totalMisses += misses[i]
		}
		return totalHits, totalMisses
	}

	hits, misses := bench(c, gen)
	elapsed := time.Since(start)
	keys = nil
	alloc2 := memAlloc()

	return &BenchmarkResult{
		CacheName: c.Name(),
		Duration:  elapsed,
		Hits:      hits,
		Misses:    misses,
		Bytes:     int64(alloc2) - int64(alloc1),
	}
}

func memAlloc() uint64 {
	runtime.GC()
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	return m.Alloc
}
