package main

import "solod.dev/so/time"

func times() {
	{
		// time.Date and time.Time properties.
		t := time.Date(2021, time.May, 10, 12, 33, 44, 777888999, time.UTC)
		if t.Year() != 2021 {
			panic("unexpected Time.Year")
		}
		if t.Month() != time.May {
			panic("unexpected Time.Month")
		}
		if t.Day() != 10 {
			panic("unexpected Time.Day")
		}
		if t.Hour() != 12 {
			panic("unexpected Time.Hour")
		}
		if t.Minute() != 33 {
			panic("unexpected Time.Minute")
		}
		if t.Second() != 44 {
			panic("unexpected Time.Second")
		}
		if t.Nanosecond() != 777888999 {
			panic("unexpected Time.Nanosecond")
		}
		println(t.Year(), t.Month(), t.Day(), t.Hour(), t.Minute(), t.Second(), t.Nanosecond())
	}
	{
		// Time.Now.
		t := time.Now()
		if t.IsZero() {
			panic("unexpected Time.IsZero")
		}
		buf := make([]byte, time.RFC3339NanoLen)
		println("UTC:", t.String(buf))
		utc5 := time.Offset(5 * 3600)
		println("UTC+5:", t.Format(buf, time.RFC3339Nano, utc5))
	}
}
