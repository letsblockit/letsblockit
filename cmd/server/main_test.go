package main

import (
	"testing"

	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/server"
	"github.com/stretchr/testify/assert"
)

func TestServerDryRun(t *testing.T) {
	assert.Equal(t, server.ErrDryRunFinished, server.NewServer(&server.Options{
		AuthMethod:    "kratos",
		AuthKratosUrl: "http://localhost:4000/.ory",
		DatabaseUrl:   db.GetTestDatabaseURL(),
		StatsdTarget:  "localhost:8125",
		DryRun:        true,
		HotReload:     true,
	}).Start())
}
