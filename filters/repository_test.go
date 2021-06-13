package filters

import (
	"os"
	"testing"

	"github.com/aymerick/raymond"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

func TestLoadFilters(t *testing.T) {
	repo, err := LoadFilters()
	require.NoError(t, err, "Filter parsing error")
	require.NotNil(t, repo, "Filter repository is empty")
	require.Greater(t, len(repo.filters), 0, "Expected at least one filter")
}

type SimpleFilterSuite struct {
	suite.Suite
	repository *Repository
	filterName string
}

func (s *SimpleFilterSuite) SetupTest() {
	var err error
	s.filterName = "____simple____filter____"
	s.repository, err = LoadFilters()
	require.NoError(s.T(), err)
	file, err := os.Open("testdata/simple.yaml")
	require.NoError(s.T(), err)
	defer file.Close()
	s.repository.filters[s.filterName], err = parseFilter(s.filterName, file)
	require.NoError(s.T(), err)
}

func (s *SimpleFilterSuite) TestRenderFilter() {
	params := map[string][]string{
		"string_list": {"one line", "another line"},
	}
	out, err := s.repository.RenderFilter(s.filterName, params)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), "one line\nanother line\n", out)
}

func (s *SimpleFilterSuite) TestRenderPage() {
	template := raymond.MustParse(`
{{#each params}}
{{ name }}: {{ description }}
{{/each}}
`)

	out, err := s.repository.RenderPage(s.filterName, template)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), `
boolean_param: A boolean parameter
another_boolean: A disabled boolean parameter
string_param: A string parameter
string_list: A list of strings
`, out)
}

func TestSimpleFilterSuite(t *testing.T) {
	suite.Run(t, new(SimpleFilterSuite))
}
