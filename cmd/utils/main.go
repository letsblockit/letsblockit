package main

import "github.com/alecthomas/kong"

var cli struct {
	UpdatePresets updatePresetsCmd `cmd:"" help:"Update template preset values."`
}

func main() {
	ctx := kong.Parse(&cli)
	ctx.FatalIfErrorf(ctx.Run())
}
