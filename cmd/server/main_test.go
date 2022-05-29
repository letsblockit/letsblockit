package main

import (
	"os"
	"testing"

	"github.com/letsblockit/letsblockit/src/server"
	"github.com/stretchr/testify/assert"
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
