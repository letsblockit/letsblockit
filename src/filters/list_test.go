package filters

import (
	"embed"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/letsblockit/letsblockit/src/filters/mocks"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"gopkg.in/yaml.v3"
)

//go:embed testdata/templates
var testTemplates embed.FS

//go:embed testdata/list.yaml
var testList []byte

type ListTestSuite struct {
	suite.Suite
	repository *Repository
	logger     *mocks.Mocklogger
	expectL    *mocks.MockloggerMockRecorder
}

func (s *ListTestSuite) SetupTest() {
	c := gomock.NewController(s.T())
	s.logger = mocks.NewMocklogger(c)
	s.expectL = s.logger.EXPECT()

	var err error
	s.repository, err = Load(testTemplates, testTemplates)
	s.NoError(err)
}

func (s *ListTestSuite) TestRenderEmpty() {
	buf := &strings.Builder{}
	list := &List{}
	s.NoError(list.Render(buf, s.logger, s.repository))
	s.Equal(`! Title: letsblock.it - 
! Expires: 12 hours
! Homepage: https://letsblock.it
! License: https://github.com/letsblockit/letsblockit/blob/main/LICENSE.txt
`, buf.String())
}

func (s *ListTestSuite) TestRenderOK() {
	var list List
	require.NoError(s.T(), yaml.Unmarshal(testList, &list))

	buf := &strings.Builder{}

	s.expectL.Warnf(gomock.Any(), "unknown", gomock.Any())
	s.NoError(list.Render(buf, s.logger, s.repository))
	s.Equal(`! Title: letsblock.it - Test list
! Expires: 12 hours
! Homepage: https://letsblock.it
! License: https://github.com/letsblockit/letsblockit/blob/main/LICENSE.txt

! hello
Hello
! hello
Hello
! unknown

! simple
one
two
`, buf.String())
}

func (s *ListTestSuite) TestValidateOK() {
	list := &List{
		Title: "Test list",
		Instances: []*Instance{{
			Template: "hello",
		}, {
			Template: "simple",
			Params: map[string]interface{}{
				"string_list": []string{"one", "two"},
			},
		}},
	}
	s.NoError(list.Validate())
}

func TestListTestSuite(t *testing.T) {
	suite.Run(t, new(ListTestSuite))
}
