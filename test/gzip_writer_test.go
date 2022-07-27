package main

import (
	"compress/gzip"
	"io/ioutil"
	"testing"
)

func BenchmarkWriters(b *testing.B) {
	w := ioutil.Discard
	d := make([]byte, 1024*1024)
	for n := 0; n < b.N; n++ {
		z, _ := gzip.NewWriterLevel(w, gzip.BestCompression)
		z.Write(d)
		z.Close()
	}
}
