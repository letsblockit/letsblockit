-- name: GetInstancesForUser :many
SELECT template_name, params, test_mode
FROM filter_instances
WHERE user_id = $1;

-- name: CreateInstance :exec
INSERT INTO filter_instances (list_id, user_id, template_name, params, test_mode)
VALUES ((SELECT id FROM filter_lists WHERE user_id = $1), $1, $2, $3, $4);

-- name: UpdateInstance :exec
UPDATE filter_instances
SET params     = $3,
    test_mode  = $4,
    updated_at = NOW()
WHERE (user_id = $1 AND template_name = $2);

-- name: GetInstance :one
SELECT params, test_mode
FROM filter_instances
WHERE (user_id = $1 AND template_name = $2);

-- name: CountInstances :one
SELECT COUNT(*)
FROM filter_instances
WHERE (user_id = $1 AND template_name = $2);

-- name: DeleteInstance :exec
DELETE
FROM filter_instances
WHERE (user_id = $1 AND template_name = $2);

-- name: GetInstancesForList :many
SELECT template_name, params, test_mode
FROM filter_instances
WHERE list_id = $1
ORDER BY template_name ASC;
