package main

import "github.com/alecthomas/kong"

var cli struct {
	ExtractIcons  extractIconsCmd  `cmd:"" help:"Extract icon data into yaml data."`
	FilterLint    filterLintCmd    `cmd:"" help:"Run lints and tests on filter data."`
	UpdatePresets updatePresetsCmd `cmd:"" help:"Update template preset values."`
}

func main() {
	ctx := kong.Parse(&cli)
	ctx.FatalIfErrorf(ctx.Run())
}
