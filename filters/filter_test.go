package filters

import (
	"fmt"
	"io/fs"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateFilters(t *testing.T) {
	validate := buildValidator(t)

	err := fs.WalkDir(inputFiles, "data", func(path string, d fs.DirEntry, err error) error {
		if d.IsDir() || !strings.HasSuffix(d.Name(), filenameSuffix){
			return nil
		}
		name := strings.TrimSuffix(d.Name(), filenameSuffix)
		var filter *filterAndTests
		t.Run("Parse | "+name, func(t *testing.T) {
			file, err := inputFiles.Open(path)
			require.NoError(t, err)
			filter, err = parseFilterAndTest(name, file)
			require.NoError(t, file.Close())
			require.NoError(t, err, "Filter did not parse OK")
			assert.NoError(t, validate.Struct(filter), "Filter did no pass input validation")
		})

		for i, tc := range filter.Tests {
			t.Run(fmt.Sprintf("Test | %s | %d", name, i), func(t *testing.T) {
				out, err := filter.Parsed.Exec(tc.Params)
				assert.NoError(t, err)
				assert.Equal(t, tc.Output, out)
			})
		}

		return nil
	})
	assert.NoError(t, err)
}
