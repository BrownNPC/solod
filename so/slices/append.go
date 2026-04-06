package slices

import (
	"solod.dev/so/c"
	"solod.dev/so/mem"
)

// Append appends elements to a heap-allocated slice, growing it if needed.
// If the allocator is nil, uses the system allocator.
// Returns an updated allocated slice; the caller owns it.
//
//so:extern
func Append[T any](a mem.Allocator, s []T, elems ...T) []T {
	return append(s, elems...)
}

// Extend appends all elements from another heap-allocated slice, growing if needed.
// If the allocator is nil, uses the system allocator.
// Returns an updated allocated slice; the caller owns it.
//
//so:extern
func Extend[T any](a mem.Allocator, s []T, other []T) []T {
	return append(s, other...)
}

// nextcap computes the capacity for a grown slice using Go's growth
// formula: 2x for small slices (< 256 elements), transitioning to ~1.25x
// for larger ones.
//
//so:inline
func slices_nextcap(newLen, oldCap uintptr) uintptr {
	newCap := oldCap
	doubleCap := newCap + newCap
	if newLen > doubleCap {
		return newLen
	}
	const threshold = 256
	if oldCap < threshold {
		return doubleCap
	}
	for {
		newCap += (newCap + 3*threshold) >> 2
		if newCap >= newLen {
			break
		}
	}
	return newCap
}

// grow grows a slice's backing allocation to hold at least newLen elements.
// Returns a result with the updated slice or an error if reallocation fails.
// If the allocator is nil, uses the system allocator.
//
//so:inline
func slices_grow(a mem.Allocator, s sliceHeader, newLen, elemSize, elemAlign uintptr) sliceResult {
	if a == nil {
		a = mem.System
	}
	if newLen <= s.cap {
		return sliceResult{val: s, err: nil}
	}
	newCap := slices_nextcap(newLen, s.cap)
	oldSize := s.cap * elemSize
	newSize := newCap * elemSize
	newPtr, err := a.Realloc(s.ptr, int(oldSize), int(newSize), int(elemAlign))
	if err != nil {
		return sliceResult{val: s, err: err}
	}
	s.ptr = newPtr.(*byte)
	s.cap = newCap
	return sliceResult{val: s, err: nil}
}

// tryExtend appends all elements from another slice, growing if needed.
// Returns a result with the updated slice or an error if reallocation fails.
// If the allocator is nil, uses the system allocator.
//
//so:inline
func slices_tryExtend(a mem.Allocator, s sliceHeader, other sliceHeader, elemSize, elemAlign uintptr) sliceResult {
	res := slices_grow(a, s, s.len+other.len, elemSize, elemAlign)
	if res.err != nil {
		return res
	}
	s = res.val
	mem.Copy(c.PtrAdd(s.ptr, int(s.len*elemSize)), other.ptr, int(other.len*elemSize))
	s.len += other.len
	return sliceResult{val: s, err: nil}
}

// extend appends all elements from another slice, growing if needed.
// Returns the updated slice or panics on allocation failure.
// If the allocator is nil, uses the system allocator.
//
//so:inline
func slices_extend(a mem.Allocator, s sliceHeader, other sliceHeader, elemSize, elemAlign uintptr) sliceHeader {
	res := slices_tryExtend(a, s, other, elemSize, elemAlign)
	if res.err != nil {
		panic(res.err)
	}
	return res.val
}
