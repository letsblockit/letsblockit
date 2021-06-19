package server

import (
	"strings"
)

type echoInterface interface {
	Reverse(name string, params ...interface{}) string
}

func buildHelpers(e echoInterface) map[string]interface{} {
	return map[string]interface{}{
		"href": func(route string, args string) string {
			return href(e, route, args)
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
