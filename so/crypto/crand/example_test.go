// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package crand_test

import (
	"solod.dev/so/crypto/crand"
	"solod.dev/so/fmt"
)

func ExampleRead() {
	// Note that no error handling is necessary, as Read always succeeds.
	key := make([]byte, 32)
	crand.Read(key)
	_ = key
}

func ExampleText() {
	buf := make([]byte, 26)
	key := crand.Text(buf)
	// The key is base32 and safe to display.
	fmt.Println(key)
}
