package main

import (
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/olekukonko/tablewriter"
)

type BenchmarkResult struct {
	CacheName string
	Duration  time.Duration
	Hits      int64
	Misses    int64
	Bytes     int64
}

func (br *BenchmarkResult) hitRate() float64 {
	return float64(br.Hits) / float64(br.Hits+br.Misses) * 100
}

type Benchmark struct {
	ItemSize            int
	CacheSizeMultiplier float64
	ZipfAlpha           float64
	Concurrency         int
	Results             []*BenchmarkResult
}

func (b *Benchmark) AddResult(r *BenchmarkResult) {
	b.Results = append(b.Results, r)
}

func (b *Benchmark) WriteToConsole() {
	b.sortResults()

	workloads := b.ItemSize * workloadMultiplier

	fmt.Printf("results:\n")
	fmt.Printf("itemSize=%d, workloads=%d, cacheSize=%.2f%%, zipf's alpha=%.2f, concurrency=%d\n\n",
		b.ItemSize,
		workloads,
		b.CacheSizeMultiplier*100,
		b.ZipfAlpha,
		b.Concurrency)

	headers := []string{"Cache", "HitRate", "Memory", "QPS", "Hits", "Misses"}
	table := tablewriter.NewWriter(os.Stdout)
	for _, ret := range b.Results {
		qps := float64(ret.Hits+ret.Misses) / float64(ret.Duration.Milliseconds())
		qps = qps * 1000

		table.Append([]string{
			ret.CacheName,
			fmt.Sprintf("%.2f%%", ret.hitRate()),
			fmt.Sprintf("%.2fMiB", float64(ret.Bytes)/1000/1000),
			fmt.Sprintf("%.f", qps),
			fmt.Sprintf("%d", ret.Hits),
			fmt.Sprintf("%d", ret.Misses),
		})
	}
	table.SetHeader(headers)
	table.SetBorder(false)
	table.Render()

	fmt.Printf("\n\n")
}

func (b *Benchmark) Clean() {
	b.Results = []*BenchmarkResult{}
}

func (b *Benchmark) sortResults() {
	// sort by hit rate
	sort.Slice(b.Results, func(i, j int) bool {
		return b.Results[i].hitRate() > b.Results[j].hitRate()
	})
}
