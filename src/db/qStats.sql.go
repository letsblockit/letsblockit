// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: qStats.sql

package db

import (
	"context"
)

const getInstanceStats = `-- name: GetInstanceStats :many
SELECT COUNT(*) as total,
       SUM(case when l.downloaded_at >= NOW() - INTERVAL '7 DAYS' then 1 else 0 end) as fresh,
       filter_name
FROM filter_instances
         JOIN filter_lists AS l ON (filter_list_id = l.id)
GROUP BY filter_name
`

type GetInstanceStatsRow struct {
	Total      int64
	Fresh      int64
	FilterName string
}

func (q *Queries) GetInstanceStats(ctx context.Context) ([]GetInstanceStatsRow, error) {
	rows, err := q.db.Query(ctx, getInstanceStats)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetInstanceStatsRow
	for rows.Next() {
		var i GetInstanceStatsRow
		if err := rows.Scan(&i.Total, &i.Fresh, &i.FilterName); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getStats = `-- name: GetStats :one
SELECT (SELECT COUNT(*) FROM filter_lists)                                                  as lists_total,
       (SELECT COUNT(*) FROM filter_lists WHERE downloaded IS TRUE)                         as lists_active,
       (SELECT COUNT(*) FROM filter_lists WHERE downloaded_at >= NOW() - INTERVAL '7 DAYS') as lists_fresh
`

type GetStatsRow struct {
	ListsTotal  int64
	ListsActive int64
	ListsFresh  int64
}

func (q *Queries) GetStats(ctx context.Context) (GetStatsRow, error) {
	row := q.db.QueryRow(ctx, getStats)
	var i GetStatsRow
	err := row.Scan(&i.ListsTotal, &i.ListsActive, &i.ListsFresh)
	return i, err
}
