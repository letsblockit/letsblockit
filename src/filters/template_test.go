package filters

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"strconv"
	"strings"
	"testing"

	"github.com/letsblockit/letsblockit/data"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateTemplates(t *testing.T) {
	contributors, err := data.ParseContributors()
	require.NoErrorf(t, err, "failed to load contributors")

	repo, err := Load(data.Templates, data.Presets)
	assert.NoError(t, err)

	validate := buildValidator(t)
	seen := make(map[string]struct{}) // Ensure uniqueness of template names

	err = data.Walk(data.Templates, filenameSuffix, func(name string, file io.Reader) error {
		t.Run("Name/"+name, func(t *testing.T) {
			if name != strings.ToLower(name) {
				assert.Fail(t, "name can only be lowercase", name)
			}
			if _, found := seen[name]; found {
				assert.Fail(t, "duplicate name found", name)
			}
			seen[name] = struct{}{}
		})

		var filter *Template
		var e error
		t.Run("Parse/"+name, func(t *testing.T) {
			filter, e = parseTemplate(name, file)
			require.NoError(t, e, "Template did not parse OK")
			require.NoError(t, checkRedundantPresetValues(filter, data.Presets), "Found redundant preset values")
			require.NoError(t, parsePresets(filter, data.Presets), "Preset values did not parse OK")
			assert.NoError(t, validate.Struct(filter), "Template did no pass input validation")
		})
		require.NotNil(t, filter, "Template did not parse OK")

		t.Run("Contributors/"+name, func(t *testing.T) {
			for _, c := range filter.Contributors {
				_, found := contributors.Get(c)
				assert.Truef(t, found, "unknown contributor %s", c)
			}
			for _, c := range filter.Sponsors {
				_, found := contributors.Get(c)
				assert.Truef(t, found, "unknown sponsor %s", c)
			}
		})

		for i, tc := range filter.Tests {
			t.Run(fmt.Sprintf("Test/%s/%d", name, i), func(t *testing.T) {
				var buf strings.Builder
				ctx := make(map[string]interface{})
				for k, v := range tc.Params {
					ctx[k] = v
				}
				assert.NoError(t, repo.Render(&buf, &Instance{
					Template: filter.Name,
					Params:   ctx,
					TestMode: false,
				}))
				assert.Equal(t, tc.Output, buf.String())
			})
		}

		return nil
	})
	assert.NoError(t, err)
}

func TestRenderRawRules(t *testing.T) {
	repo, err := Load(testTemplates, testTemplates)
	require.NoError(t, err)
	tpl, err := repo.Get("raw-rules")
	require.NoError(t, err)

	cases := []struct {
		params   map[string]interface{}
		expected string
	}{{
		params:   nil,
		expected: "",
	}, {
		params: map[string]interface{}{
			"boolean_param":   false,
			"another_boolean": false,
		},
		expected: "",
	}, {
		params: map[string]interface{}{
			"boolean_param":   true,
			"another_boolean": false,
		},
		expected: "helloA\nhelloB\n",
	}, {
		params: map[string]interface{}{
			"boolean_param":   false,
			"another_boolean": true,
		},
		expected: "helloC\n",
	}, {
		params: map[string]interface{}{
			"boolean_param":   true,
			"another_boolean": true,
		},
		expected: "helloA\nhelloB\nhelloC\n",
	}}

	for i, tc := range cases {
		t.Run(strconv.Itoa(i), func(t *testing.T) {
			buf := bytes.Buffer{}
			assert.NoError(t, tpl.renderRawRules(&buf, tc.params))
			assert.Equal(t, tc.expected, buf.String())
		})
	}
}

func checkRedundantPresetValues(f *Template, presets fs.FS) error {
	for _, param := range f.Params {
		for _, preset := range param.Presets {
			if len(preset.Values) == 0 {
				continue
			}

			filename := fmt.Sprintf(presetFilePattern, f.Name, preset.Name)
			file, err := presets.Open(filename)
			if err == nil {
				_ = file.Close()
				return fmt.Errorf("preset %s has values in both a preset file and the YAML, remove one", preset.Name)
			} else if !errors.Is(err, fs.ErrNotExist) {
				return fmt.Errorf("unexpected error opening %s: %w", filename, err)
			}
		}
	}
	return nil
}
