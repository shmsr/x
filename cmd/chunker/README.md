# chunker
chunker breaks the slice in chunks and each chunks is processed by goroutine which is assigned to each chunk. For big slices, I have seen 2-3x improvement than the normal method where a single goroutine adds up elements for the same slice.

## Install
* Install in `GOBIN` or `~/go/bin`:
```
go get github.com/shmsr/x/cmd/chunker
```
* Install manually:
```
go build
```

## Example
Generates a slice with random elements (used a random number generator). By default 100000000 element are there in slice which is broken into 10 chunks and 10 goroutines are used to process them.
```sh
chunkAdder 
```
