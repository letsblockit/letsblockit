package data

import (
	_ "embed"
	"encoding/json"
	"fmt"

	"github.com/samber/lo"
)

//go:embed contributors.json
var contributorsFile []byte

type Contributors struct {
	all      []*Contributor
	sponsors []*Contributor
	byLogin  map[string]*Contributor
}

type Contributor struct {
	Asset         string
	Login         string
	Name          string
	AvatarUrl     string `json:"avatar_url"`
	Profile       string
	Contributions []string
}

type contributorList struct {
	Contributors []*Contributor
}

func ParseContributors() (*Contributors, error) {
	var list contributorList
	if err := json.Unmarshal(contributorsFile, &list); err != nil {
		return nil, err
	}
	output := &Contributors{
		all:     list.Contributors,
		byLogin: make(map[string]*Contributor),
	}
	for _, c := range list.Contributors {
		c.Asset = fmt.Sprintf("/assets/images/contributors/%s.png", c.Login)
		output.byLogin[c.Login] = c
		if lo.Contains(c.Contributions, "financial") {
			output.sponsors = append(output.sponsors, c)
		}
	}

	return output, nil
}

func (c *Contributors) Get(login string) (*Contributor, bool) {
	item, found := c.byLogin[login]
	return item, found
}

func (c *Contributors) GetAll() []*Contributor {
	return c.all
}

func (c *Contributors) GetSponsors() []*Contributor {
	return c.sponsors
}
