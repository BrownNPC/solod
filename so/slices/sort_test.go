// Copyright 2023 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package slices_test

import (
	"cmp"
	"fmt"
	"math/rand"
	"testing"

	"solod.dev/so/c"
	"solod.dev/so/math"
	. "solod.dev/so/slices"
)

var ints = [...]int{74, 59, 238, -784, 9845, 959, 905, 0, 0, 42, 7586, -5467984, 7586}
var float64s = [...]float64{74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3, math.Inf(-1), 9845.768, -959.7485, 905, 7.8, 7.8, 74.3, 59.0, math.Inf(1), 238.2, -784.0, 2.3}
var strs = [...]string{"", "Hello", "foo", "bar", "foo", "f00", "%*&^*&^&", "***"}

func TestSortIntSlice(t *testing.T) {
	data := Clone(nil, ints[:])
	Sort(data)
	if !IsSorted(data) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestSortFuncIntSlice(t *testing.T) {
	data := Clone(nil, ints[:])
	SortFunc(data, func(a, b any) int {
		va := *c.PtrAs[int](a)
		vb := *c.PtrAs[int](b)
		return va - vb
	})
	if !IsSorted(data) {
		t.Errorf("sorted %v", ints)
		t.Errorf("   got %v", data)
	}
}

func TestSortFloat64Slice(t *testing.T) {
	data := Clone(nil, float64s[:])
	Sort(data)
	if !IsSorted(data) {
		t.Errorf("sorted %v", float64s)
		t.Errorf("   got %v", data)
	}
}

func TestSortStringSlice(t *testing.T) {
	data := Clone(nil, strs[:])
	Sort(data)
	if !IsSorted(data) {
		t.Errorf("sorted %v", strs)
		t.Errorf("   got %v", data)
	}
}

func TestSortLarge_Random(t *testing.T) {
	n := 10000
	data := make([]int, n)
	for i := 0; i < len(data); i++ {
		data[i] = rand.Intn(100)
	}
	if IsSorted(data) {
		t.Fatalf("terrible rand.rand")
	}
	Sort(data)
	if !IsSorted(data) {
		t.Errorf("sort didn't sort - 1M ints")
	}
}

type intPair struct {
	a, b int
}

type intPairs []intPair

// Pairs compare on a only.
func intPairCmp(x, y any) int {
	vx := *c.PtrAs[intPair](x)
	vy := *c.PtrAs[intPair](y)
	return vx.a - vy.a
}

// Record initial order in B.
func (d intPairs) initB() {
	for i := range d {
		d[i].b = i
	}
}

// InOrder checks if a-equal elements were not reordered.
// If reversed is true, expect reverse ordering.
func (d intPairs) inOrder(reversed bool) bool {
	lastA, lastB := -1, 0
	for i := 0; i < len(d); i++ {
		if lastA != d[i].a {
			lastA = d[i].a
			lastB = d[i].b
			continue
		}
		if !reversed {
			if d[i].b <= lastB {
				return false
			}
		} else {
			if d[i].b >= lastB {
				return false
			}
		}
		lastB = d[i].b
	}
	return true
}

func TestStability(t *testing.T) {
	n, m := 1000, 100
	data := make(intPairs, n)

	// random distribution
	for i := 0; i < len(data); i++ {
		data[i].a = rand.Intn(m)
	}
	if IsSortedFunc(data, intPairCmp) {
		t.Fatalf("terrible rand.rand")
	}
	data.initB()
	SortStableFunc(data, intPairCmp)
	if !IsSortedFunc(data, intPairCmp) {
		t.Errorf("Stable didn't sort %d ints", n)
	}
	if !data.inOrder(false) {
		t.Errorf("Stable wasn't stable on %d ints", n)
	}

	// already sorted
	data.initB()
	SortStableFunc(data, intPairCmp)
	if !IsSortedFunc(data, intPairCmp) {
		t.Errorf("Stable shuffled sorted %d ints (order)", n)
	}
	if !data.inOrder(false) {
		t.Errorf("Stable shuffled sorted %d ints (stability)", n)
	}

	// sorted reversed
	for i := 0; i < len(data); i++ {
		data[i].a = len(data) - i
	}
	data.initB()
	SortStableFunc(data, intPairCmp)
	if !IsSortedFunc(data, intPairCmp) {
		t.Errorf("Stable didn't sort %d ints", n)
	}
	if !data.inOrder(false) {
		t.Errorf("Stable wasn't stable on %d ints", n)
	}
}

type S struct {
	a int
	b string
}

func cmpS(s1, s2 any) int {
	v1 := *c.PtrAs[S](s1)
	v2 := *c.PtrAs[S](s2)
	return cmp.Compare(v1.a, v2.a)
}

func TestMinMax(t *testing.T) {
	intCmp := func(a, b any) int {
		va := *c.PtrAs[int](a)
		vb := *c.PtrAs[int](b)
		return va - vb
	}

	tests := []struct {
		data    []int
		wantMin int
		wantMax int
	}{
		{[]int{7}, 7, 7},
		{[]int{1, 2}, 1, 2},
		{[]int{2, 1}, 1, 2},
		{[]int{1, 2, 3}, 1, 3},
		{[]int{3, 2, 1}, 1, 3},
		{[]int{2, 1, 3}, 1, 3},
		{[]int{2, 2, 3}, 2, 3},
		{[]int{3, 2, 3}, 2, 3},
		{[]int{0, 2, -9}, -9, 2},
	}
	for _, tt := range tests {
		t.Run(fmt.Sprintf("%v", tt.data), func(t *testing.T) {
			gotMin := Min(tt.data)
			if gotMin != tt.wantMin {
				t.Errorf("Min got %v, want %v", gotMin, tt.wantMin)
			}

			gotMinFunc := MinFunc(tt.data, intCmp)
			if gotMinFunc != tt.wantMin {
				t.Errorf("MinFunc got %v, want %v", gotMinFunc, tt.wantMin)
			}

			gotMax := Max(tt.data)
			if gotMax != tt.wantMax {
				t.Errorf("Max got %v, want %v", gotMax, tt.wantMax)
			}

			gotMaxFunc := MaxFunc(tt.data, intCmp)
			if gotMaxFunc != tt.wantMax {
				t.Errorf("MaxFunc got %v, want %v", gotMaxFunc, tt.wantMax)
			}
		})
	}

	svals := []S{
		{1, "a"},
		{2, "a"},
		{1, "b"},
		{2, "b"},
	}

	gotMin := MinFunc(svals, cmpS)
	wantMin := S{1, "a"}
	if gotMin != wantMin {
		t.Errorf("MinFunc(%v) = %v, want %v", svals, gotMin, wantMin)
	}

	gotMax := MaxFunc(svals, cmpS)
	wantMax := S{2, "a"}
	if gotMax != wantMax {
		t.Errorf("MaxFunc(%v) = %v, want %v", svals, gotMax, wantMax)
	}
}

func TestMinMaxPanics(t *testing.T) {
	intCmp := func(a, b any) int {
		va := *c.PtrAs[int](a)
		vb := *c.PtrAs[int](b)
		return va - vb
	}
	emptySlice := []int{}

	if !panics(func() { _ = Min(emptySlice) }) {
		t.Errorf("Min([]): got no panic, want panic")
	}

	if !panics(func() { _ = Max(emptySlice) }) {
		t.Errorf("Max([]): got no panic, want panic")
	}

	if !panics(func() { _ = MinFunc(emptySlice, intCmp) }) {
		t.Errorf("MinFunc([]): got no panic, want panic")
	}

	if !panics(func() { _ = MaxFunc(emptySlice, intCmp) }) {
		t.Errorf("MaxFunc([]): got no panic, want panic")
	}
}

func panics(f func()) (b bool) {
	defer func() {
		if x := recover(); x != nil {
			b = true
		}
	}()
	f()
	return false
}
