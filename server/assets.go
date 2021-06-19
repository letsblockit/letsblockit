package server

import (
	"embed"
	"io/fs"
)

//go:embed assets
var assetFiles embed.FS

func openAssets() fs.FS {
	return &wrappedAssets{
		root:  assetFiles,
		isDir: make(map[string]bool),
	}
}

// wrappedAssets wraps the fs.Open call to forbid access to folders
type wrappedAssets struct {
	root  fs.FS
	isDir map[string]bool
}

func (w wrappedAssets) Open(name string) (fs.File, error) {
	isDir, found := w.isDir[name]
	if isDir {
		return nil, fs.ErrNotExist
	}
	file, err := w.root.Open(name)
	if err != nil {
		return nil, err
	}
	if !found {
		stat, err := file.Stat()
		if err != nil {
			return nil, err
		}
		w.isDir[name] = stat.IsDir()
		if stat.IsDir() {
			return nil, fs.ErrNotExist
		}
	}
	return file, nil
}
