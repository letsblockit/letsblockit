package main

import "github.com/alecthomas/kong"

var cli struct {
	DownloadAvatars downloadAvatarsCmd `cmd:"" help:"Download all contributor avatars."`
	ExtractIcons    extractIconsCmd    `cmd:"" help:"Extract icon data into yaml data."`
	FilterLint      filterLintCmd      `cmd:"" help:"Run lints and tests on filter data."`
	HashAssets      hashAssetsCmd      `cmd:"" help:"Update assets-hashes.json."`
	UpdatePresets   updatePresetsCmd   `cmd:"" help:"Update template preset values."`
}

func main() {
	ctx := kong.Parse(&cli)
	ctx.FatalIfErrorf(ctx.Run())
}
