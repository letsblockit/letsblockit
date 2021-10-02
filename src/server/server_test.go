package server

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xvello/letsblockit/src/filters"
)

// Checks that the server can instantiate all its components.
// This includes parsing all pages and filters.
func TestServerDryRun(t *testing.T) {
	dir, err := ioutil.TempDir("", "lbi")
	require.NoError(t, err)
	defer os.RemoveAll(dir)
	server := NewServer(&Options{
		DryRun:     true,
		Migrations: true,
		DataFolder: dir,
	})
	assert.Equal(t, ErrDryRunFinished, server.Start())
}

func BenchmarkServerStart(b *testing.B) {
	dir, err := ioutil.TempDir("", "lbi")
	require.NoError(b, err)
	defer os.RemoveAll(dir)
	o := &Options{
		DryRun:     true,
		Migrations: false,
		DataFolder: dir,
	}
	for n := 0; n < b.N; n++ {
		_ = NewServer(o).Start()
	}
}

func BenchmarkLoadPages(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = loadPages()
	}
}

func BenchmarkLoadFilters(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = filters.LoadFilters()
	}
}
