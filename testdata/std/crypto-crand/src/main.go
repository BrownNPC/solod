package main

import "solod.dev/so/crypto/crand"

func main() {
	{
		// Read.
		buf := make([]byte, 16)
		n, err := crand.Read(buf)
		if err != nil {
			panic("failed to read random data")
		}
		if n != len(buf) {
			panic("short read of random data")
		}
	}
	{
		// Read empty slice.
		buf := make([]byte, 0)
		n, err := crand.Read(buf)
		if err != nil {
			panic("failed to read random data")
		}
		if n != 0 {
			panic("non-zero read of empty slice")
		}
	}
	{
		// Reader.
		buf := make([]byte, 16)
		n, err := crand.Reader.Read(buf)
		if err != nil {
			panic("failed to read random data")
		}
		if n != len(buf) {
			panic("short read of random data")
		}
	}
	{
		// Text.
		buf := make([]byte, 26)
		s := crand.Text(buf)
		if len(s) != 26 {
			panic("unexpected length of random text")
		}
		println(s)
	}
}
