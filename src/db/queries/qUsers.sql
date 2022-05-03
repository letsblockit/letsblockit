-- name: GetBannedUsers :many
SELECT user_id
from banned_users
WHERE lifted_at IS NULL;

-- name: InitUserPreferences :one
INSERT INTO user_preferences (user_id)
VALUES ($1)
RETURNING *;

-- name: GetUserPreferences :one
SELECT *
FROM user_preferences
WHERE user_id = $1;

-- name: BumpLatestNews :exec
UPDATE user_preferences
SET latest_news = NOW()
WHERE user_id = $1;
