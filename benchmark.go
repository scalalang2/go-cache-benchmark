package main

import (
	"sort"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
)

type BenchmarkWriter struct {
	chart *charts.Bar
}

func NewBenchmarkWriter(title string, xAxis []string) *BenchmarkWriter {
	c := charts.NewBar()
	c.SetGlobalOptions(charts.WithTitleOpts(opts.Title{
		Title: title,
	}))
	c = c.SetXAxis(xAxis)

	return &BenchmarkWriter{
		chart: c,
	}
}

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
	CacheSizeMultiplier float64
	Concurrency         int
	Results             []*BenchmarkResult
}

func (b *Benchmark) AddResult(r *BenchmarkResult) {
	b.Results = append(b.Results, r)
}

func (b *Benchmark) WriteToConsole() {
	b.sortResults()
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
