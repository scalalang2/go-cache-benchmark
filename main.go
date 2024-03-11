package main

import (
	"runtime"
	"sync"
	"time"

	"go-cache-benchmark/cache"
)

// totalItems is the number of items used in the benchmark
const totalItems = 500000

// zipfAlpha is the skewness of the zipfian distribution
const zipfAlpha = 0.99

type NewCacheFunc func(size int) cache.Cache

func main() {
	concurrencies := []int{1, 2, 4, 8, 16}
	cacheSizeMultiplier := []float64{0.001, 0.01, 0.1}
	caches := []NewCacheFunc{
		cache.NewSieve,
		cache.NewS3FIFO,
		cache.NewOtter,
		cache.NewLRU,
		cache.NewTwoQueue,
		cache.NewLRUGroupCache,
		cache.NewTinyLFU,
		cache.NewSLRU,
		cache.NewS4LRU,
		cache.NewClock,
		cache.NewFreeLRUSynced,
		cache.NewFreeLRUSharded,
	}

	for _, multiplier := range cacheSizeMultiplier {
		for _, curr := range concurrencies {
			runBenchmark(multiplier, caches, curr)
		}
	}
}

func runBenchmark(cacheMultiplier float64, caches []NewCacheFunc, concurrency int) {
	b := &Benchmark{
		CacheSizeMultiplier: cacheMultiplier,
		Concurrency:         concurrency,
		Results:             make([]*BenchmarkResult, 0),
	}

	for _, newCache := range caches {
		b.Results = append(b.Results, run(newCache, cacheMultiplier, concurrency))
	}

	b.WriteToConsole()
}

func run(newCache NewCacheFunc, cacheSizeMultiplier float64, concurrency int) *BenchmarkResult {
	gen := NewZipfGenerator(uint64(totalItems), zipfAlpha)

	total := totalItems
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

	cacheSize := int(float64(totalItems) * cacheSizeMultiplier)
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
