// Copyright 2019, Joshua J Baker
// <https://github.com/tidwall/hashmap>

// Permission to use, copy, modify, and/or distribute this software for any
// purpose with or without fee is hereby granted, provided that the above
// copyright notice and this permission notice appear in all copies.

// THE SOFTWARE IS PROVIDED "AS IS" AND THE AUTHOR DISCLAIMS ALL WARRANTIES
// WITH REGARD TO THIS SOFTWARE INCLUDING ALL IMPLIED WARRANTIES OF
// MERCHANTABILITY AND FITNESS. IN NO EVENT SHALL THE AUTHOR BE LIABLE FOR ANY
// SPECIAL, DIRECT, INDIRECT, OR CONSEQUENTIAL DAMAGES OR ANY DAMAGES
// WHATSOEVER RESULTING FROM LOSS OF USE, DATA OR PROFITS, WHETHER IN AN ACTION
// OF CONTRACT, NEGLIGENCE OR OTHER TORTIOUS ACTION, ARISING OUT OF OR IN
// CONNECTION WITH THE USE OR PERFORMANCE OF THIS SOFTWARE.

package maps

import (
	"unsafe"

	"solod.dev/so/c"
	"solod.dev/so/mem"
)

const loadFactor = 0.85 // must be above 50%

// ByteMap is a Robin Hood hashmap operating on raw byte keys and values.
// Most users will want to use the generic [Map] wrapper type instead.
type ByteMap struct {
	a mem.Allocator

	hdib  []uint64 // one per bucket: bitfield { hash:48 dib:16 }
	keys  []byte   // ksize bytes per bucket, packed
	vals  []byte   // vsize bytes per bucket, packed
	ksize int
	vsize int

	seed   uint64 // per-instance hash seed
	len    int    // number of items in the map
	mask   int    // mask for indexing into buckets
	growAt int    // length at which to grow the map
}

// NewByteMap creates a new ByteMap with the given initial capacity,
// key size, and value size, using the provided allocator (or the
// default allocator if nil). The map automatically grows as needed.
//
// If the allocator is nil, uses the system allocator.
// The caller is responsible for freeing map resources
// with [ByteMap.Free] when done using it.
func NewByteMap(a mem.Allocator, size, ksize, vsize int) ByteMap {
	m := ByteMap{a: a, ksize: ksize, vsize: vsize, seed: seed()}
	sz := 8
	// The map must be large enough to hold size entries without resizing.
	for int(float64(sz)*loadFactor) < size {
		sz *= 2
	}
	m.hdib = mem.AllocSlice[uint64](m.a, sz, sz)
	m.keys = mem.AllocSlice[byte](m.a, sz*ksize, sz*ksize)
	m.vals = mem.AllocSlice[byte](m.a, sz*vsize, sz*vsize)
	m.mask = sz - 1
	m.growAt = int(float64(sz) * loadFactor)
	return m
}

// Len returns the number of key-value pairs in the map.
func (m *ByteMap) Len() int {
	return m.len
}

// Clear removes all key-value pairs from the map, resetting
// it to an empty state. Does not free map resources;
// the map can be reused after Clear.
func (m *ByteMap) Clear() {
	clear(m.hdib)
	m.len = 0
}

// Free frees internal resources used by the map.
// If the map is already freed, does nothing.
// The map must not be used after Free.
func (m *ByteMap) Free() {
	if len(m.hdib) == 0 {
		return
	}
	mem.FreeSlice(m.a, m.hdib)
	mem.FreeSlice(m.a, m.keys)
	mem.FreeSlice(m.a, m.vals)
	m.hdib = nil
	m.keys = nil
	m.vals = nil
	m.len = 0
}

// Resize grows or reallocates the map to hold at least size entries.
func (m *ByteMap) Resize(size int) {
	nmap := NewByteMap(m.a, size, m.ksize, m.vsize)
	nmap.seed = m.seed // preserve seed so stored hashes remain valid
	rehash(&nmap, m)
	m.Free()
	*m = nmap
}

// rehash moves all entries from src into dst.
func rehash(dst, src *ByteMap) {
	hdib := unsafe.SliceData(src.hdib)
	keys := unsafe.SliceData(src.keys)
	vals := unsafe.SliceData(src.vals)
	ksize := src.ksize
	vsize := src.vsize
	n := len(src.hdib)
	for i := range n {
		hdI := c.PtrAt(hdib, i)
		if *hdI&0xFFFF > 0 {
			insert(dst, int(*hdI>>16),
				c.PtrAdd(keys, i*ksize),
				c.PtrAdd(vals, i*vsize))
		}
	}
}

// insert does byte-level Robin Hood insertion into a map.
// Used during rehash only - skips equality check since keys are unique.
func insert(m *ByteMap, h int, key any, val any) {
	ehdib := (uint64(h) << 16) | 1
	ksize := m.ksize
	vsize := m.vsize
	hdib := unsafe.SliceData(m.hdib)
	keys := unsafe.SliceData(m.keys)
	vals := unsafe.SliceData(m.vals)
	ekey := c.Alloca[byte](ksize)
	eval := c.Alloca[byte](vsize)
	mem.Copy(ekey, key, ksize)
	mem.Copy(eval, val, vsize)
	i := h & m.mask
	tmpk := c.Alloca[byte](ksize)
	tmpv := c.Alloca[byte](vsize)
	for {
		hdI := c.PtrAt(hdib, i)
		if *hdI&0xFFFF == 0 {
			*hdI = ehdib
			mem.Copy(c.PtrAdd(keys, i*ksize), ekey, ksize)
			mem.Copy(c.PtrAdd(vals, i*vsize), eval, vsize)
			m.len++
			return
		}
		if *hdI&0xFFFF < ehdib&0xFFFF {
			te := ehdib
			ehdib = *hdI
			*hdI = te
			mem.Copy(tmpk, ekey, ksize)
			mem.Copy(ekey, c.PtrAdd(keys, i*ksize), ksize)
			mem.Copy(c.PtrAdd(keys, i*ksize), tmpk, ksize)
			mem.Copy(tmpv, eval, vsize)
			mem.Copy(eval, c.PtrAdd(vals, i*vsize), vsize)
			mem.Copy(c.PtrAdd(vals, i*vsize), tmpv, vsize)
		}
		i = (i + 1) & m.mask
		ehdib++
	}
}
