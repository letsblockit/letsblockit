package server

import (
	"fmt"
	"strings"
	"testing"

	"github.com/aymerick/raymond"
	"github.com/stretchr/testify/assert"
)

type mockedEcho struct{}

func (m *mockedEcho) Reverse(name string, params ...interface{}) string {
	var builder strings.Builder
	builder.WriteString(fmt.Sprintf("//%s", name))
	for _, p := range params {
		builder.WriteString(fmt.Sprintf("//%s", p))
	}
	return builder.String()
}

func TestHelpers(t *testing.T) {
	tests := map[string]struct {
		input    string
		ctx      map[string]interface{}
		expected string
	}{
		"href_noarg": {
			input:    `{{href "name" ""}}`,
			expected: "//name",
		},
		"href_one_arg": {
			input:    `{{href "name" "one"}}`,
			expected: "//name//one",
		},
		"href_three_arg": {
			input:    `{{href "name" "one/two/three"}}`,
			expected: "//name//one//two//three",
		},
		"lookup_list_strings": {
			input: `-{{#each (lookup_list params "name")}}{{this}}-{{/each}}`,
			ctx: map[string]interface{}{
				"params": map[string]interface{}{
					"simple": "single",
					"name":   []string{"one", "two"},
				},
			},
			expected: "-one-two-",
		},
		"lookup_list_interface": {
			input: `-{{#each (lookup_list params "name")}}{{this}}-{{/each}}`,
			ctx: map[string]interface{}{
				"params": map[string]interface{}{
					"simple": "single",
					"name":   []interface{}{"one", 2},
				},
			},
			expected: "-one-2-",
		},
	}
	helpers := buildHelpers(&mockedEcho{})

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tpl := raymond.MustParse(tc.input)
			tpl.RegisterHelpers(helpers)
			assert.Equal(t, tc.expected, tpl.MustExec(tc.ctx))
		})
	}
}
