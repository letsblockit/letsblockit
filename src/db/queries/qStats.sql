-- name: GetStats :one
SELECT (SELECT COUNT(*) FROM filter_lists)                                                  as lists_total,
       (SELECT COUNT(*) FROM filter_lists WHERE downloaded IS TRUE)                         as lists_active,
       (SELECT COUNT(*) FROM filter_lists WHERE downloaded_at >= NOW() - INTERVAL '7 DAYS') as lists_fresh;

-- name: GetInstanceStats :many
SELECT COUNT(*) as total,
       SUM(case when l.downloaded_at >= NOW() - INTERVAL '7 DAYS' then 1 else 0 end) as fresh,
       template_name
FROM filter_instances
         JOIN filter_lists AS l ON (list_id = l.id)
GROUP BY template_name;
