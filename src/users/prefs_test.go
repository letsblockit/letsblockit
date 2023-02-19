package users

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/random"
	"github.com/letsblockit/letsblockit/src/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

const maxTimeSkew = 100 * time.Millisecond

func pastNow(hours int64) time.Time {
	return time.Now().Add(time.Duration(-1*hours) * time.Hour).Round(time.Second)
}

type PreferenceManagerSuite struct {
	suite.Suite
	store db.Store
	prefs *PreferenceManager
	user  string
	ctx   echo.Context
}

func (s *PreferenceManagerSuite) SetupTest() {
	s.store = db.NewTestStore(s.T())
	s.user = random.String(12)
	s.ctx = echo.New().NewContext(httptest.NewRequest(http.MethodGet, "/", nil), httptest.NewRecorder())

	var err error
	s.prefs, err = NewPreferenceManager(s.store)
	require.NoError(s.T(), err)
}

func (s *PreferenceManagerSuite) TestInitIfNotFound() {
	got, err := s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), s.user, got.UserID)
	assert.WithinDuration(s.T(), time.Now(), got.NewsCursor, maxTimeSkew)
}

func (s *PreferenceManagerSuite) TestGetCached() {
	expected := db.UserPreference{
		UserID:     s.user,
		NewsCursor: pastNow(10),
		ColorMode:  db.ColorModeAuto,
	}
	_, err := s.store.InitUserPreferences(context.Background(), s.user)
	require.NoError(s.T(), err)
	require.NoError(s.T(), s.store.UpdateNewsCursor(context.Background(), db.UpdateNewsCursorParams{
		UserID:     s.user,
		NewsCursor: expected.NewsCursor,
	}))

	got, err := s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &expected, got)

	// Out-of-band DB update will be ignored due to the cache
	require.NoError(s.T(), s.store.UpdateNewsCursor(context.Background(), db.UpdateNewsCursorParams{
		UserID:     s.user,
		NewsCursor: time.Now(),
	}))

	// Second get returns the cached value
	got, err = s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &expected, got)
}

func (s *PreferenceManagerSuite) TestUpdateNewsCursor() {
	initial := db.UserPreference{
		UserID:     s.user,
		NewsCursor: pastNow(10),
		ColorMode:  db.ColorModeAuto,
	}
	updated := db.UserPreference{
		UserID:     s.user,
		NewsCursor: pastNow(1),
		ColorMode:  db.ColorModeAuto,
	}
	_, err := s.store.InitUserPreferences(context.Background(), s.user)
	require.NoError(s.T(), err)
	require.NoError(s.T(), s.store.UpdateNewsCursor(context.Background(), db.UpdateNewsCursorParams{
		UserID:     s.user,
		NewsCursor: initial.NewsCursor,
	}))

	got, err := s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &initial, got)

	// Update the value through the manager
	assert.NoError(s.T(), s.prefs.UpdateNewsCursor(s.ctx, s.user, updated.NewsCursor))

	// Cache has been invalidated
	got, err = s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &updated, got)
}

func (s *PreferenceManagerSuite) TestUpdatePreferences() {
	initial := db.UserPreference{
		UserID:       s.user,
		NewsCursor:   pastNow(10),
		BetaFeatures: false,
		ColorMode:    db.ColorModeAuto,
	}
	withBeta := db.UserPreference{
		UserID:       s.user,
		NewsCursor:   pastNow(10),
		BetaFeatures: true,
		ColorMode:    db.ColorModeAuto,
	}
	withDark := db.UserPreference{
		UserID:       s.user,
		NewsCursor:   pastNow(10),
		BetaFeatures: false,
		ColorMode:    db.ColorModeDark,
	}

	_, err := s.store.InitUserPreferences(context.Background(), s.user)
	require.NoError(s.T(), err)
	require.NoError(s.T(), s.store.UpdateNewsCursor(context.Background(), db.UpdateNewsCursorParams{
		UserID:     s.user,
		NewsCursor: initial.NewsCursor,
	}))

	got, err := s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &initial, got)

	// Update the value through the manager
	assert.NoError(s.T(), s.prefs.UpdatePreferences(s.ctx, db.UpdateUserPreferencesParams{
		UserID:       s.user,
		ColorMode:    "auto",
		BetaFeatures: true,
	}))

	// Cache has been invalidated
	got, err = s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &withBeta, got)

	// Update the value through the manager
	assert.NoError(s.T(), s.prefs.UpdatePreferences(s.ctx, db.UpdateUserPreferencesParams{
		UserID:       s.user,
		ColorMode:    "dark",
		BetaFeatures: false,
	}))

	// Cache has been invalidated
	got, err = s.prefs.Get(s.ctx, s.user)
	assert.NoError(s.T(), err)
	assert.EqualValues(s.T(), &withDark, got)
}

func TestPreferenceManagerSuite(t *testing.T) {
	suite.Run(t, new(PreferenceManagerSuite))
}
