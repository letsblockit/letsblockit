package main

import (
	"fmt"
	"io"
	"os"
	"regexp"

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
