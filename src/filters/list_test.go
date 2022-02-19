package filters

import (
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/suite"
	"github.com/xvello/letsblockit/src/filters/mocks"
)

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
	s.repository, err = load(os.DirFS("testdata"))
	s.NoError(err)
}

func (s *ListTestSuite) TestRenderEmpty() {
	buf := &strings.Builder{}
	list := &List{}
	s.NoError(list.Render(buf, s.logger, s.repository))
	s.Equal(`! Title: letsblock.it - 
! Expires: 1 day
! Homepage: https://letsblock.it
! License: https://github.com/xvello/letsblockit/blob/main/LICENSE.txt
`, buf.String())
}

func (s *ListTestSuite) TestRenderOK() {
	buf := &strings.Builder{}
	list := &List{
		Title: "Test list",
		Instances: []*Instance{{
			Filter: "hello",
		}, {
			Filter: "hello",
		}, {
			Filter: "unknown",
		}, {
			Filter: "simple",
			Params: map[string]interface{}{
				"string_list": []string{"one", "two"},
			},
		}},
	}
	s.expectL.Warnf(gomock.Any(), "unknown", gomock.Any())
	s.NoError(list.Render(buf, s.logger, s.repository))
	s.Equal(`! Title: letsblock.it - Test list
! Expires: 1 day
! Homepage: https://letsblock.it
! License: https://github.com/xvello/letsblockit/blob/main/LICENSE.txt

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
			Filter: "hello",
		}, {
			Filter: "simple",
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
