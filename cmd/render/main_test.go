package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRenderFromFile(t *testing.T) {
	out, err := strings.Builder{}, strings.Builder{}
	stdout = &out
	stderr = &err

	cmd := &renderCmd{
		Input: "testdata/input.yaml",
	}
	assert.NoError(t, cmd.Run())

	expected, e := os.ReadFile("testdata/expected.txt")
	assert.NoError(t, e)
	assert.Equal(t, string(expected), out.String())
	assert.Equal(t, "WARNING: skipping unknown: template 'unknown' not found\n", err.String())
}
