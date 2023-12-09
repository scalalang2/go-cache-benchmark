package main

import (
	"fmt"
	"go-cache-benchmark/zipf"
	"math/rand"
)

type ZipfGenerator struct {
	gen *zipf.ZipfGenerator
}

func NewZipfGenerator(size uint64, theta float64) *ZipfGenerator {
	src := rand.NewSource(19931203)
	r := rand.New(src)
	gen, err := zipf.NewZipfGenerator(r, 0, size, theta, false)

	if err != nil {
		panic(fmt.Errorf("could not create zipf generator: %v", err))
	}
	return &ZipfGenerator{
		gen: gen,
	}
}

func (z *ZipfGenerator) Name() string {
	return "zipf"
}

func (z *ZipfGenerator) Next() string {
	return fmt.Sprintf("%d", z.gen.Uint64())
}
