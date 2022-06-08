package users

import "context"

type banQuerier interface {
	GetBannedUsers(ctx context.Context) ([]string, error)
}

type BanManager struct {
	bans map[string]struct{}
}

func LoadUserBans(store banQuerier) (*BanManager, error) {
	users, err := store.GetBannedUsers(context.Background())
	if err != nil {
		return nil, err
	}
	bans := make(map[string]struct{}, len(users))
	for _, u := range users {
		bans[u] = struct{}{}
	}
	return &BanManager{bans: bans}, nil
}

func (m *BanManager) IsBanned(id string) bool {
	if m == nil {
		return false // For unit tests
	}
	_, found := m.bans[id]
	return found
}
