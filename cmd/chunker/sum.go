package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	elements = 100000000
	chunks   = 10
)

// chunkSum sums up the chunked slice
func chunkSum(s []int64, c chan<- int64) {
	var sum int64
	for _, v := range s {
		sum += v
	}
	c <- sum
}

// chunker chunks the slice and spawns goroutines
// for each of the chunk to processed by chunkSum
func chunker(slice []int64, chunks int) int64 {
	var length = len(slice)
	var collector = make(chan int64, chunks)
	// If the length of slice if lesser than chunks, then don't chunk
	if length < chunks {
		go chunkSum(slice, collector)
		return <-collector
	}

	// Chunk and add until end doesn't reach boundary
	var buckets = length / chunks
	var begin = 0
	var end = buckets
	for end <= length {
		go chunkSum(slice[begin:end], collector)
		begin = end
		end += buckets
	}

	// If some elements in the slice are still left
	var routine int
	if length%chunks != 0 {
		routine++
		go chunkSum(slice[begin:], collector)
	}

	// Recieve from every goroutine
	var sum int64 = 0
	for i := 0; i < chunks+routine; i++ {
		sum += <-collector
	}

	// Close the channel
	close(collector)
	return sum
}

// sliceGenerator generates a slice of size: size
func sliceGenerator(size int) []int64 {
	slice := make([]int64, size)
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		slice[i] = rand.Int63n(999) - rand.Int63n(999)
	}
	return slice
}

func main() {
	stream := sliceGenerator(elements)
	fmt.Printf("Result: %d\n", chunker(stream, chunks))
}
