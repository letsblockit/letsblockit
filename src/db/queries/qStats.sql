-- name: GetStats :one
SELECT (SELECT COUNT(*) FROM filter_lists)                                                  as lists_total,
       (SELECT COUNT(*) FROM filter_lists WHERE downloaded IS TRUE)                         as lists_active,
       (SELECT COUNT(*) FROM filter_lists WHERE downloaded_at >= NOW() - INTERVAL '7 DAYS') as lists_fresh;

-- name: GetInstanceStats :many
SELECT COUNT(*), filter_name
FROM filter_instances
GROUP BY filter_name;
