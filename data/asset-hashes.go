package data

import (
	_ "embed"
	"encoding/json"
	"fmt"
)

//go:embed asset-hashes.json
var assetHashesFile []byte

type AssetHashes struct {
	hashes map[string]string
}

func ParseAssetHashes() (*AssetHashes, error) {
	output := &AssetHashes{
		hashes: make(map[string]string),
	}
	return output, json.Unmarshal(assetHashesFile, &output.hashes)
}

func (h *AssetHashes) BuildURL(path string) string {
	if hash, found := h.hashes[path]; found {
		return fmt.Sprintf("/assets/%s?h=%s", path, hash)
	}
	return fmt.Sprintf("/assets/%s", path)
}
