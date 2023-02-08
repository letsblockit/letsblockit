package filters

import (
	"fmt"
	"regexp"

	"github.com/samber/lo"
)

var maybeRegExpPattern = regexp.MustCompile(`^\/.+\/[dgimsuy]*$`)

func buildRegexForWords(values interface{}) []string {
	switch values := values.(type) {
	case []string:
		fmt.Println("called with []string")
		return lo.Map(values, func(item string, _ int) string {
			return buildRegexForWord(item)
		})
	case []interface{}:
		fmt.Println("called with []interface{}")
		return lo.Map(values, func(item interface{}, _ int) string {
			return buildRegexForWord(fmt.Sprint(item))
		})
	default:
		return nil
	}
}

func buildRegexForWord(value string) string {
	if maybeRegExpPattern.MatchString(value) {
		return value
	}
	return fmt.Sprintf(`/\b%s\b/i`, value)
}
