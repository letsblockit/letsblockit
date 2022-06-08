package users

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

type mockedBanStore struct {
	bans []string
}

func (s *mockedBanStore) GetBannedUsers(_ context.Context) ([]string, error) {
	return s.bans, nil
}

func TestLoadBannedUsers(t *testing.T) {
	bans, err := LoadUserBans(&mockedBanStore{[]string{"one", "two", "one", "four"}})
	assert.NoError(t, err)

	assert.True(t, bans.IsBanned("one"))
	assert.True(t, bans.IsBanned("two"))
	assert.False(t, bans.IsBanned("three"))
}
