// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import "solod.dev/so/strconv"

// -- bool Value
type boolValue bool

func (b *boolValue) Set(s string) error {
	v, err := strconv.ParseBool(s)
	if err != nil {
		err = ErrParse
	}
	*b = boolValue(v)
	return err
}

func (b *boolValue) Get() any { return (*bool)(b) }

func (*boolValue) Type() string { return "bool" }

// -- int Value
type intValue int

func (i *intValue) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
	*i = intValue(v)
	return err
}

func (i *intValue) Get() any { return (*int)(i) }

func (*intValue) Type() string { return "int" }

// -- int64 Value
type int64Value int64

func (i *int64Value) Set(s string) error {
	v, err := strconv.ParseInt(s, 0, 64)
	if err != nil {
		err = numError(err)
	}
	*i = int64Value(v)
	return err
}

func (i *int64Value) Get() any { return (*int64)(i) }

func (*int64Value) Type() string { return "int" }

// -- uint Value
type uintValue uint

func (i *uintValue) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, strconv.IntSize)
	if err != nil {
		err = numError(err)
	}
	*i = uintValue(v)
	return err
}

func (i *uintValue) Get() any { return (*uint)(i) }

func (*uintValue) Type() string { return "uint" }

// -- uint64 Value
type uint64Value uint64

func (i *uint64Value) Set(s string) error {
	v, err := strconv.ParseUint(s, 0, 64)
	if err != nil {
		err = numError(err)
	}
	*i = uint64Value(v)
	return err
}

func (i *uint64Value) Get() any { return (*uint64)(i) }

func (*uint64Value) Type() string { return "uint" }

// -- float64 Value
type float64Value float64

func (f *float64Value) Set(s string) error {
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		err = numError(err)
	}
	*f = float64Value(v)
	return err
}

func (f *float64Value) Get() any { return (*float64)(f) }

func (*float64Value) Type() string { return "float" }

// -- string Value
type stringValue string

func (s *stringValue) Set(val string) error {
	*s = stringValue(val)
	return nil
}

func (s *stringValue) Get() any { return (*string)(s) }

func (*stringValue) Type() string { return "string" }

func numError(err error) error {
	if err == strconv.ErrSyntax {
		return ErrParse
	}
	if err == strconv.ErrRange {
		return ErrRange
	}
	return err
}
