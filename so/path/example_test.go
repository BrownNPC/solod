// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package path_test

import (
	"solod.dev/so/fmt"
	"solod.dev/so/mem"
	"solod.dev/so/path"
)

func ExampleBase() {
	fmt.Println(path.Base("/a/b"))
	fmt.Println(path.Base("/"))
	fmt.Println(path.Base(""))
	// Output:
	// b
	// /
	// .
}

func ExampleClean() {
	cleaned := path.Clean(nil, "/opt/app/../config.json")
	fmt.Println(cleaned)
	mem.FreeString(nil, cleaned)
	// Output:
	// /opt/config.json
}

func ExampleDir() {
	dir := path.Dir(nil, "/opt/app/config.json")
	fmt.Println(dir)
	mem.FreeString(nil, dir)
	// Output:
	// /opt/app
}

func ExampleExt() {
	ext := path.Ext("/opt/app/config.json")
	fmt.Println(ext)
	// Output:
	// .json
}

func ExampleIsAbs() {
	fmt.Printf("%v\n", path.IsAbs("/dev/null"))
	// Output: true
}

func ExampleJoin() {
	joined := path.Join(nil, "opt", "app", "config.json")
	fmt.Println(joined)
	mem.FreeString(nil, joined)
	// Output:
	// opt/app/config.json
}

func ExampleMatch() {
	const pattern = "/opt/*/*.js?n"
	ok, err := path.Match(pattern, "/opt/app/config.json")
	fmt.Printf("%v %v\n", ok, err)
	// Output:
	// true <nil>
}

func ExampleSplit() {
	dir, file := path.Split("/opt/app/config.json")
	fmt.Printf("dir = %s, file = %s\n", dir, file)
	// Output:
	// dir = /opt/app/, file = config.json
}
