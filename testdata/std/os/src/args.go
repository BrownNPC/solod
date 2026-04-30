package main

import "solod.dev/so/os"

func argsTest() {
	// os.Args should be populated.
	if len(os.Args) == 0 {
		panic("os.Args: empty")
	}
	// First arg (program name) should be non-empty.
	if len(os.Args[0]) == 0 || os.Args[0] == "" {
		panic("os.Args[0]: empty")
	}
}
