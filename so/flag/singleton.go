// Copyright 2009 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package flag

import (
	"solod.dev/so/os"
)

// CommandLine is the default set of command-line flags, parsed from [os.Args].
// The top-level functions such as [BoolVar], [Arg], and so on are wrappers for the
// methods of CommandLine.
var commandLine FlagSet
var CommandLine *FlagSet

// Args returns the non-flag command-line arguments.
func Args() []string {
	initCommandLine()
	return CommandLine.args
}

// Parse parses the command-line flags from [os.Args][1:]. Must be called
// after all flags are defined and before flags are accessed by the program.
func Parse() {
	initCommandLine()
	// Ignore errors; CommandLine is set for ExitOnError.
	CommandLine.Parse(os.Args[1:])
}

// Usage prints a usage message documenting all defined command-line flags
// to [CommandLine]'s output, which by default is [os.Stderr].
// It is called when an error occurs while parsing flags.
// The function is a variable that may be changed to point to a custom function.
// By default it prints a simple header and calls [PrintDefaults]; for details about the
// format of the output and how to control it, see the documentation for [PrintDefaults].
// Custom usage functions may choose to exit the program; by default exiting
// happens anyway as the command line's error handling strategy is set to
// [ExitOnError].
func Usage() {
	initCommandLine()
	CommandLine.Usage()
}

// BoolVar defines a bool flag with specified name, default value, and usage string.
// The argument p points to a bool variable in which to store the value of the flag.
func BoolVar(p *bool, name string, value bool, usage string) {
	initCommandLine()
	CommandLine.BoolVar(p, name, value, usage)
}

// IntVar defines an int flag with specified name, default value, and usage string.
// The argument p points to an int variable in which to store the value of the flag.
func IntVar(p *int, name string, value int, usage string) {
	initCommandLine()
	CommandLine.IntVar(p, name, value, usage)
}

// Int64Var defines an int64 flag with specified name, default value, and usage string.
// The argument p points to an int64 variable in which to store the value of the flag.
func Int64Var(p *int64, name string, value int64, usage string) {
	initCommandLine()
	CommandLine.Int64Var(p, name, value, usage)
}

// UintVar defines a uint flag with specified name, default value, and usage string.
// The argument p points to a uint variable in which to store the value of the flag.
func UintVar(p *uint, name string, value uint, usage string) {
	initCommandLine()
	CommandLine.UintVar(p, name, value, usage)
}

// Uint64Var defines a uint64 flag with specified name, default value, and usage string.
// The argument p points to a uint64 variable in which to store the value of the flag.
func Uint64Var(p *uint64, name string, value uint64, usage string) {
	initCommandLine()
	CommandLine.Uint64Var(p, name, value, usage)
}

// Float64Var defines a float64 flag with specified name, default value, and usage string.
// The argument p points to a float64 variable in which to store the value of the flag.
func Float64Var(p *float64, name string, value float64, usage string) {
	initCommandLine()
	CommandLine.Float64Var(p, name, value, usage)
}

// StringVar defines a string flag with specified name, default value, and usage string.
// The argument p points to a string variable in which to store the value of the flag.
func StringVar(p *string, name string, value string, usage string) {
	initCommandLine()
	CommandLine.StringVar(p, name, value, usage)
}

// Var defines a flag with the specified name and usage string. The type and
// value of the flag are represented by the first argument, of type [Value], which
// typically holds a user-defined implementation of [Value]. For instance, the
// caller could create a flag that turns a comma-separated string into a slice
// of strings by giving the slice the methods of [Value]; in particular, [Set] would
// decompose the comma-separated string into the slice.
func Var(value Value, name string, usage string) {
	initCommandLine()
	CommandLine.Var(value, name, usage)
}

// initCommandLine initializes CommandLine if it has not already been initialized.
func initCommandLine() {
	if CommandLine != nil {
		return
	}
	// It's possible for execl to hand us an empty os.Args.
	if len(os.Args) == 0 {
		commandLine = NewFlagSet("", ExitOnError)
	} else {
		commandLine = NewFlagSet(os.Args[0], ExitOnError)
	}
	CommandLine = &commandLine
}
