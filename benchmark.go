package main

import (
	"fmt"
	"os"
	"time"

	"github.com/go-echarts/go-echarts/v2/charts"
	"github.com/go-echarts/go-echarts/v2/opts"
	"github.com/samber/lo"
)

type BenchmarkWriter struct {
	chart *charts.Bar
	xAxis []string
}

type BenchmarkResult struct {
	Latencies []time.Duration
	Hits      int64
	Misses    int64
	Bytes     int64
}

func (br *BenchmarkResult) hitRate() float64 {
	return float64(br.Hits) / float64(br.Hits+br.Misses) * 100
}

func NewBenchmarkWriter(title string, subtitle string, xAxis []string) *BenchmarkWriter {
	legendOption := charts.WithLegendOpts(opts.Legend{
		Show:   true,
		Bottom: "50%",
		Orient: "vertical",
		X:      "left",
		Y:      "left",
	})

	titleOption := charts.WithTitleOpts(opts.Title{
		Title:    title,
		Subtitle: subtitle,
	})

	c := charts.NewBar()
	c.SetGlobalOptions(legendOption, titleOption)
	c = c.SetXAxis(xAxis)

	return &BenchmarkWriter{
		chart: c,
		xAxis: xAxis,
	}
}

func (b *BenchmarkWriter) AppendResult(name string, result []*BenchmarkResult) {
	if len(result) != len(b.xAxis) {
		panic("the result doesn't match the number of X axis.")
	}

	data := lo.Map(result, func(item *BenchmarkResult, index int) opts.BarData {
		return opts.BarData{
			Value: item.hitRate(),
		}
	})
	b.chart = b.chart.AddSeries(name, data)
}

func (b *BenchmarkWriter) Write(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("unable to create a new file %s, %v", filename, err)
	}

	return b.chart.Render(f)
}
