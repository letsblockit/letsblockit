package filters

import (
	"testing"

	"github.com/letsblockit/letsblockit/data"
	"github.com/stretchr/testify/require"
)

// Check that all template definitions parse OK
func TestLoad(t *testing.T) {
	repo, err := Load(data.Templates, data.Presets)
	require.NoError(t, err, "Template parsing error")
	require.NotNil(t, repo, "Template repository is empty")
	require.Greater(t, len(repo.templateMap), 0, "Expected at least one template")
	require.Greater(t, len(repo.templateList), 0, "Expected at least one template")
	require.Greater(t, len(repo.tagList), 0, "Expected at least one tag")
}
