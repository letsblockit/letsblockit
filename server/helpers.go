package server

import (
	"fmt"
	"strings"
)

type echoInterface interface {
	Reverse(name string, params ...interface{}) string
}

func buildHelpers(e echoInterface) map[string]interface{} {
	assetHash := computeAssetHash()
	if assetHash != "" {
		assetHash = "?h=" + assetHash
	}

	return map[string]interface{}{
		"assetHash": func() string {
			return assetHash
		},
		"href": func(route string, args string) string {
			return href(e, route, args)
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
	}
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
