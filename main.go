package main

import (
	"github.com/xvello/weblock/server"
)

func main() {
	e, err := server.SetupRouter()
	if err != nil {
		panic(err)
	}
	e.Logger.Fatal(e.Start(":8080"))
}
