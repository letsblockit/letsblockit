package main

import (
	"encoding/json"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/alecthomas/kong"
)

const (
	hashOutputFile = "data/asset-hashes.json"
)

type hashAssetsCmd struct {
	All bool `help:"Update assets-hashes.json."`
}

func (c *hashAssetsCmd) Run(k *kong.Context) error {
	hashes := make(map[string]string)
	cwd, err := os.Getwd()
	k.FatalIfErrorf(err)
	cmd := exec.Command("git", "ls-files", ".", "--abbrev", "--format", "%(objectname) %(path)")
	cmd.Dir = path.Join(cwd, "data/assets")

	output, err := cmd.Output()
	k.FatalIfErrorf(err)

	entries := strings.Split(string(output), "\n")
	for _, entry := range entries {
		parts := strings.Split(entry, " ")
		if len(parts) < 2 {
			continue
		}
		hashes[parts[1]] = parts[0]
		hashes[strings.TrimSuffix(parts[1], ".gz")] = parts[0]
	}

	outputFile, err := os.OpenFile(hashOutputFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	k.FatalIfErrorf(err)
	enc := json.NewEncoder(outputFile)
	enc.SetIndent("", "    ")
	k.FatalIfErrorf(enc.Encode(&hashes))
	return outputFile.Close()
}
