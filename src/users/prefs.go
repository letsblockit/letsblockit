package users

import (
	"context"
	"time"

	"zgo.at/zcache/v2"

	"github.com/labstack/echo/v4"
	"github.com/letsblockit/letsblockit/src/db"
)

type PreferenceManager struct {
	cache *zcache.Cache[string, *db.UserPreference]
	store db.Store
}

func NewPreferenceManager(store db.Store) (*PreferenceManager, error) {
	return &PreferenceManager{
		cache: zcache.New[string, *db.UserPreference](30*time.Minute, 10*time.Minute),
		store: store,
	}, nil
}

// Get retrieves the preferences from cache or DB. If no prefs are in DB, create a row with default values.
func (m *PreferenceManager) Get(c echo.Context, user string) (*db.UserPreference, error) {
	if entry, ok := m.cache.Get(user); ok {
		return entry, nil
	}
	var prefs db.UserPreference
	if err := m.store.RunTx(c, func(ctx context.Context, q db.Querier) error {
		var err error
		prefs, err = q.GetUserPreferences(ctx, user)
		if err == db.NotFound {
			prefs, err = q.InitUserPreferences(ctx, user)
		}
		return err
	}); err != nil {
		return nil, err
	}
	m.cache.Set(user, &prefs)
	return &prefs, nil
}

func (m *PreferenceManager) UpdateNewsCursor(c echo.Context, user string, at time.Time) error {
	if _, err := m.Get(c, user); err != nil {
		return err
	}
	err := m.store.UpdateNewsCursor(c.Request().Context(), db.UpdateNewsCursorParams{
		UserID:     user,
		NewsCursor: at,
	})
	m.cache.Delete(user)
	return err
}

func (m *PreferenceManager) UpdatePreferences(c echo.Context, params db.UpdateUserPreferencesParams) error {
	if _, err := m.Get(c, params.UserID); err != nil {
		return err
	}
	err := m.store.UpdateUserPreferences(c.Request().Context(), params)
	m.cache.Remove(params.UserID)
	return err
}
