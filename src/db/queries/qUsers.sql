-- name: GetBannedUsers :many
SELECT user_id
from banned_users
WHERE lifted_at IS NULL;