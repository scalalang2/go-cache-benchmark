package main

import (
	"go-cache-benchmark/cache"
	"runtime"
	"time"
)

type NewCacheFunc func(size int) cache.Cache

func main() {
	zipfAlphas := []float64{0.99}
	cacheSizes := []int{1e3, 1e4, 1e5}
	multiplier := []int{10, 100, 1000}
	caches := []NewCacheFunc{
		cache.NewLRU,
		cache.NewS3FIFO,
	}

	for _, sz := range cacheSizes {
		for _, mul := range multiplier {
			cacheSize := sz
			numberOfKeys := sz * mul
			for _, alpha := range zipfAlphas {
				runBenchmark(cacheSize, uint64(numberOfKeys), alpha, caches)
			}
		}
	}
}

func runBenchmark(cacheSize int, numberOfKeys uint64, zipfAlpha float64, caches []NewCacheFunc) {
	b := &Benchmark{
		CacheSize: cacheSize,
		NumKey:    numberOfKeys,
		ZipfAlpha: zipfAlpha,
		Results:   make([]*BenchmarkResult, 0),
	}

	for _, newCache := range caches {
		b.Results = append(b.Results, run(newCache, cacheSize, numberOfKeys, zipfAlpha))
	}

	b.WriteToConsole()
}

func run(newCache NewCacheFunc, cacheSize int, numberOfKeys uint64, zipfAlpha float64) *BenchmarkResult {
	gen := NewZipfGenerator(numberOfKeys, zipfAlpha)

	alloc1 := memAlloc()
	c := newCache(cacheSize)
	defer c.Close()

	start := time.Now()
	bench := func(c cache.Cache, gen *ZipfGenerator) (hits, misses int64) {
		for i := 0; i < 1e6; i++ {
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
