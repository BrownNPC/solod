package main

import "solod.dev/so/mem"

func arenaTest() {
	{
		// Arena allocator.
		buf := make([]byte, 1024)
		arena := mem.NewArena(buf)
		var a mem.Allocator = &arena

		// Allocate a Point.
		p, err := mem.TryAlloc[Point](a)
		if err != nil {
			panic("initial allocation failed")
		}
		p.x = 11
		p.y = 22
		if p.x != 11 || p.y != 22 {
			panic("unexpected p.x or p.y")
		}

		// Free last allocation reclaims space.
		mem.Free(a, p)

		// Allocate again: should reuse the same space.
		p2, err := mem.TryAlloc[Point](a)
		if err != nil {
			panic("allocation after free failed")
		}
		// Memory should be zeroed.
		if p2.x != 0 || p2.y != 0 {
			panic("memory not zeroed after free")
		}
		p2.x = 33
		p2.y = 44

		// Free non-last allocation is a no-op.
		p3, err := mem.TryAlloc[Point](a)
		if err != nil {
			panic("allocation for p3 failed")
		}
		p3.x = 55
		p3.y = 66
		mem.Free(a, p2) // not last, no-op

		// Reset and reallocate.
		arena.Reset()
		p4, err := mem.TryAlloc[Point](a)
		if err != nil {
			panic("allocation after reset failed")
		}
		if p4.x != 0 || p4.y != 0 {
			panic("memory not zeroed after reset")
		}
		p4.x = 77
		p4.y = 88
	}
}
