package news

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/russross/blackfriday/v2"
)

const (
	GithubReleasesEndpoint string = "https://api.github.com/repos/xvello/letsblockit/releases?per_page=20"
	cacheFileName          string = "lbi-releases.json"
)

var (
	githubUserRegex = regexp.MustCompile(`(\W)@([0-9A-Za-z-]+)(\W)`)
	githubUserLink  = `$1[**@$2**](https://github.com/$2)$3`
)

type templateProvider interface {
	Has(name string) bool
}

type githubRelease struct {
	HtmlUrl     string    `json:"html_url"`
	Id          int       `json:"id"`
	Draft       bool      `json:"draft"`
	Prerelease  bool      `json:"prerelease"`
	TagName     string    `json:"tag_name"`
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
	TagName     string
	GithubUrl   string
}

func (r Release) Date() string {
	return r.CreatedAt.Format("2006-01-02")
}

// ReleaseClient fetches and parses the github releases for a repository.
// The parsed results are cached in memory until the next restart.
type ReleaseClient struct {
	sync.Mutex
	url              string
	cacheDir         string
	officialInstance bool
	templateProvider templateProvider
	latestAt         time.Time
	releases         []*Release
	etag             string
}

func NewReleaseClient(url string, cacheDir string, officialInstance bool, tp templateProvider) *ReleaseClient {
	return &ReleaseClient{
		url:              url,
		cacheDir:         cacheDir,
		officialInstance: officialInstance,
		templateProvider: tp,
	}
}

func (c *ReleaseClient) GetReleases() ([]*Release, string, error) {
	c.Lock()
	defer c.Unlock()
	if c.releases != nil {
		return c.releases, c.etag, nil
	}

	if err := c.populate(); err != nil {
		return nil, "", err
	}
	return c.releases, c.etag, nil
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
	contents, err := download(c.url, c.cacheDir, cacheFileName)
	if err != nil {
		return err
	}
	defer func() { _ = contents.Close() }()

	var githubReleases []githubRelease
	if err = json.NewDecoder(contents).Decode(&githubReleases); err != nil {
		return err
	}

	etagHasher := fnv.New64()
	renderer := initRenderer(c.officialInstance, c.templateProvider)
	c.releases = make([]*Release, 0, len(githubReleases))
	for _, r := range githubReleases {
		if r.Prerelease || r.Draft {
			continue
		}
		_, _ = etagHasher.Write([]byte(r.Body))
		// Cleanup \r that mess up with blackfriday parsing
		body := strings.ReplaceAll(r.Body, "\r\n", "\n")
		// Insert links for github users
		body = githubUserRegex.ReplaceAllString(body, githubUserLink)
		desc := blackfriday.Run([]byte(body), renderer)
		c.releases = append(c.releases, &Release{
			Id:          r.Id,
			Link:        r.HtmlUrl,
			Description: string(desc),
			CreatedAt:   r.CreatedAt,
			PublishedAt: r.PublishedAt,
			TagName:     r.TagName,
			GithubUrl:   r.HtmlUrl,
		})
		if r.CreatedAt.After(c.latestAt) {
			c.latestAt = r.CreatedAt
		}
	}
	c.etag = strconv.FormatUint(etagHasher.Sum64(), 36)
	return nil
}

func download(url string, cacheDir, cacheFileName string) (io.ReadCloser, error) {
	// Try opening the cache file
	if cacheDir != "" {
		if file, err := os.Open(path.Join(cacheDir, cacheFileName)); err == nil {
			return file, nil
		}
	}

	// Else, download it from the server
	resp, err := retryablehttp.Get(url)
	if err != nil {
		return nil, err
	}

	// If caching, copy the bytes to a file and pass that file as contents
	if cacheDir != "" {
		file, err := os.OpenFile(path.Join(cacheDir, cacheFileName), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0644)
		if err != nil {
			return nil, fmt.Errorf("cannot create changelog cache file: %w", err)
		}
		if _, err = io.Copy(file, resp.Body); err != nil {
			return nil, fmt.Errorf("cannot write changelog to cache: %w", err)
		}
		if err = resp.Body.Close(); err != nil {
			return nil, fmt.Errorf("error closing the http response: %w", err)
		}
		if _, err = file.Seek(0, io.SeekStart); err != nil {
			return nil, fmt.Errorf("cannot rewind cachelog cache file: %w", err)
		}
		return file, nil
	}

	// If not caching, return the response body directly
	return resp.Body, nil
}
