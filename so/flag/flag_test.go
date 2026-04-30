// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"testing"

	"solod.dev/so/io"
	"solod.dev/so/os"
	"solod.dev/so/strings"

	"github.com/nalgeon/be"
)

func TestDefine(t *testing.T) {
	flags := NewFlagSet("test", ContinueOnError)

	var b bool
	flags.BoolVar(&b, "test_bool", false, "bool value")
	var i int
	flags.IntVar(&i, "test_int", 0, "int value")
	var i64 int64
	flags.Int64Var(&i64, "test_int64", 0, "int64 value")
	var u uint
	flags.UintVar(&u, "test_uint", 0, "uint value")
	var u64 uint64
	flags.Uint64Var(&u64, "test_uint64", 0, "uint64 value")
	var f float64
	flags.Float64Var(&f, "test_float64", 0, "float64 value")
	var s string
	flags.StringVar(&s, "test_string", "0", "string value")

	be.Equal(t, flags.NFlag(), 7)
}

func TestGetSet(t *testing.T) {
	flags := NewFlagSet("test", ContinueOnError)

	var b bool
	flags.BoolVar(&b, "test_bool", false, "bool value")
	var i int
	flags.IntVar(&i, "test_int", 0, "int value")
	var i64 int64
	flags.Int64Var(&i64, "test_int64", 0, "int64 value")
	var u uint
	flags.UintVar(&u, "test_uint", 0, "uint value")
	var u64 uint64
	flags.Uint64Var(&u64, "test_uint64", 0, "uint64 value")
	var f float64
	flags.Float64Var(&f, "test_float64", 0, "float64 value")
	var s string
	flags.StringVar(&s, "test_string", "0", "string value")

	var err error
	err = flags.Set("test_bool", "true")
	be.Err(t, err, nil)
	be.Equal(t, b, true)

	err = flags.Set("test_int", "11")
	be.Err(t, err, nil)
	be.Equal(t, i, 11)

	err = flags.Set("test_int64", "22")
	be.Err(t, err, nil)
	be.Equal(t, i64, 22)

	err = flags.Set("test_uint", "33")
	be.Err(t, err, nil)
	be.Equal(t, u, 33)

	err = flags.Set("test_uint64", "44")
	be.Err(t, err, nil)
	be.Equal(t, u64, 44)

	err = flags.Set("test_float64", "1.1")
	be.Err(t, err, nil)
	be.Equal(t, f, 1.1)

	err = flags.Set("test_string", "hello")
	be.Err(t, err, nil)
	be.Equal(t, s, "hello")

	fb := flags.Lookup("test_bool")
	be.Equal(t, fb.Value.Get(), any(&b))
	fi := flags.Lookup("test_int")
	be.Equal(t, fi.Value.Get(), any(&i))
	fi64 := flags.Lookup("test_int64")
	be.Equal(t, fi64.Value.Get(), any(&i64))
	fu := flags.Lookup("test_uint")
	be.Equal(t, fu.Value.Get(), any(&u))
	fu64 := flags.Lookup("test_uint64")
	be.Equal(t, fu64.Value.Get(), any(&u64))
	ff := flags.Lookup("test_float64")
	be.Equal(t, ff.Value.Get(), any(&f))
	fs := flags.Lookup("test_string")
	be.Equal(t, fs.Value.Get(), any(&s))
}

func TestParse(t *testing.T) {
	flags := NewFlagSet("test", ContinueOnError)
	if flags.Parsed() {
		t.Error("flags.Parse() = true before Parse")
	}

	var b bool
	flags.BoolVar(&b, "bool", false, "bool value")
	var b2 bool
	flags.BoolVar(&b2, "bool2", false, "bool2 value")
	var i int
	flags.IntVar(&i, "int", 0, "int value")
	var i64 int64
	flags.Int64Var(&i64, "int64", 0, "int64 value")
	var u uint
	flags.UintVar(&u, "uint", 0, "uint value")
	var u64 uint64
	flags.Uint64Var(&u64, "uint64", 0, "uint64 value")
	var f float64
	flags.Float64Var(&f, "float64", 0, "float64 value")
	var s string
	flags.StringVar(&s, "string", "0", "string value")

	extra := "one-extra-argument"
	args := []string{
		"-bool",
		"-bool2=true",
		"--int", "22",
		"--int64", "0x23",
		"-uint", "24",
		"--uint64", "25",
		"-string", "hello",
		"-float64", "2718e28",
		extra,
	}
	err := flags.Parse(args)
	be.Err(t, err, nil)
	be.Equal(t, flags.Parsed(), true)

	be.Equal(t, b, true)
	be.Equal(t, b2, true)
	be.Equal(t, i, 22)
	be.Equal(t, i64, 0x23)
	be.Equal(t, u, 24)
	be.Equal(t, u64, 25)
	be.Equal(t, f, 2718e28)
	be.Equal(t, s, "hello")

	be.Equal(t, len(flags.Args()), 1)
	be.Equal(t, flags.Args()[0], extra)
}

// Declare a user-defined flag type.
type flagVar []string

func (f *flagVar) Get() any {
	return []string(*f)
}

func (f *flagVar) Set(value string) error {
	*f = append(*f, value)
	return nil
}

func (f *flagVar) Type() string {
	return "list"
}

func TestUserDefined(t *testing.T) {
	var flags FlagSet
	flags.Init("test", ContinueOnError)
	flags.SetOutput(io.Discard)

	var v flagVar
	flags.Var(&v, "v", "usage")

	err := flags.Parse([]string{"-v", "1", "-v", "2", "-v=3"})
	be.Err(t, err, nil)
	be.Equal(t, len(v), 3)

	be.Equal(t, v[0], "1")
	be.Equal(t, v[1], "2")
	be.Equal(t, v[2], "3")
}

func TestSetOutput(t *testing.T) {
	var flags FlagSet
	var buf strings.Builder
	flags.SetOutput(&buf)
	flags.Init("test", ContinueOnError)
	flags.Parse([]string{"-unknown"})
	if out := buf.String(); !strings.Contains(out, "-unknown") {
		t.Errorf("expected output mentioning unknown; got %q", out)
	}
}

const defaultOutput = `  -A	for bootstrapping, allow 'any' type
  -Alongflagname
    	disable bounds checking
  -C	a boolean defaulting to true (default true)
  -D string
    	set relative path for local imports
  -E string
    	issue 23543 (default "0")
  -F float
    	a non-zero number (default 2.7)
  -G float
    	a float that defaults to zero
  -N int
    	a non-zero int (default 27)
  -V list
    	a list of strings (default [a b])
  -Z int
    	an int that defaults to zero
`

func TestPrintDefaults(t *testing.T) {
	fs := NewFlagSet("print defaults test", ContinueOnError)
	var buf strings.Builder
	fs.SetOutput(&buf)
	var a bool
	fs.BoolVar(&a, "A", false, "for bootstrapping, allow 'any' type")
	var alongflagname bool
	fs.BoolVar(&alongflagname, "Alongflagname", false, "disable bounds checking")
	var c bool
	fs.BoolVar(&c, "C", true, "a boolean defaulting to true (default true)")
	var d string
	fs.StringVar(&d, "D", "", "set relative path for local imports")
	var e string
	fs.StringVar(&e, "E", "0", "issue 23543 (default \"0\")")
	var f float64
	fs.Float64Var(&f, "F", 2.7, "a non-zero number (default 2.7)")
	var g float64
	fs.Float64Var(&g, "G", 0, "a float that defaults to zero")
	var n int
	fs.IntVar(&n, "N", 27, "a non-zero int (default 27)")
	var v flagVar
	fs.Var(&v, "V", "a list of strings (default [a b])")
	var z int
	fs.IntVar(&z, "Z", 0, "an int that defaults to zero")
	fs.PrintDefaults()
	got := buf.String()
	if got != defaultOutput {
		t.Errorf("got:\n%q\nwant:\n%q", got, defaultOutput)
	}
}

func TestGetters(t *testing.T) {
	expectedName := "flag set"
	expectedErrorHandling := ContinueOnError
	expectedOutput := io.Writer(os.Stderr)
	fs := NewFlagSet(expectedName, expectedErrorHandling)

	if fs.Name() != expectedName {
		t.Errorf("unexpected name: got %s, expected %s", fs.Name(), expectedName)
	}
	if fs.ErrorHandling() != expectedErrorHandling {
		t.Errorf("unexpected ErrorHandling: got %d, expected %d", fs.ErrorHandling(), expectedErrorHandling)
	}
	if fs.Output() != expectedOutput {
		t.Errorf("unexpected output: got %#v, expected %#v", fs.Output(), expectedOutput)
	}

	expectedName = "gopher"
	expectedErrorHandling = ExitOnError
	expectedOutput = os.Stdout
	fs.Init(expectedName, expectedErrorHandling)
	fs.SetOutput(expectedOutput)

	if fs.Name() != expectedName {
		t.Errorf("unexpected name: got %s, expected %s", fs.Name(), expectedName)
	}
	if fs.ErrorHandling() != expectedErrorHandling {
		t.Errorf("unexpected ErrorHandling: got %d, expected %d", fs.ErrorHandling(), expectedErrorHandling)
	}
	if fs.Output() != expectedOutput {
		t.Errorf("unexpected output: got %v, expected %v", fs.Output(), expectedOutput)
	}
}

func TestParseError(t *testing.T) {
	var b bool
	var i int
	var i64 int64
	var u uint
	var u64 uint64
	var f float64
	for _, typ := range []string{"bool", "int", "int64", "uint", "uint64", "float64"} {
		fs := NewFlagSet("parse error test", ContinueOnError)
		fs.SetOutput(io.Discard)
		fs.BoolVar(&b, "bool", false, "")
		fs.IntVar(&i, "int", 0, "")
		fs.Int64Var(&i64, "int64", 0, "")
		fs.UintVar(&u, "uint", 0, "")
		fs.Uint64Var(&u64, "uint64", 0, "")
		fs.Float64Var(&f, "float64", 0, "")
		// Strings cannot give errors.
		args := []string{"-" + typ + "=x"}
		err := fs.Parse(args) // x is not a valid setting for any flag.
		if err != ErrParse {
			t.Errorf("Parse(%q)=%v; expected parse error", args, err)
		}
	}
}

func TestRangeError(t *testing.T) {
	bad := []string{
		"-int=123456789012345678901",
		"-int64=123456789012345678901",
		"-uint=123456789012345678901",
		"-uint64=123456789012345678901",
		"-float64=1e1000",
	}
	var i int
	var i64 int64
	var u uint
	var u64 uint64
	var f float64
	for _, arg := range bad {
		fs := NewFlagSet("parse error test", ContinueOnError)
		fs.SetOutput(io.Discard)
		fs.IntVar(&i, "int", 0, "")
		fs.Int64Var(&i64, "int64", 0, "")
		fs.UintVar(&u, "uint", 0, "")
		fs.Uint64Var(&u64, "uint64", 0, "")
		fs.Float64Var(&f, "float64", 0, "")
		// Strings cannot give errors, and bools and durations do not return strconv.NumError.
		err := fs.Parse([]string{arg})
		if err != ErrRange {
			t.Errorf("Parse(%q)=%v; expected range error", arg, err)
		}
	}
}
