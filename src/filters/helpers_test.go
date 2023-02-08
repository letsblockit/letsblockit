package filters

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBuildRegexForWords(t *testing.T) {
	cases := map[string]string{
		`word`:      `/\bword\b/i`,
		`/regexp/`:  `/regexp/`,
		`/regexp/i`: `/regexp/i`,
	}

	for in, out := range cases {
		t.Run("buildRegexForWord:"+in, func(t *testing.T) {
			assert.Equal(t, out, buildRegexForWord(in))
		})
	}

	inStrings, inInterfaces := make([]string, 0, len(cases)), make([]interface{}, 0, len(cases))
	outStrings := make([]string, 0, len(cases))
	for in, out := range cases {
		inStrings = append(inStrings, in)
		inInterfaces = append(inInterfaces, interface{}(in))
		outStrings = append(outStrings, out)
	}
	t.Run("buildRegexForWords", func(t *testing.T) {
		assert.Equal(t, outStrings, buildRegexForWords(inStrings))
		assert.Equal(t, outStrings, buildRegexForWords(inInterfaces))
		assert.Empty(t, buildRegexForWords([]int{1, 2, 3}))
	})
}
