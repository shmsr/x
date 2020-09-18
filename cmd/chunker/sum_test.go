package main

import (
	"fmt"
	"reflect"
	"testing"
)

func plainSum(slice []int64) int64 {
	var sum int64 = 0
	for _, num := range slice {
		sum += num
	}
	return sum
}

func BenchmarkAdders(b *testing.B) {
	benchmarks := map[string]struct {
		input  []int64
		chunks int
	}{
		"small":  {input: sliceGenerator(10), chunks: 5},
		"medium": {input: sliceGenerator(10045), chunks: 10},
		"large":  {input: sliceGenerator(10000010), chunks: 15},
	}
	for name, bm := range benchmarks {
		b.Run(fmt.Sprintf("BenchmarkChunker: %s", name), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = chunker(bm.input, bm.chunks)
			}
		})
		b.Run(fmt.Sprintf("BenchmarkPlainSum: %s", name), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				_ = plainSum(bm.input)
			}
		})
	}
}

func TestChunker(t *testing.T) {
	tests := map[string]struct {
		input  []int64
		chunks int
	}{
		"small":  {input: sliceGenerator(10), chunks: 5},
		"medium": {input: sliceGenerator(1004), chunks: 10},
		"large":  {input: sliceGenerator(1000001), chunks: 15},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got := chunker(tc.input, tc.chunks)
			if s := plainSum(tc.input); !reflect.DeepEqual(s, got) {
				t.Fatalf("expected: %v, got: %v", s, got)
			}
		})
	}
}
