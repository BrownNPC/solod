package main

import "solod.dev/so/mem"

func withDefer() {
	p := mem.Alloc[Point](nil)
	defer mem.Free(nil, p)

	p.x = 11
	p.y = 22
	if p.x != 11 || p.y != 22 {
		panic("unexpected value")
	}
}

func allocTest() {
	{
		// TryAlloc and Free.
		p, err := mem.TryAlloc[Point](mem.System)
		if err != nil {
			panic("Alloc: allocation failed")
		}
		p.x = 11
		p.y = 22
		if p.x != 11 || p.y != 22 {
			panic("Alloc: unexpected value")
		}
		mem.Free(mem.System, p)
	}
	{
		// TryAllocSlice and FreeSlice.
		slice, err := mem.TryAllocSlice[int](mem.System, 3, 3)
		if err != nil {
			panic("AllocSlice: allocation failed")
		}
		slice[0] = 11
		slice[1] = 22
		slice[2] = 33
		if slice[0] != 11 || slice[1] != 22 || slice[2] != 33 {
			panic("AllocSlice: unexpected value")
		}
		mem.FreeSlice(mem.System, slice)
	}
	{
		// Alloc/Free with default allocator.
		p := mem.Alloc[Point](nil)
		p.x = 11
		p.y = 22
		if p.x != 11 || p.y != 22 {
			panic("New: unexpected value")
		}
		mem.Free(nil, p)
	}
	{
		// AllocSlice/FreeSlice with default allocator.
		slice := mem.AllocSlice[int](nil, 3, 3)
		slice[0] = 11
		slice[1] = 22
		slice[2] = 33
		if slice[0] != 11 || slice[1] != 22 || slice[2] != 33 {
			panic("NewSlice: unexpected value")
		}
		mem.FreeSlice(nil, slice)
	}
	{
		// TryReallocSlice with explicit allocator.
		slice, err := mem.TryAllocSlice[int](mem.System, 3, 3)
		if err != nil {
			panic("ReallocSlice: initial allocation failed")
		}
		slice[0] = 11
		slice[1] = 22
		slice[2] = 33
		slice, err = mem.TryReallocSlice(mem.System, slice, 3, 6)
		if err != nil {
			panic("ReallocSlice: reallocation failed")
		}
		if len(slice) != 3 || cap(slice) != 6 {
			panic("ReallocSlice: unexpected len/cap")
		}
		if slice[0] != 11 || slice[1] != 22 || slice[2] != 33 {
			panic("ReallocSlice: data not preserved")
		}
		mem.FreeSlice(mem.System, slice)
	}
	{
		// ReallocSlice with default allocator.
		slice := mem.AllocSlice[int](nil, 2, 2)
		slice[0] = 44
		slice[1] = 55
		slice = mem.ReallocSlice(nil, slice, 4, 8)
		if len(slice) != 4 || cap(slice) != 8 {
			panic("ReallocSlice default: unexpected len/cap")
		}
		if slice[0] != 44 || slice[1] != 55 {
			panic("ReallocSlice default: data not preserved")
		}
		// New elements should be zeroed.
		if slice[2] != 0 || slice[3] != 0 {
			panic("ReallocSlice default: new elements not zeroed")
		}
		mem.FreeSlice(nil, slice)
	}
	{
		// ReallocSlice from empty slice.
		var empty []int
		slice := mem.ReallocSlice(nil, empty, 3, 4)
		if len(slice) != 3 || cap(slice) != 4 {
			panic("ReallocSlice empty: unexpected len/cap")
		}
		if slice[0] != 0 || slice[1] != 0 || slice[2] != 0 {
			panic("ReallocSlice empty: not zeroed")
		}
		mem.FreeSlice(nil, slice)
	}
	{
		// Free with nil or an empty slice.
		var p *Point
		mem.Free(nil, p)
		var empty []int
		mem.FreeSlice(nil, empty)
	}
	{
		// Free string.
		b := mem.AllocSlice[byte](nil, 3, 3)
		b[0] = 'h'
		b[1] = 'i'
		b[2] = '!'
		s1 := string(b)
		mem.FreeString(nil, s1)
		s2 := ""
		mem.FreeString(nil, s2)
	}
	{
		// Free with defer.
		withDefer()
	}
	{
		// Tracking allocator.
		track := &mem.Tracker{Allocator: mem.System}
		p := mem.Alloc[Point](track)
		mem.Free(track, p)
		if track.Stats.Alloc != 0 {
			panic("Tracker: Stats.Alloc != 0")
		}
	}
}
