package main

import (
	"solod.dev/so/mem"
	"solod.dev/so/path"
)

func main() {
	{
		// Clean.
		cleaned := path.Clean(nil, "/opt/app/../config.json")
		if cleaned != "/opt/config.json" {
			panic("unexpected cleaned path: " + cleaned)
		}
		mem.FreeString(nil, cleaned)
	}
	{
		// Split.
		dir, file := path.Split("/opt/app/config.json")
		if dir != "/opt/app/" {
			panic("unexpected dir: " + dir)
		}
		if file != "config.json" {
			panic("unexpected file: " + file)
		}
	}
	{
		// Join.
		joined := path.Join(nil, "opt", "app", "config.json")
		if joined != "opt/app/config.json" {
			panic("unexpected path: " + joined)
		}
		mem.FreeString(nil, joined)
	}
	{
		// IsAbs.
		if !path.IsAbs("/opt/app/config.json") {
			panic("want absolute")
		}
		if path.IsAbs("opt/app/config.json") {
			panic("want not absolute")
		}
	}
	{
		// Dir.
		dir := path.Dir(nil, "/opt/app/config.json")
		if dir != "/opt/app" {
			panic("unexpected dir: " + dir)
		}
		mem.FreeString(nil, dir)
	}
	{
		// Base.
		base := path.Base("/opt/app/config.json")
		if base != "config.json" {
			panic("unexpected base: " + base)
		}
	}
	{
		// Ext.
		ext := path.Ext("/opt/app/config.json")
		if ext != ".json" {
			panic("unexpected ext: " + ext)
		}
	}
	{
		// Match.
		ok, err := path.Match("/opt/*/*.js?n", "/opt/app/config.json")
		if err != nil {
			panic(err)
		}
		if !ok {
			panic("want match")
		}
	}
}
