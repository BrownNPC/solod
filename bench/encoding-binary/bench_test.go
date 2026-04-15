package main

import (
	"encoding/binary"
	"testing"
)

var putbuf = []byte{0, 0, 0, 0, 0, 0, 0, 0}

func Benchmark_BE_PutUint64(b *testing.B) {
	b.SetBytes(8)
	for i := range b.N {
		binary.BigEndian.PutUint64(putbuf[:8], uint64(i))
	}
}

func Benchmark_BE_AppendUint64(b *testing.B) {
	b.SetBytes(8)
	for i := range b.N {
		putbuf = binary.BigEndian.AppendUint64(putbuf[:0], uint64(i))
	}
}

func Benchmark_LE_PutUint64(b *testing.B) {
	b.SetBytes(8)
	for i := range b.N {
		binary.LittleEndian.PutUint64(putbuf[:8], uint64(i))
	}
}

func Benchmark_LE_AppendUint64(b *testing.B) {
	b.SetBytes(8)
	for i := range b.N {
		putbuf = binary.LittleEndian.AppendUint64(putbuf[:0], uint64(i))
	}
}
