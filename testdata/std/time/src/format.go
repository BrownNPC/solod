package main

import "solod.dev/so/time"

func format() {
	t := time.Date(2024, time.March, 15, 14, 30, 45, 0, time.UTC)
	{
		// RFC3339.
		buf := make([]byte, time.RFC3339Len)
		s := t.Format(buf, time.RFC3339, time.UTC)
		if s != "2024-03-15T14:30:45Z" {
			panic("unexpected RFC3339 format")
		}
	}
	{
		// RFC3339Nano.
		buf := make([]byte, time.RFC3339NanoLen)
		t = time.Date(2024, time.March, 15, 14, 30, 45, 123456789, time.UTC)
		s := t.Format(buf, time.RFC3339Nano, time.UTC)
		if s != "2024-03-15T14:30:45.123456789Z" {
			panic("unexpected RFC3339Nano format")
		}
	}
	{
		// DateTime.
		buf := make([]byte, time.DateTimeLen)
		s := t.Format(buf, time.DateTime, time.UTC)
		if s != "2024-03-15 14:30:45" {
			panic("unexpected DateTime format")
		}
	}
	{
		// DateOnly.
		buf := make([]byte, time.DateOnlyLen)
		s := t.Format(buf, time.DateOnly, time.UTC)
		if s != "2024-03-15" {
			panic("unexpected DateOnly format")
		}
	}
	{
		// TimeOnly.
		buf := make([]byte, time.TimeOnlyLen)
		s := t.Format(buf, time.TimeOnly, time.UTC)
		if s != "14:30:45" {
			panic("unexpected TimeOnly format")
		}
	}
	{
		// Custom format.
		buf := make([]byte, len("15.03.2024")+1)
		s := t.Format(buf, "%d.%m.%Y", time.UTC)
		if s != "15.03.2024" {
			panic("unexpected custom format")
		}
	}
	{
		// Time.String.
		buf := make([]byte, time.RFC3339Len)
		s := t.String(buf)
		if s != "2024-03-15T14:30:45Z" {
			panic("unexpected String format")
		}
	}
}
