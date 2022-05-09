package users

import (
	"context"
	"time"

	"github.com/google/uuid"
	lru "github.com/hashicorp/golang-lru"
	"github.com/labstack/echo/v4"
	"github.com/xvello/letsblockit/src/db"
)

const prefCacheSize = 2048

type PreferenceManager struct {
	cache *lru.Cache
	store db.Store
}

func NewPreferenceManager(store db.Store) (*PreferenceManager, error) {
	cache, err := lru.New(prefCacheSize)
	if err != nil {
		return nil, err
	}
	return &PreferenceManager{
		cache: cache,
		store: store,
	}, nil
}

// Get retrieves the preferences from cache or DB. If no prefs are in DB, create a row with default values.
func (m *PreferenceManager) Get(c echo.Context, user uuid.UUID) (*db.UserPreference, error) {
	if entry, ok := m.cache.Get(user); ok {
		if prefs, ok := entry.(*db.UserPreference); ok {
			return prefs, nil
		}
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
	m.cache.Add(user, &prefs)
	return &prefs, nil
}

func (m *PreferenceManager) UpdateNewsCursor(c echo.Context, user uuid.UUID, at time.Time) error {
	err := m.store.UpdateNewsCursor(c.Request().Context(), db.UpdateNewsCursorParams{
		UserID:     user,
		NewsCursor: at,
	})
	m.cache.Remove(user)
	return err
}
