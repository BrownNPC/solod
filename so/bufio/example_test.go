// Copyright 2013 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bufio_test

import (
	"solod.dev/so/bufio"
	"solod.dev/so/bytes"
	"solod.dev/so/fmt"
	"solod.dev/so/io"
	"solod.dev/so/mem"
	"solod.dev/so/os"
	"solod.dev/so/strings"
)

func ExampleReader() {
	s := "Hello, world!\nThis is a bufio.Reader example."
	sr := strings.NewReader(s)
	r := bufio.NewReader(nil, &sr)
	defer r.Free()

	line, err := r.ReadString('\n')
	defer mem.FreeString(nil, line)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Line: %s", line)

	rest, err := r.ReadString(0) // Read until EOF
	defer mem.FreeString(nil, rest)
	if err != nil && err != io.EOF {
		panic(err)
	}
	fmt.Printf("Rest: %s", rest)
	// Output:
	// Line: Hello, world!
	// Rest: This is a bufio.Reader example.
}

func ExampleWriter() {
	w := bufio.NewWriter(nil, os.Stdout)
	defer w.Free()
	fmt.Fprintf(&w, "Hello, ")
	fmt.Fprintf(&w, "world!")
	w.Flush() // Don't forget to flush!
}

// ExampleWriter_ReadFrom demonstrates how to use the ReadFrom method of Writer.
func ExampleWriter_ReadFrom() {
	buf := bytes.NewBuffer(nil, nil)
	defer buf.Free()
	writer := bufio.NewWriter(nil, &buf)
	defer writer.Free()

	data := "Hello, world!\nThis is a ReadFrom example."
	reader := strings.NewReader(data)

	n, err := writer.ReadFrom(&reader)
	if err != nil {
		panic(err)
	}

	if err = writer.Flush(); err != nil {
		panic(err)
	}

	fmt.Printf("Bytes written: %d\n", n)
	fmt.Printf("Buffer contents: %s\n", buf.String())
	// Output:
	// Bytes written: 41
	// Buffer contents: Hello, world!
	// This is a ReadFrom example.
}
