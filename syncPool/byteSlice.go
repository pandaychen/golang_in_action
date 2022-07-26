package main

import (
	"sync"
)

var (
	maxSize           = 20000 // 20k
	defaultBufferPool = newBufferPool(maxSize)
)

type bufferPool struct {
	bp      *sync.Pool
	maxSize int
}

func newBufferPool(maxSize int) *bufferPool {
	return &bufferPool{
		bp: &sync.Pool{
			New: func() interface{} {
				b := make([]byte, 2048)
				return b
			},
		},
		maxSize: maxSize,
	}
}

func (bpool *bufferPool) get() []byte {
	return bpool.bp.Get().([]byte)
}

func (bpool *bufferPool) put(buf []byte) {
	if maxSize != 0 && cap(buf) > maxSize {
		return
	}

	// reset length to 0
	bpool.bp.Put(buf[:0])
}

func main() {
	array := defaultBufferPool.get()
	array[0] = 1
	defaultBufferPool.put(array)
}
