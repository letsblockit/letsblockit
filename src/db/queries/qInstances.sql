-- name: GetActiveFiltersForUser :many
SELECT filter_name, params
FROM filter_instances
WHERE user_id = $1;

-- name: CreateInstanceForUserAndFilter :exec
INSERT INTO filter_instances (filter_list_id, user_id, filter_name, params, test_mode)
VALUES ((SELECT id FROM filter_lists WHERE user_id = $1), $1, $2, $3, $4);

-- name: UpdateInstanceForUserAndFilter :exec
UPDATE filter_instances
SET params     = $3,
    test_mode  = $4,
    updated_at = NOW()
WHERE (user_id = $1 AND filter_name = $2);

-- name: GetInstanceForUserAndFilter :one
SELECT params, test_mode
FROM filter_instances
WHERE (user_id = $1 AND filter_name = $2);

-- name: CountInstanceForUserAndFilter :one
SELECT COUNT(*)
FROM filter_instances
WHERE (user_id = $1 AND filter_name = $2);

-- name: DeleteInstanceForUserAndFilter :exec
DELETE
FROM filter_instances
WHERE (user_id = $1 AND filter_name = $2);

-- name: GetInstancesForList :many
SELECT filter_name, params, test_mode
FROM filter_instances
WHERE filter_list_id = $1
ORDER BY filter_name ASC;
