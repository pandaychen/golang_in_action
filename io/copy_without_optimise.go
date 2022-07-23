package main

import (
	"io"
	"os"
)

type OnlyWriter struct {
	io.Writer
}

func CopyFile(src, dst string) (int64, error) {
	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.CopyBuffer(OnlyWriter{destination}, source, make([]byte, 1024*1024))
	return nBytes, err
}

func main() {
	CopyFile("srcfile", "dstfile")
}
