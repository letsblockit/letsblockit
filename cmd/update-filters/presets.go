package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"sort"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/letsblockit/letsblockit/data"
	"github.com/letsblockit/letsblockit/src/filters"
)

const (
	templateFilePattern string = "data/filters/templates/%s.yaml"
	presetFilePattern   string = "data/filters/presets/%s/%s.txt"

	uodfSourcePrefix string = "https://github.com/quenhus/uBlock-Origin-dev-filter/blob/"
	uodfRawPrefix    string = "https://raw.githubusercontent.com/quenhus/uBlock-Origin-dev-filter/"
)

var targets = map[string]func(file *filters.Template) error{
	"search-results": updateSearchResults,
}

type PresetsCmd struct{}

func (c *PresetsCmd) Run(k *kong.Context) error {
	repo, err := filters.Load(data.Templates, data.Presets)
	k.FatalIfErrorf(err)
	for name, f := range targets {
		k.Printf("Updating %s...", name)
		template, err := repo.Get(name)
		k.FatalIfErrorf(err, "unknown template")
		k.FatalIfErrorf(f(template), "cannot update presets")
		k.Printf("OK")
	}
	return nil
}

func updateSearchResults(template *filters.Template) error {
	for _, param := range template.Params {
		for _, preset := range param.Presets {
			if strings.HasPrefix(preset.Source, uodfSourcePrefix) {
				values, err := fetchUodf(preset.Source)
				if err != nil {
					return fmt.Errorf("error fetching %s: %w", preset.Name, err)
				}
				if err = savePreset(template.Name, preset.Name, values); err != nil {
					return err
				}
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

func savePreset(template, preset string, values []string) error {
	fileName := fmt.Sprintf(presetFilePattern, template, preset)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	fmt.Println("writing", fileName)
	for _, v := range values {
		if _, err = fmt.Fprintln(file, v); err != nil {
			return err
		}
	}
	return nil
}
