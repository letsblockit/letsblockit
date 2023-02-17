package data

import (
	"embed"
	"fmt"
	"hash/fnv"
	"io"
	"io/fs"
	"strconv"
	"strings"
)

//go:embed assets
var Assets embed.FS

//go:embed filters/templates
var Templates embed.FS

//go:embed filters/presets
var Presets embed.FS

//go:embed pages/*
var Pages embed.FS

// Walk warps fs.WalkDir with simpler invocation pattern:
//   - only files with a given suffix are passed opened
//   - the file is automatically opened and closed
//   - only the shortened file name (no folder, no suffix) and io.Reader are passed down
func Walk(input fs.FS, suffix string, fn func(string, io.Reader) error) error {
	return fs.WalkDir(input, ".", func(path string, d fs.DirEntry, _ error) error {
		if d.IsDir() || !strings.HasSuffix(d.Name(), suffix) {
			return nil
		}
		name := strings.TrimSuffix(d.Name(), suffix)
		file, err := input.Open(path)
		if err != nil {
			return fmt.Errorf("cannot open %s: %w", path, err)
		}
		defer file.Close()
		if err = fn(name, file); err != nil {
			return fmt.Errorf("cannot process %s: %w", path, err)
		}
		return nil
	})
}

// HashFiles computes a fnv hash of the given filesystems, for cache invalidation purposes
func HashFiles(input ...fs.FS) (string, error) {
	hasher := fnv.New64()
	for _, i := range input {
		if err := Walk(i, "", func(s string, reader io.Reader) error {
			_, err := io.Copy(hasher, reader)
			return err
		}); err != nil {
			return "", err
		}
	}
	return strconv.FormatUint(hasher.Sum64(), 36), nil
}
