package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/alecthomas/kong"
)

const (
	uodfSourcePrefix string = "https://github.com/quenhus/uBlock-Origin-dev-filter/blob/"
	uodfRawPrefix    string = "https://raw.githubusercontent.com/quenhus/uBlock-Origin-dev-filter/"
)

var templates = map[string]func(file *filterAndDescription) error{
	"search-results": updateSearchResults,
}

type PresetsCmd struct{}

func (c *PresetsCmd) Run(k *kong.Context) error {
	for name, f := range templates {
		k.Printf("Updating %s...", name)
		filename := filterFolder + name + filterExtension
		file, err := os.Open(filename)
		k.FatalIfErrorf(err, "cannot open target file for reading")
		filter, err := read(file)
		k.FatalIfErrorf(err, "cannot decode target file")
		k.FatalIfErrorf(file.Close())

		k.FatalIfErrorf(f(filter), "cannot update presets")

		file, err = os.OpenFile(filename, os.O_WRONLY|os.O_TRUNC, 0)
		k.FatalIfErrorf(err, "cannot open target file for writing")
		k.FatalIfErrorf(filter.write(file), "cannot write result")
		k.FatalIfErrorf(file.Close())
		k.Printf("OK")
	}
	return nil
}

func updateSearchResults(filter *filterAndDescription) error {
	for i, param := range filter.filter.Params {
		for j, preset := range param.Presets {
			if strings.HasPrefix(preset.Source, uodfSourcePrefix) {
				values, err := fetchUodf(preset.Source)
				if err != nil {
					return fmt.Errorf("error fetching %s: %w", preset.Name, err)
				}
				filter.filter.Params[i].Presets[j].Values = values
			}
		}
	}
	return nil
}

func fetchUodf(url string) ([]string, error) {
	file := uodfRawPrefix + strings.TrimPrefix(url, uodfSourcePrefix)
	fmt.Println("downloading", file)
	res, err := http.Get(file)
	if err != nil {
		return nil, err
	}

	var values []string
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "!") || line == "" {
			continue
		}
		line = strings.TrimPrefix(line, "*://*")
		line = strings.TrimPrefix(line, "*://")
		line = strings.TrimSuffix(line, "*")

		values = append(values, line)
	}

	sort.Slice(values, func(i, j int) bool {
		return strings.Compare(values[i], values[j]) < 0
	})
	return values, scanner.Err()
}
