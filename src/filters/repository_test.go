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
	require.Greater(t, len(repo.fMap), 0, "Expected at least one filter")
	require.Greater(t, len(repo.fList), 0, "Expected at least one filter")
}
