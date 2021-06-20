package filters

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xvello/weblock/utils"
)

func TestValidateFilters(t *testing.T) {
	validate := buildValidator(t)
	seen := make(map[string]struct{}) // Ensure uniqueness of filter names

	err := utils.Walk(definitionFiles, filenameSuffix, func(name string, file io.Reader) error {
		t.Run("Name/"+name, func(t *testing.T) {
			if name != strings.ToLower(name) {
				assert.Fail(t, "name can only be lowercase", name)
			}
			if _, found := seen[name]; found {
				assert.Fail(t, "duplicate name found", name)
			}
			seen[name] = struct{}{}
		})

		var filter *filterAndTests
		var e error
		t.Run("Parse/"+name, func(t *testing.T) {
			filter, e = parseFilterAndTest(name, file)
			require.NoError(t, e, "Filter did not parse OK")
			assert.NoError(t, validate.Struct(filter), "Filter did no pass input validation")
		})

		t.Run("Desc/"+name, func(t *testing.T) {
			assert.True(t, strings.HasPrefix(filter.Description, "<h2>"), "Description must start with a second-level header")
		})

		for i, tc := range filter.Tests {
			t.Run(fmt.Sprintf("Test/%s/%d", name, i), func(t *testing.T) {
				var buf strings.Builder
				assert.NoError(t, filter.Parsed.Execute(&buf, tc.Params))
				assert.Equal(t, tc.Output, buf.String())
			})
		}

		return nil
	})
	assert.NoError(t, err)
}
