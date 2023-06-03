package news

import (
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"io/fs"
	"os"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-retryablehttp"
	"github.com/russross/blackfriday/v2"
	"github.com/samber/lo"
)

const (
	GithubReleasesEndpoint string = "https://api.github.com/repos/xvello/letsblockit/releases?per_page=20"
	cacheFileName          string = "lbi-releases.json"
)

var (
	githubUserRegex = regexp.MustCompile(`(\W)@([0-9A-Za-z-]+)(\W)`)
	githubUserLink  = `$1[**@$2**](https://github.com/$2)$3`
)

type templateExists = func(name string) bool

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

// Releases fetches and parses the github releases for a repository.
// The parsed results are cached in memory until the next restart.
type Releases struct {
	latestAt time.Time
	releases []*Release
	etag     string
}

func (c *Releases) GetReleases() ([]*Release, string) {
	return c.releases, c.etag
}

func (c *Releases) GetLatestAt() time.Time {
	return c.latestAt
}

func DownloadReleases(url string, cacheDir string, officialInstance bool, templates fs.ReadDirFS) (*Releases, error) {
	contents, err := download(url, cacheDir, cacheFileName)
	if err != nil {
		return nil, err
	}
	defer func() { _ = contents.Close() }()

	var githubReleases []githubRelease
	if err = json.NewDecoder(contents).Decode(&githubReleases); err != nil {
		return nil, err
	}

	availableTemplates, err := enumerateTemplates(templates)
	if err != nil {
		return nil, err
	}
	etagHasher := fnv.New64()
	renderer := initRenderer(officialInstance, availableTemplates)
	output := &Releases{
		releases: make([]*Release, 0, len(githubReleases)),
	}
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
		output.releases = append(output.releases, &Release{
			Id:          r.Id,
			Link:        r.HtmlUrl,
			Description: string(desc),
			CreatedAt:   r.CreatedAt,
			PublishedAt: r.PublishedAt,
			TagName:     r.TagName,
			GithubUrl:   r.HtmlUrl,
		})
		if r.CreatedAt.After(output.latestAt) {
			output.latestAt = r.CreatedAt
		}
	}
	output.etag = strconv.FormatUint(etagHasher.Sum64(), 36)
	return output, nil
}

func enumerateTemplates(templates fs.ReadDirFS) (templateExists, error) {
	entries, err := templates.ReadDir(".")
	if err != nil {
		return nil, err
	}

	present := lo.SliceToMap(entries, func(item fs.DirEntry) (string, bool) {
		return strings.TrimSuffix(item.Name(), ".yaml"), true
	})
	return func(name string) bool { return present[name] }, nil
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
