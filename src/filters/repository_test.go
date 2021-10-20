package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Check that all filter definitions parse OK
func TestLoadFilters(t *testing.T) {
	repo, err := LoadFilters()
	require.NoError(t, err, "Filter parsing error")
	require.NotNil(t, repo, "Filter repository is empty")
	require.Greater(t, len(repo.filterMap), 0, "Expected at least one filter")
	require.Greater(t, len(repo.filterList), 0, "Expected at least one filter")
	require.Greater(t, len(repo.tagList), 0, "Expected at least one tag")
}
