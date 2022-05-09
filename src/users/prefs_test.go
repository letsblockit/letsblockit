package users

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/xvello/letsblockit/src/db"
	"github.com/xvello/letsblockit/src/server/mocks"
)

type PreferenceManagerSuite struct {
	suite.Suite
	expectQ *mocks.MockQuerierMockRecorder
	prefs   *PreferenceManager
	user    uuid.UUID
	ctx     echo.Context
}

func (s *PreferenceManagerSuite) SetupTest() {
	c := gomock.NewController(s.T())
	qm := mocks.NewMockQuerier(c)
	s.expectQ = qm.EXPECT()
	s.user = uuid.New()
	s.ctx = echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

	var err error
	s.prefs, err = NewPreferenceManager(mocks.NewMockStore(qm))
	require.NoError(s.T(), err)
}

func (s *PreferenceManagerSuite) TestInitIfNotFound() {
	expected := db.UserPreference{
		UserID:     s.user,
		NewsCursor: time.Now(),
	}
	s.expectQ.GetUserPreferences(gomock.Any(), s.user).Return(db.UserPreference{}, db.NotFound)
	s.expectQ.InitUserPreferences(gomock.Any(), s.user).Return(expected, nil)

	got, err := s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &expected, got)
}

func (s *PreferenceManagerSuite) TestGetCached() {
	expected := db.UserPreference{
		UserID:     s.user,
		NewsCursor: time.Now(),
	}
	s.expectQ.GetUserPreferences(gomock.Any(), s.user).Return(expected, nil).MaxTimes(1)
	got, err := s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &expected, got)

	// Second get hits the cache
	got, err = s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &expected, got)
}

func (s *PreferenceManagerSuite) TestUpdateNewsCursor() {
	expected := db.UserPreference{
		UserID:     s.user,
		NewsCursor: time.Now(),
	}
	updated := db.UserPreference{
		UserID:     s.user,
		NewsCursor: expected.NewsCursor.Add(time.Hour),
	}
	s.expectQ.GetUserPreferences(gomock.Any(), s.user).Return(expected, nil).MaxTimes(1)
	got, err := s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &expected, got)

	// Update the value
	s.expectQ.UpdateNewsCursor(gomock.Any(), db.UpdateNewsCursorParams{
		UserID:     s.user,
		NewsCursor: updated.NewsCursor,
	})
	assert.NoError(s.T(), s.prefs.UpdateNewsCursor(s.ctx, s.user, updated.NewsCursor))

	// Cache has been invalidated
	s.expectQ.GetUserPreferences(gomock.Any(), s.user).Return(updated, nil).MaxTimes(1)
	got, err = s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &updated, got)
}

func TestPreferenceManagerSuite(t *testing.T) {
	suite.Run(t, new(PreferenceManagerSuite))
}
