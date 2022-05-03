package news

import (
	"encoding/json"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/russross/blackfriday/v2"
)

const (
	GithubReleasesEndpoint string = "https://api.github.com/repos/xvello/lbi-release-test/releases?per_page=100"
)

type githubRelease struct {
	HtmlUrl     string    `json:"html_url"`
	Id          int       `json:"id"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	CreatedAt   time.Time `json:"created_at"`
	PublishedAt time.Time `json:"published_at"`
	Body        string    `json:"body"`
}

type Release struct {
	Id          int
	Link        string
	Description string
	CreatedAt   time.Time
	PublishedAt time.Time
}

func (r Release) Date() string {
	return r.CreatedAt.Format("02 Jan. 2006")
}

// ReleaseClient fetches and parses the github releases for a repository.
// The parsed results are cached in memory until the next restart.
type ReleaseClient struct {
	sync.Mutex
	url      string
	latestAt time.Time
	releases []*Release
}

func NewReleaseClient(url string) *ReleaseClient {
	return &ReleaseClient{url: url}
}

func (c *ReleaseClient) GetReleases() ([]*Release, error) {
	c.Lock()
	defer c.Unlock()
	if c.releases != nil {
		return c.releases, nil
	}

	if err := c.populate(); err != nil {
		return nil, err
	}
	return c.releases, nil
}

func (c *ReleaseClient) GetLatestAt() (time.Time, error) {
	c.Lock()
	defer c.Unlock()
	if !c.latestAt.IsZero() {
		return c.latestAt, nil
	}

	if err := c.populate(); err != nil {
		return time.Time{}, err
	}
	return c.latestAt, nil
}

func (c *ReleaseClient) populate() error {
	resp, err := retryablehttp.Get(c.url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var githubReleases []githubRelease
	if err = json.NewDecoder(resp.Body).Decode(&githubReleases); err != nil {
		return err
	}

	renderer := initRenderer()
	c.releases = make([]*Release, 0, len(githubReleases))
	for _, r := range githubReleases {
		if r.Prerelease || r.Draft {
			continue
		}
		// Cleanup \r that mess up with blackfriday parsing
		body := strings.ReplaceAll(r.Body, "\r\n", "\n")
		desc := blackfriday.Run([]byte(body), renderer)
		c.releases = append(c.releases, &Release{
			Id:          r.Id,
			Link:        r.HtmlUrl,
			Description: string(desc),
			CreatedAt:   r.CreatedAt,
			PublishedAt: r.PublishedAt,
		})
		if r.CreatedAt.After(c.latestAt) {
			c.latestAt = r.CreatedAt
		}
	}
	return nil
}
