package main

import (
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/letsblockit/letsblockit/data"
)

const (
	inputFileName    string = "src/assets/node_modules/@tabler/icons/tabler-sprite-nostroke.svg"
	outputFileName   string = "data/tabler-icons.yaml"
	outputFileHeader string = `# Extract of the SVG paths from tabler-icons
# Source: https://github.com/tabler/tabler-icons
# License: MIT`
)

var (
	iconRegex   = regexp.MustCompile(`{{\s?>icon\s+name="(.+?)"`)
	symbolRegex = regexp.MustCompile(`<symbol id="tabler-(.+?)".+?>(<path stroke="none" d="M0 0h24v24H0z" fill="none"/>)?(.+?)</symbol>`)
	extraIcons  = map[string][]string{
		"custom-open-collective": {
			`<path fill="currentColor" stroke=none d="M20.22 6.305c2.365 3.312 2.365 8.078 0 11.39l-2.59-2.59c1.063-1.89 1.063-4.32 0-6.21l2.59-2.59z"/>`,
			`<path fill="currentColor" stroke=none d="m17.695 3.78-2.59 2.59c-2.606-1.505-6.198-.815-8.066 1.54-1.985 2.312-1.94 6.035.1 8.297 1.89 2.267 5.406 2.898 7.966 1.423l2.59 2.59c-3.344 2.392-8.168 2.358-11.484-.065-3.119-2.162-4.768-6.197-4.055-9.925.662-3.983 3.977-7.339 7.951-8.051 2.61-.508 5.408.07 7.588 1.6z"/>`,
		},
	}
)

type extractIconsCmd struct {
	All bool `help:"extract all icons from the source."`
}

func (c *extractIconsCmd) Run(k *kong.Context) error {
	neededIcons := make(map[string]bool)

	if !c.All {
		k.FatalIfErrorf(data.Walk(data.Pages, ".hbs", func(name string, file io.Reader) error {
			source, err := io.ReadAll(file)
			if err != nil {
				return nil
			}
			for _, match := range iconRegex.FindAllSubmatch(source, -1) {
				neededIcons[string(match[1])] = false
			}
			return nil
		}))
		if len(neededIcons) == 0 {
			k.Fatalf("found no icon in the page templates")
		}
		k.Printf("found %d icons to extract", len(neededIcons))
	} else {
		k.Printf("extracting all icons in input file")
	}

	input, err := os.ReadFile(inputFileName)
	k.FatalIfErrorf(err)
	symbols := symbolRegex.FindAllSubmatch(input, -1)
	if len(symbols) == 0 {
		k.Fatalf("no icon was found in input file")
	}

	outputFile, err := os.OpenFile(outputFileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	k.FatalIfErrorf(err)
	_, err = fmt.Fprintln(outputFile, outputFileHeader)
	k.FatalIfErrorf(err)

	for _, match := range symbols {
		if !c.All {
			name := string(match[1])
			if found, needed := neededIcons[name]; !needed || found {
				continue
			}
			neededIcons[name] = true
		}
		_, err := fmt.Fprintf(outputFile, "%s: %s\n", match[1], match[3])
		k.FatalIfErrorf(err)
	}

	for name, lines := range extraIcons {
		if !c.All {
			if found, needed := neededIcons[name]; !needed || found {
				continue
			}
			neededIcons[name] = true
		}
		_, err := fmt.Fprintf(outputFile, "%s: %s\n", name, strings.Join(lines, ""))
		k.FatalIfErrorf(err)
	}

	if !c.All {
		for name, found := range neededIcons {
			if !found {
				k.Fatalf("failed to extract icon %s", name)
			}
		}
		k.Printf("extracted %d icons", len(neededIcons))
	}

	return outputFile.Close()
}
