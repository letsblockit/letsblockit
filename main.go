package main

import (
	"github.com/xvello/weblock/server"
)

func main() {
	r, err := server.SetupRouter()
	if err != nil {
		panic(err)
	}
	r.Run(":8080")
}
