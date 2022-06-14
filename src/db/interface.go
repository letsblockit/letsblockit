// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0

package db

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgtype"
)

type Querier interface {
	CountInstanceForUserAndFilter(ctx context.Context, arg CountInstanceForUserAndFilterParams) (int64, error)
	CountListsForUser(ctx context.Context, userID string) (int64, error)
	CreateInstanceForUserAndFilter(ctx context.Context, arg CreateInstanceForUserAndFilterParams) error
	CreateListForUser(ctx context.Context, userID string) (uuid.UUID, error)
	DeleteInstanceForUserAndFilter(ctx context.Context, arg DeleteInstanceForUserAndFilterParams) error
	GetActiveFiltersForUser(ctx context.Context, userID string) ([]GetActiveFiltersForUserRow, error)
	GetBannedUsers(ctx context.Context) ([]string, error)
	GetInstanceForUserAndFilter(ctx context.Context, arg GetInstanceForUserAndFilterParams) (pgtype.JSONB, error)
	GetInstanceStats(ctx context.Context) ([]GetInstanceStatsRow, error)
	GetInstancesForList(ctx context.Context, filterListID int32) ([]GetInstancesForListRow, error)
	GetListForToken(ctx context.Context, token uuid.UUID) (GetListForTokenRow, error)
	GetListForUser(ctx context.Context, userID string) (GetListForUserRow, error)
	GetStats(ctx context.Context) (GetStatsRow, error)
	GetUserPreferences(ctx context.Context, userID string) (UserPreference, error)
	HasUserDownloadedList(ctx context.Context, userID string) (bool, error)
	InitUserPreferences(ctx context.Context, userID string) (UserPreference, error)
	MarkListDownloaded(ctx context.Context, id int32) error
	RotateListToken(ctx context.Context, arg RotateListTokenParams) error
	UpdateInstanceForUserAndFilter(ctx context.Context, arg UpdateInstanceForUserAndFilterParams) error
	UpdateNewsCursor(ctx context.Context, arg UpdateNewsCursorParams) error
}

var _ Querier = (*Queries)(nil)
