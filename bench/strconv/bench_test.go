package main

import (
	"fmt"
	"strconv"
	"testing"
)

func Benchmark_Atof64_Decimal(b *testing.B) {
	for b.Loop() {
		f, _ := strconv.ParseFloat("33909", 64)
		sinkFloat += f
	}
}

func Benchmark_Atof64_Float(b *testing.B) {
	for b.Loop() {
		f, _ := strconv.ParseFloat("339.7784", 64)
		sinkFloat += f
	}
}

func Benchmark_Atof64_Exp(b *testing.B) {
	for b.Loop() {
		f, _ := strconv.ParseFloat("-5.09e75", 64)
		sinkFloat += f
	}
}

func Benchmark_Atof64_Big(b *testing.B) {
	for b.Loop() {
		f, _ := strconv.ParseFloat("1844674407370955", 64)
		sinkFloat += f
	}
}

func Benchmark_ParseInt_7bit(b *testing.B) {
	s := fmt.Sprintf("%d", 1<<7-1)
	for b.Loop() {
		n, _ := strconv.ParseInt(s, 10, 64)
		sinkInt += int(n)
	}
}

func Benchmark_ParseInt_26bit(b *testing.B) {
	s := fmt.Sprintf("%d", 1<<26-1)
	for b.Loop() {
		n, _ := strconv.ParseInt(s, 10, 64)
		sinkInt += int(n)
	}
}

func Benchmark_ParseInt_31bit(b *testing.B) {
	s := fmt.Sprintf("%d", 1<<31-1)
	for b.Loop() {
		n, _ := strconv.ParseInt(s, 10, 64)
		sinkInt += int(n)
	}
}

func Benchmark_ParseInt_56bit(b *testing.B) {
	s := fmt.Sprintf("%d", 1<<56-1)
	for b.Loop() {
		n, _ := strconv.ParseInt(s, 10, 64)
		sinkInt += int(n)
	}
}

func Benchmark_ParseInt_62bit(b *testing.B) {
	s := fmt.Sprintf("%d", 1<<62-1)
	for b.Loop() {
		n, _ := strconv.ParseInt(s, 10, 64)
		sinkInt += int(n)
	}
}

func Benchmark_FormatFloat_Decimal(b *testing.B) {
	for b.Loop() {
		s := strconv.FormatFloat(33909, 'f', -1, 64)
		sinkInt += len(s)
	}
}

func Benchmark_FormatFloat_Float(b *testing.B) {
	for b.Loop() {
		s := strconv.FormatFloat(339.7784, 'f', -1, 64)
		sinkInt += len(s)
	}
}

func Benchmark_FormatFloat_Exp(b *testing.B) {
	for b.Loop() {
		s := strconv.FormatFloat(-5.09e75, 'e', -1, 64)
		sinkInt += len(s)
	}
}

func Benchmark_FormatFloat_Big(b *testing.B) {
	for b.Loop() {
		s := strconv.FormatFloat(1844674407370955, 'f', -1, 64)
		sinkInt += len(s)
	}
}

func Benchmark_FormatInt_7bit(b *testing.B) {
	n := 1<<7 - 1
	for b.Loop() {
		s := strconv.FormatInt(int64(n), 10)
		sinkInt += len(s)
	}
}

func Benchmark_FormatInt_26bit(b *testing.B) {
	n := 1<<26 - 1
	for b.Loop() {
		s := strconv.FormatInt(int64(n), 10)
		sinkInt += len(s)
	}
}

func Benchmark_FormatInt_31bit(b *testing.B) {
	n := 1<<31 - 1
	for b.Loop() {
		s := strconv.FormatInt(int64(n), 10)
		sinkInt += len(s)
	}
}

func Benchmark_FormatInt_56bit(b *testing.B) {
	n := 1<<56 - 1
	for b.Loop() {
		s := strconv.FormatInt(int64(n), 10)
		sinkInt += len(s)
	}
}

func Benchmark_FormatInt_62bit(b *testing.B) {
	n := 1<<62 - 1
	for b.Loop() {
		s := strconv.FormatInt(int64(n), 10)
		sinkInt += len(s)
	}
}
