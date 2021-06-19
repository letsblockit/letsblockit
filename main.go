package main

import (
	"github.com/xvello/weblock/server"
)

func main() {
	if err := server.NewServer().Start(); err != nil {
		panic(err)
	}
}
