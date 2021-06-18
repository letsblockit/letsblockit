package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Check that all web templates parse OK
func TestLoadTemplates(t *testing.T) {
	templates, err := loadTemplates()
	assert.NoError(t, err)
	require.Greater(t, len(templates.templates), 0, "Expected at least one template")
}
