package server

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/letsblockit/letsblockit/data"
	"github.com/letsblockit/letsblockit/src/db"
	"github.com/letsblockit/letsblockit/src/filters"
	"github.com/letsblockit/letsblockit/src/pages"
)

type echoInterface interface {
	Reverse(name string, params ...interface{}) string
}

func buildHelpers(e echoInterface) (map[string]interface{}, error) {
	contributors, err := data.ParseContributors()
	if err != nil {
		return nil, err
	}

	return map[string]interface{}{
		"eq": func(a string, b string) bool {
			return strings.Compare(a, b) == 0
		},
		"tag": func(name string) string {
			return fmt.Sprintf(
				`<a href="%s" class="badge rounded-pill bg-secondary text-decoration-none me-2">%s</a>`,
				href(e, "filters-for-tag", name), name)
		},
		"href": func(route string, args string) string {
			return href(e, route, args)
		},
		"abp_subscribe_href": func(listUrl string) string {
			return fmt.Sprintf("abp:subscribe?location=%s&title=%s",
				url.QueryEscape(listUrl),
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
		"preset_name": func(param filters.Parameter, preset filters.Preset) string {
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
		"beta_features": func(c *pages.Context) bool {
			return c.Preferences != nil && c.Preferences.BetaFeatures
		},
		"is_color_mode": func(c *pages.Context, mode string) bool {
			if c.Preferences == nil {
				return mode == string(db.ColorModeAuto)
			}
			return mode == string(c.Preferences.ColorMode)
		},
		"avatars": func(names []string) []*data.Contributor {
			var output []*data.Contributor
			for _, name := range names {
				if c, found := contributors.Get(name); found {
					output = append(output, c)
				}
			}
			return output
		},
		"all_avatars": func() []*data.Contributor {
			return contributors.GetAll()
		},
		"all_sponsors": func() []*data.Contributor {
			return contributors.GetSponsors()
		},
	}, nil
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
