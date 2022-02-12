package main

import (
	"github.com/alecthomas/kong"
)

const (
	filterFolder    string = "data/filters/"
	filterExtension string = ".yaml"
)

func main() {
	options := struct {
		Presets PresetsCmd `cmd:"" help:"Sync filter presets from external sources"`
	}{}
	k := kong.Parse(&options)
	k.FatalIfErrorf(k.Run())
}
