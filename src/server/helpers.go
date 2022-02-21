package server

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/xvello/letsblockit/src/filters"
	"github.com/xvello/letsblockit/src/pages"
)

type echoInterface interface {
	Reverse(name string, params ...interface{}) string
}

func buildHelpers(e echoInterface, assetHash string) map[string]interface{} {
	assetHashQuery := "?h=" + assetHash

	return map[string]interface{}{
		"eq": func(a string, b string) bool {
			return strings.Compare(a, b) == 0
		},
		"assetHash": func() string {
			return assetHashQuery
		},
		"tag": func(name string) string {
			return fmt.Sprintf(
				`<a href="%s" class="badge rounded-pill bg-secondary text-decoration-none me-2">%s</a>`,
				href(e, "filters-for-tag", name), name)
		},
		"href": func(route string, args string) string {
			return href(e, route, args)
		},
		"list_href": func(token string) string {
			return listDownloadRef(e, token)
		},
		"abp_subscribe_href": func(token string) string {
			return fmt.Sprintf("abp:subscribe?location=%s&title=%s",
				url.QueryEscape(listDownloadRef(e, token)),
				url.QueryEscape("letsblock.it - My filters"))
		},
		"lookup_list": func(obj map[string]interface{}, key string) []string {
			switch values := obj[key].(type) {
			case []string:
				return values
			case []interface{}:
				var converted []string
				for _, v := range values {
					converted = append(converted, fmt.Sprintf("%v", v))
				}
				return converted
			default:
				return nil
			}
		},
		"preset_name": func(param filters.FilterParam, preset filters.Preset) string {
			return param.BuildPresetParamName(preset.Name)
		},
		"csrf": func(c *pages.Context) string {
			return fmt.Sprintf(
				`<input type="hidden" name="%s" value="%s"/>`,
				csrfLookup, c.CSRFToken)
		},
		"http_root": func(c *pages.Context) string {
			return fmt.Sprintf("%s://%s", c.RequestInfo.Scheme(), c.RequestInfo.Request().Host)
		},
	}
}

func listDownloadRef(e echoInterface, token string) string {
	return "https://get.letsblock.it" + href(e, "render-filterlist", token)
}

func href(e echoInterface, route string, args string) string {
	if len(args) == 0 {
		return e.Reverse(route) // Quick path if no args
	}
	if !strings.Contains(args, "/") {
		return e.Reverse(route, args) // Quick path for a single arg
	}
	// Multi-arg slow path
	argList := strings.Split(args, "/")
	outList := make([]interface{}, len(argList))
	for i, a := range argList {
		outList[i] = a
	}
	return e.Reverse(route, outList...)
}
