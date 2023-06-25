package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strings"

	"github.com/alecthomas/kong"
	"github.com/labstack/gommon/log"
	"github.com/letsblockit/letsblockit/data"
	"github.com/letsblockit/letsblockit/src/filters"
)

const (
	presetFilePattern string = "data/filters/presets/%s/%s.txt"
	uodfSourcePrefix  string = "https://github.com/quenhus/uBlock-Origin-dev-filter/blob/"
	bmrSource         string = "https://github.com/rotgruengelb/BlockModReposting/blob/main/list.txt"
	githubRawPrefix   string = "https://raw.githubusercontent.com"
)

var (
	githubPathPattern = regexp.MustCompile(`https://github.com/([-\w]+)/([-\w]+)/blob/(.*)`)
	targets           = map[string]func(file *filters.Template) error{
		"search-results": updateSearchResults,
	}
)

type updatePresetsCmd struct{}

func (c *updatePresetsCmd) Run(k *kong.Context) error {
	repo, err := filters.Load(data.Templates, data.Presets)
	if err != nil {
		log.Fatal(err)
	}
	for name, f := range targets {
		k.Printf("Updating presets for %s...", name)
		template, err := repo.Get(name)
		k.FatalIfErrorf(err, "unknown template %s", name)
		k.FatalIfErrorf(f(template), "failed to process %s", name)
	}
	return nil
}

func updateSearchResults(template *filters.Template) error {
	for _, param := range template.Params {
		for _, preset := range param.Presets {
			switch {
			case strings.HasPrefix(preset.Source, uodfSourcePrefix):
				values, err := fetchUodf(preset.Source)
				if err != nil {
					return fmt.Errorf("error fetching %s: %w", preset.Name, err)
				}
				if err = saveValues(template.Name, preset.Name, values); err != nil {
					return err
				}
			case preset.Source == bmrSource:
				values, err := fetchNetworkRules(preset.Source)
				if err != nil {
					return fmt.Errorf("error fetching %s: %w", preset.Name, err)
				}
				if err = saveValues(template.Name, preset.Name, values); err != nil {
					return err
				}
			default:
				return fmt.Errorf("error fetching %s: %s", preset.Name, "unknown source format")
			}
		}
	}
	return nil
}

func fetchUodf(url string) ([]string, error) {
	file := buildGithubRawUrl(url)
	fmt.Println("  downloading", file)
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

func fetchNetworkRules(url string) ([]string, error) {
	file := buildGithubRawUrl(url)
	fmt.Println("  downloading", file)
	res, err := http.Get(file)
	if err != nil {
		return nil, err
	}

	ruleRe := regexp.MustCompile(`^\|\|(.*)\^\$all$`)
	var values []string
	scanner := bufio.NewScanner(res.Body)
	for scanner.Scan() {
		if match := ruleRe.FindStringSubmatch(scanner.Text()); len(match) == 2 {
			values = append(values, match[1]+"/")
		}
	}

	sort.Slice(values, func(i, j int) bool {
		return strings.Compare(values[i], values[j]) < 0
	})
	return values, scanner.Err()
}

func buildGithubRawUrl(url string) string {
	if parts := githubPathPattern.FindStringSubmatch(url); len(parts) == 4 {
		parts[0] = githubRawPrefix
		return strings.Join(parts, "/")
	}
	return url
}

func saveValues(template, preset string, values []string) error {
	fileName := fmt.Sprintf(presetFilePattern, template, preset)
	file, err := os.OpenFile(fileName, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	fmt.Println("  writing", fileName)
	for _, v := range values {
		if _, err = fmt.Fprintln(file, v); err != nil {
			return err
		}
	}
	return nil
}
