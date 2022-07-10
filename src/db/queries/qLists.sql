-- name: CreateListForUser :one
INSERT INTO filter_lists (user_id)
VALUES ($1)
RETURNING token;

-- name: GetListForUser :one
SELECT token,
       downloaded,
       (SELECT COUNT(*) FROM filter_instances WHERE filter_instances.user_id = $1) AS instance_count
FROM filter_lists
WHERE filter_lists.user_id = $1
LIMIT 1;

-- name: CountListsForUser :one
SELECT COUNT(*)
FROM filter_lists
WHERE user_id = $1;

-- name: RotateListToken :exec
UPDATE filter_lists
SET token      = gen_random_uuid(),
    downloaded = false
WHERE user_id = $1
  AND token = $2;

-- name: HasUserDownloadedList :one
SELECT downloaded
FROM filter_lists
WHERE filter_lists.user_id = $1
LIMIT 1;

-- name: GetListForToken :one
SELECT id, user_id, downloaded
FROM filter_lists
WHERE token = $1
LIMIT 1;

-- name: MarkListDownloaded :exec
UPDATE filter_lists
SET downloaded    = true,
    downloaded_at = NOW()
WHERE token = $1;
