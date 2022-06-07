package main

import (
	"testing"

	"github.com/letsblockit/letsblockit/src/server"
	"github.com/stretchr/testify/assert"
)

func TestServerDryRun(t *testing.T) {
	assert.Equal(t, server.ErrDryRunFinished, server.NewServer(&server.Options{
		AuthMethod:   "kratos",
		DatabaseUrl:  "postgresql:///letsblockit",
		StatsdTarget: "localhost:8125",
		DryRun:       true,
		HotReload:    true,
	}).Start())
}
