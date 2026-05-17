package main

import (
	"solod.dev/so/strconv"
)

func main() {
	{
		// AppendBool.
		buf := make([]byte, 0, strconv.MaxBoolLen)
		b := strconv.AppendBool(buf, true)
		if string(b) != "true" {
			panic("AppendBool")
		}
	}
	{
		// AppendFloat.
		buf := make([]byte, 0, strconv.MaxFloat64Len)
		b := strconv.AppendFloat(buf, 3.1415926535, 'E', -1, 32)
		if string(b) != "3.1415927E+00" {
			panic("AppendFloat 32")
		}
		b = strconv.AppendFloat(buf, 3.1415926535, 'E', -1, 64)
		if string(b) != "3.1415926535E+00" {
			panic("AppendFloat 64")
		}
	}
	{
		// AppendInt.
		buf := make([]byte, 0, strconv.MaxIntBase10Len)
		b := strconv.AppendInt(buf, -42, 10)
		if string(b) != "-42" {
			panic("AppendInt base 10")
		}
		b = strconv.AppendInt(buf, -42, 16)
		if string(b) != "-2a" {
			panic("AppendInt base 16")
		}
	}
	{
		// AppendUint.
		buf := make([]byte, 0, strconv.MaxUintBase10Len)
		b := strconv.AppendUint(buf, 42, 10)
		if string(b) != "42" {
			panic("AppendUint base 10")
		}
		b = strconv.AppendUint(buf, 42, 16)
		if string(b) != "2a" {
			panic("AppendUint base 16")
		}
	}
	{
		// Atof.
		f, err := strconv.ParseFloat("1844674407370955", 64)
		if err != nil {
			panic("Atof error")
		}
		if f != float64(1844674407370955) {
			panic("Atof value")
		}
	}
	{
		// Atoi.
		s, err := strconv.Atoi("10")
		if err != nil {
			panic("Atoi error")
		}
		if s != 10 {
			panic("Atoi value")
		}
	}
	{
		// FormatBool.
		s := strconv.FormatBool(true)
		if s != "true" {
			panic("FormatBool")
		}
	}
	{
		// FormatFloat.
		buf := make([]byte, strconv.MaxFloat64Len)
		s := strconv.FormatFloat(buf, 3.1415926535, 'E', -1, 32)
		if s != "3.1415927E+00" {
			panic("FormatFloat 32")
		}
		s = strconv.FormatFloat(buf, 3.1415926535, 'E', -1, 64)
		if s != "3.1415926535E+00" {
			panic("FormatFloat 64")
		}
		s = strconv.FormatFloat(buf, 3.1415926535, 'g', -1, 64)
		if s != "3.1415926535" {
			panic("FormatFloat g")
		}
		s = strconv.FormatFloat(buf, 1844674407370955, 'f', -1, 64)
		if s != "1844674407370955" {
			panic("FormatFloat big")
		}
	}
	{
		// FormatInt.
		buf := make([]byte, strconv.MaxIntBase10Len)
		s := strconv.FormatInt(buf, -42, 10)
		if s != "-42" {
			panic("FormatInt base 10")
		}
		s = strconv.FormatInt(buf, -42, 16)
		if s != "-2a" {
			panic("FormatInt base 16")
		}
		s = strconv.FormatInt(buf, int64(1<<31-1), 10)
		if s != "2147483647" {
			panic("FormatInt 31bit")
		}
		s = strconv.FormatInt(buf, int64(1<<56-1), 10)
		if s != "72057594037927935" {
			panic("FormatInt 56bit")
		}
		s = strconv.FormatInt(buf, int64(1<<62-1), 10)
		if s != "4611686018427387903" {
			panic("FormatInt 62bit")
		}
	}
	{
		// FormatUint.
		buf := make([]byte, strconv.MaxUintBase10Len)
		s := strconv.FormatUint(buf, 42, 10)
		if s != "42" {
			panic("FormatUint base 10")
		}
		s = strconv.FormatUint(buf, 42, 16)
		if s != "2a" {
			panic("FormatUint base 16")
		}
	}
	{
		// Itoa.
		buf := make([]byte, strconv.MaxIntBase10Len)
		s := strconv.Itoa(buf, 10)
		if s != "10" {
			panic("Itoa")
		}
	}
	{
		// ParseBool.
		s, err := strconv.ParseBool("true")
		if err != nil {
			panic("ParseBool error")
		}
		if !s {
			panic("ParseBool value")
		}
	}
	{
		// ParseFloat.
		buf := make([]byte, strconv.MaxFloat64Len)
		s, err := strconv.ParseFloat("3.1415926535", 32)
		if err != nil {
			panic("ParseFloat 32 error")
		}
		r := strconv.FormatFloat(buf, s, 'E', -1, 32)
		if r != "3.1415927E+00" {
			panic("ParseFloat 32 value")
		}
		s, err = strconv.ParseFloat("3.1415926535", 64)
		if err != nil {
			panic("ParseFloat 64 error")
		}
		if s != 3.1415926535 {
			panic("ParseFloat 64 value")
		}
		// NaN.
		s, err = strconv.ParseFloat("NaN", 32)
		if err != nil {
			panic("ParseFloat NaN error")
		}
		if s == s {
			panic("ParseFloat NaN value")
		}
		// Case insensitive.
		s, err = strconv.ParseFloat("nan", 32)
		if err != nil {
			panic("ParseFloat nan error")
		}
		if s == s {
			panic("ParseFloat nan value")
		}
		// inf.
		s, err = strconv.ParseFloat("inf", 32)
		if err != nil {
			panic("ParseFloat inf error")
		}
		r = strconv.FormatFloat(buf, s, 'g', -1, 64)
		if r != "+Inf" {
			panic("ParseFloat inf value")
		}
		// +Inf.
		s, err = strconv.ParseFloat("+Inf", 32)
		if err != nil {
			panic("ParseFloat +Inf error")
		}
		r = strconv.FormatFloat(buf, s, 'g', -1, 64)
		if r != "+Inf" {
			panic("ParseFloat +Inf value")
		}
		// -Inf.
		s, err = strconv.ParseFloat("-Inf", 32)
		if err != nil {
			panic("ParseFloat -Inf error")
		}
		r = strconv.FormatFloat(buf, s, 'g', -1, 64)
		if r != "-Inf" {
			panic("ParseFloat -Inf value")
		}
		// -0.
		s, err = strconv.ParseFloat("-0", 32)
		if err != nil {
			panic("ParseFloat -0 error")
		}
		r = strconv.FormatFloat(buf, s, 'g', -1, 64)
		if r != "-0" {
			panic("ParseFloat -0 value")
		}
		// +0.
		s, err = strconv.ParseFloat("+0", 32)
		if err != nil {
			panic("ParseFloat +0 error")
		}
		if s != 0 {
			panic("ParseFloat +0 value")
		}
	}
	{
		// ParseInt.
		s, err := strconv.ParseInt("-354634382", 10, 32)
		if err != nil {
			panic("ParseInt 32 error")
		}
		if s != -354634382 {
			panic("ParseInt 32 value")
		}
		s, err = strconv.ParseInt("-3546343826724305832", 10, 64)
		if err != nil {
			panic("ParseInt 64 error")
		}
		if s != -3546343826724305832 {
			panic("ParseInt 64 value")
		}
	}
	{
		// ParseUint.
		s, err := strconv.ParseUint("42", 10, 32)
		if err != nil {
			panic("ParseUint 32 error")
		}
		if s != 42 {
			panic("ParseUint 32 value")
		}
		s, err = strconv.ParseUint("42", 10, 64)
		if err != nil {
			panic("ParseUint 64 error")
		}
		if s != 42 {
			panic("ParseUint 64 value")
		}
	}
}
