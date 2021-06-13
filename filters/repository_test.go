package filters

import (
	"embed"
	"testing"

	"github.com/aymerick/raymond"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

//go:embed testdata
var testDefinitionFiles embed.FS

// Check that all filter definitions parse OK
func TestLoadFilters(t *testing.T) {
	repo, err := LoadFilters()
	require.NoError(t, err, "Filter parsing error")
	require.NotNil(t, repo, "Filter repository is empty")
	require.Greater(t, len(repo.fMap), 0, "Expected at least one filter")
	require.Greater(t, len(repo.fList), 0, "Expected at least one filter")
}

type SimpleFilterSuite struct {
	suite.Suite
	repository *Repository
	filterName string
}

func (s *SimpleFilterSuite) SetupTest() {
	var err error
	s.filterName = "simple"
	s.repository, err = load(testDefinitionFiles)
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

func (s *SimpleFilterSuite) TestRenderIndex() {
	template := raymond.MustParse(`
{{#each this}}
{{ name }}: {{ title }}
{{/each}}
`)

	out, err := s.repository.RenderIndex(template)
	require.NoError(s.T(), err)
	assert.Equal(s.T(), `
hello: Hello filter
simple: Filter title
`, out)
}

func TestSimpleFilterSuite(t *testing.T) {
	suite.Run(t, new(SimpleFilterSuite))
}
