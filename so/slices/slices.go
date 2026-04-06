// Package slices provides various functions useful with slices of any type.
package slices

import (
	"unsafe"

	"solod.dev/so/c"
	"solod.dev/so/mem"
)

//so:embed slices.h
var slices_h string

//so:extern so_Slice
type sliceHeader struct {
	ptr *byte
	len uintptr
	cap uintptr
}

//so:extern so_R_slice_err
type sliceResult struct {
	val sliceHeader
	err error
}

// Make allocates a slice of type T with given length using allocator a.
// If the allocator is nil, uses the system allocator.
// The returned slice is allocated; the caller owns it.
//
//so:inline
func Make[T any](a mem.Allocator, len int) []T {
	return mem.AllocSlice[T](a, len, len)
}

// MakeCap allocates a slice of type T with given length and capacity using allocator a.
// If the allocator is nil, uses the system allocator.
// The returned slice is allocated; the caller owns it.
//
//so:inline
func MakeCap[T any](a mem.Allocator, len int, cap int) []T {
	return mem.AllocSlice[T](a, len, cap)
}

// Free frees a previously allocated slice.
// If the allocator is nil, uses the system allocator.
//
//so:inline
func Free[T any](a mem.Allocator, s []T) {
	mem.FreeSlice(a, s)
}

// Clone returns a shallow copy of the slice.
// If the allocator is nil, uses the system allocator.
// The returned slice is allocated; the caller owns it.
//
//so:inline
func Clone[T any](a mem.Allocator, s []T) []T {
	_s, _slen := s, len(s)
	_elemSize := c.Sizeof[T]()
	_newSlice := mem.AllocSlice[T](a, _slen, _slen)
	mem.Copy(unsafe.SliceData(_newSlice), unsafe.SliceData(_s), _slen*_elemSize)
	return _newSlice
}

// Equal reports whether two slices are equal: the same length and all
// elements equal. Empty and nil slices are considered equal.
//
//so:extern
func Equal[T comparable](s1, s2 []T) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i := range s1 {
		if s1[i] != s2[i] {
			return false
		}
	}
	return true
}
