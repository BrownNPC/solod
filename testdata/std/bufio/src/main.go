package main

import (
	"solod.dev/so/bufio"
	"solod.dev/so/bytes"
	"solod.dev/so/mem"
	"solod.dev/so/strings"
)

func main() {
	{
		// Writer -> Buffer -> Reader pipeline.
		var buf bytes.Buffer
		buf = bytes.NewBuffer(nil, nil)
		w := bufio.NewWriter(nil, &buf)
		w.WriteString("Hello, ")
		w.WriteString("World!")
		w.WriteByte('\n')
		w.Flush()
		w.Free()

		sr := strings.NewReader(buf.String())
		r := bufio.NewReader(nil, &sr)
		line, err := r.ReadString('\n')
		if err != nil {
			panic("ReadString failed")
		}
		if line != "Hello, World!\n" {
			panic("unexpected line")
		}
		mem.FreeString(nil, line)
		r.Free()
		buf.Free()
	}
	{
		// ReadByte and UnreadByte.
		sr := strings.NewReader("abc")
		r := bufio.NewReader(nil, &sr)
		b, err := r.ReadByte()
		if err != nil || b != 'a' {
			panic("ReadByte failed")
		}
		err = r.UnreadByte()
		if err != nil {
			panic("UnreadByte failed")
		}
		b, err = r.ReadByte()
		if err != nil || b != 'a' {
			panic("UnreadByte re-read failed")
		}
		r.Free()
	}
	{
		// Peek.
		sr := strings.NewReader("hello")
		r := bufio.NewReader(nil, &sr)
		p, err := r.Peek(3)
		if err != nil || string(p) != "hel" {
			panic("Peek failed")
		}
		r.Free()
	}
	{
		// WriteRune.
		var buf bytes.Buffer
		buf = bytes.NewBuffer(nil, nil)
		w := bufio.NewWriter(nil, &buf)
		w.WriteRune('A')
		w.Flush()
		if buf.String() != "A" {
			panic("WriteRune failed")
		}
		w.Free()
		buf.Free()
	}
	{
		// Scanner.
		sr := strings.NewReader("line1\nline2\n")
		s := bufio.NewScanner(nil, &sr)
		count := 0
		for s.Scan() {
			if count == 0 && s.Text() != "line1" {
				panic("Scanner: expected line1")
			}
			if count == 1 && s.Text() != "line2" {
				panic("Scanner: expected line2")
			}
			count++
		}
		if count != 2 {
			panic("Scanner: expected 2 lines")
		}
		s.Free()
	}
}
