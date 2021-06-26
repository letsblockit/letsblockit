package main

import (
	"fmt"
	"time"

	"github.com/alexflint/go-arg"
	"github.com/xvello/weblock/src/server"
)

func main() {
	start := time.Now()
	options := &server.Options{}
	arg.MustParse(options)
	err := server.NewServer(options).Start()

	switch err {
	case server.DryRunFinished:
		fmt.Printf("Dry-run checks finished in %s\n", time.Since(start))
	default:
		panic(err)
	}
}
