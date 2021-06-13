package filters

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadFilters(t *testing.T) {
	repo, err := LoadFilters()
	require.NoError(t, err, "Filter parsing error")
	require.NotNil(t, repo, "Filter repository is empty")
	require.Greater(t, len(repo.filters), 0, "Expected at least one filter")
}

func TestRenderSimple(t *testing.T) {
	filterName := "____simple____filter____"
	repo, err := LoadFilters()
	require.NoError(t, err, "Filter parsing error")
	file, err := os.Open("testdata/simple.yaml")
	require.NoError(t, err)
	defer file.Close()
	repo.filters[filterName], err = parseFilter(filterName, file)
	require.NoError(t, err)
	params := map[string][]string {
		"string_list": []string{"one line", "another line"},
	}
	out, err := repo.Render(filterName, params)
	require.NoError(t, err)
	assert.Equal(t, "one line\nanother line\n", out)
}
