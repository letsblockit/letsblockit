package filters

import (
	"fmt"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xvello/letsblockit/data"
)

func TestValidateFilters(t *testing.T) {
	repo, err := LoadFilters()
	assert.NoError(t, err)

	validate := buildValidator(t)
	seen := make(map[string]struct{}) // Ensure uniqueness of filter names

	err = data.Walk(data.Filters, filenameSuffix, func(name string, file io.Reader) error {
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
			assert.NoError(t, validate.Struct(&filter.Filter), "Filter did no pass input validation")
		})

		for i, tc := range filter.Tests {
			t.Run(fmt.Sprintf("Test/%s/%d", name, i), func(t *testing.T) {
				var buf strings.Builder
				ctx := make(map[string]interface{})
				for k, v := range tc.Params {
					ctx[k] = v
				}
				assert.NoError(t, repo.Render(&buf, filter.Name, ctx))
				assert.Equal(t, tc.Output, buf.String())
			})
		}

		return nil
	})
	assert.NoError(t, err)
}
