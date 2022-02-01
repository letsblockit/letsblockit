package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xvello/letsblockit/src/server"
)

func TestServerDryRun(t *testing.T) {
	// Try to use the unix socket, fallback to TCP on localhost
	pgHost := "/var/run/postgresql"
	if _, err := os.Stat(pgHost); err != nil {
		pgHost = "localhost"
	}

	assert.Equal(t, server.ErrDryRunFinished, server.NewServer(&server.Options{
		DryRun:       true,
		Reload:       true,
		Statsd:       "localhost:8125",
		DatabaseName: "letsblockit",
		DatabaseHost: pgHost,
	}).Start())
}
