package main

import (
	"fmt"
	"time"

	"github.com/alecthomas/kong"
	"github.com/letsblockit/letsblockit/src/server"
)

func main() {
	start := time.Now()
	options := &server.Options{}
	k := kong.Parse(options,
		kong.Description("Read https://github.com/letsblockit/letsblockit/blob/main/cmd/server/README.md for setup instructions."),
		kong.DefaultEnvars("LETSBLOCKIT"),
	)
	err := server.NewServer(options).Start()

	if err == server.ErrDryRunFinished {
		fmt.Printf("Dry-run checks finished in %s\n", time.Since(start))
	} else {
		k.FatalIfErrorf(err)
	}
}
