package data

import (
	_ "embed"
	"encoding/json"
)

//go:embed contributors.json
var contributorsFile []byte

type Contributor struct {
	Login         string
	Name          string
	AvatarUrl     string `json:"avatar_url"`
	Profile       string
	Contributions []string
}

type contributorList struct {
	Contributors []Contributor
}

func ParseContributors() ([]Contributor, error) {
	var list contributorList
	if err := json.Unmarshal(contributorsFile, &list); err != nil {
		return nil, err
	}
	return list.Contributors, nil
}
