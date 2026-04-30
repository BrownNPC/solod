package main

import "solod.dev/so/flag"

func main() {
	flags := flag.NewFlagSet("example", flag.ContinueOnError)
	var b bool
	flags.BoolVar(&b, "b", false, "a boolean flag")
	var n int
	flags.IntVar(&n, "n", 0, "an int flag")
	var f float64
	flags.Float64Var(&f, "f", 0.0, "a float flag")
	var s string
	flags.StringVar(&s, "s", "default", "a string flag")

	err := flags.Parse([]string{"-b", "-n", "42", "-f", "3.14", "-s", "hello"})
	if err != nil {
		panic(err)
	}

	if !b {
		panic("b != true")
	}
	if n != 42 {
		panic("n != 42")
	}
	if f != 3.14 {
		panic("f != 3.14")
	}
	if s != "hello" {
		panic("s != hello")
	}
}
