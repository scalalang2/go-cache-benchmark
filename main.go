package main

import (
	"go-cache-benchmark/cache"
	"runtime"
	"time"
)

type NewCacheFunc func(size int) cache.Cache

func main() {
	zipfAlphas := []float64{0.99}
	items := []int{1e5 * 5}
	cacheSizeMultiplier := []float64{0.001, 0.01, 0.1}
	caches := []NewCacheFunc{
		cache.NewLRU,
		cache.NewS3FIFO,
		cache.NewTwoQueue,
		cache.NewLRUGroupCache,
		cache.NewTinyLFU,
		cache.NewSLRU,
		cache.NewS4LRU,
		cache.NewClock,
	}

	for _, itemSize := range items {
		for _, multiplier := range cacheSizeMultiplier {
			for _, alpha := range zipfAlphas {
				runBenchmark(itemSize, multiplier, alpha, caches)
			}
		}
	}
}

func runBenchmark(itemSize int, cacheMultiplier float64, zipfAlpha float64, caches []NewCacheFunc) {
	b := &Benchmark{
		ItemSize:            itemSize,
		CacheSizeMultiplier: cacheMultiplier,
		ZipfAlpha:           zipfAlpha,
		Results:             make([]*BenchmarkResult, 0),
	}

	for _, newCache := range caches {
		b.Results = append(b.Results, run(newCache, itemSize, cacheMultiplier, zipfAlpha))
	}

	b.WriteToConsole()
}

func run(newCache NewCacheFunc, itemSize int, cacheSizeMultiplier float64, zipfAlpha float64) *BenchmarkResult {
	gen := NewZipfGenerator(uint64(itemSize), zipfAlpha)

	alloc1 := memAlloc()
	cacheSize := int(float64(itemSize) * cacheSizeMultiplier)
	c := newCache(cacheSize)
	defer c.Close()

	start := time.Now()
	bench := func(c cache.Cache, gen *ZipfGenerator) (hits, misses int64) {
		for i := 0; i < itemSize*5; i++ {
			key := gen.Next()
			if c.Get(key) {
				hits++
			} else {
				misses++
				c.Set(key)
			}
		}
		return
	}

	hits, misses := bench(c, gen)
	elapsed := time.Since(start)
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
