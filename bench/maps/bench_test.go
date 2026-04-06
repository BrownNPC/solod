package main

import (
	"fmt"
	"testing"
)

func init() {
	for i := range nKeys {
		strKeys = append(strKeys, fmt.Sprintf("key-%d", i))
	}
}

func Benchmark_IntSet(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[int]int)
		for i := range nKeys {
			m[i] = i
		}
	}
}

func Benchmark_IntGet(b *testing.B) {
	b.ReportAllocs()
	m := make(map[int]int, nKeys)
	for i := range nKeys {
		m[i] = i
	}
	for b.Loop() {
		for i := range nKeys {
			sinkInt = m[i]
		}
	}
}

func Benchmark_IntDel(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[int]int, nKeys)
		for i := range nKeys {
			m[i] = i
		}
		for i := range nKeys {
			delete(m, i)
		}
	}
}

func Benchmark_StrSet(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[string]int)
		for i := range nKeys {
			m[strKeys[i]] = i
		}
	}
}

func Benchmark_StrGet(b *testing.B) {
	b.ReportAllocs()
	m := make(map[string]int, nKeys)
	for i := range nKeys {
		m[strKeys[i]] = i
	}
	for b.Loop() {
		for i := range nKeys {
			sinkInt = m[strKeys[i]]
		}
	}
}

func Benchmark_StrDel(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		m := make(map[string]int, nKeys)
		for i := range nKeys {
			m[strKeys[i]] = i
		}
		for i := range nKeys {
			delete(m, strKeys[i])
		}
	}
}
