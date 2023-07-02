package users

import (
	"context"
	"testing"

	"github.com/letsblockit/letsblockit/src/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadBannedUsers(t *testing.T) {
	store := db.NewTestStore(t)
	for _, user := range []string{"one", "two", "four", "five"} {
		require.NoError(t, store.AddUserBan(context.Background(), db.AddUserBanParams{
			UserID: user,
			Reason: "testing",
		}))
	}
	require.NoError(t, store.LiftUserBan(context.Background(), db.LiftUserBanParams{
		UserID: "five",
		Reason: "just testing",
	}))

	bans, err := LoadUserBans(store)
	assert.NoError(t, err)

	assert.True(t, bans.IsBanned("one"))
	assert.True(t, bans.IsBanned("two"))
	assert.False(t, bans.IsBanned("three")) // Never banned
	assert.False(t, bans.IsBanned("five"))  // Ban lifted
}
