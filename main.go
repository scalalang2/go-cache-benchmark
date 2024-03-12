package main

import (
	"fmt"
	"log/slog"
	"runtime"
	"sync"
	"time"

	"go-cache-benchmark/cache"
)

// itemSize is the number of items used in the benchmark
const itemSize = 500000

// zipfAlpha is the skewness of the zipfian distribution
const zipfAlpha = 0.99

type NewCacheFunc func(size int) cache.Cache

type CacheWithName struct {
	Name        string
	Initializer NewCacheFunc
}

func main() {
	caches := []*CacheWithName{
		{"sieve", cache.NewSieve},
		{"s3-fifo", cache.NewS3FIFO},
		{"otter", cache.NewOtter},
		{"lru-hashicorp", cache.NewLRU},
		{"two-queue", cache.NewTwoQueue},
		{"lru-groupcache", cache.NewLRUGroupCache},
		{"tinylfu", cache.NewTinyLFU},
		{"slru", cache.NewSLRU},
		{"s4-lru", cache.NewS4LRU},
		{"clock", cache.NewClock},
		//{"freelru-synced", cache.NewFreeLRUSynced},
		//{"freelru-sharded", cache.NewFreeLRUSharded},
	}
	runBenchmark(caches)
}

func runBenchmark(caches []*CacheWithName) {
	//concurrencies := []int{1, 2, 4, 8, 16}
	//cacheSizeMultiplier := []float64{0.001, 0.01, 0.1}
	concurrencies := []int{1}
	cacheSizeMultiplier := []float64{0.001}

	xAxis := make([]string, 0)
	for i := 0; i < len(cacheSizeMultiplier); i++ {
		for j := 0; j < len(concurrencies); j++ {
			xAxis = append(xAxis, fmt.Sprintf("%d/%d", int(float64(itemSize)*cacheSizeMultiplier[i]), concurrencies[j]))
		}
	}

	writer := NewBenchmarkWriter(
		"A comparison cache algorithms under Zipfian distribution.",
		"In the A/B format, A is the size of the cache and B is the number of concurrent operations.",
		xAxis)

	slog.Info("start to run the benchmark")

	for _, c := range caches {
		results := make([]*BenchmarkResult, 0)

		for i := 0; i < len(cacheSizeMultiplier); i++ {
			for j := 0; j < len(concurrencies); j++ {
				slog.Info("a test was evaluated", "cache", c.Name, "cacheSize", itemSize*cacheSizeMultiplier[i], "concurrency", concurrencies[j])
				results = append(results, run(c.Initializer, cacheSizeMultiplier[i], concurrencies[j]))
			}
		}

		writer.AppendResult(c.Name, results)
	}

	err := writer.Write("test.html")
	if err != nil {
		panic(fmt.Errorf("failed to create a report: %v", err))
	}
}

func run(newCache NewCacheFunc, cacheSizeMultiplier float64, concurrency int) *BenchmarkResult {
	workloads := itemSize * 20
	gen := NewZipfGenerator(uint64(itemSize), zipfAlpha)
	each := workloads / concurrency

	// get memory allocation before starting the testing
	alloc1 := memAlloc()

	// record a latency for a certain number of requests
	latencies := make([]time.Duration, 0)
	cnt := 0

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

	bench := func(c cache.Cache, gen *ZipfGenerator) (int64, int64) {
		var wg sync.WaitGroup
		hits := make([]int64, concurrency)
		misses := make([]int64, concurrency)

		start := time.Now()

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

					// record latencies per 1000 requests
					cnt++
					if cnt%1000 == 0 {
						cnt = 0
						elapsed := time.Since(start)
						start = time.Now()
						latencies = append(latencies, elapsed)
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
	keys = nil
	alloc2 := memAlloc()

	return &BenchmarkResult{
		Latencies: latencies,
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
