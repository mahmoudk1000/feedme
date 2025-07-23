-- name: FollowFeed :one
INSERT INTO feed_follows (id, created_at, updated_at, feed_id, user_id)
VALUES ($1, $2, $3, $4, $5)
RETURNING *;

-- name: GetFeedFollowsByUser :many
SELECT * FROM feed_follows
WHERE user_id = $1
ORDER BY created_at DESC;

-- name: CheckFeedFollow :one
SELECT EXISTS (
    SELECT 1
    FROM feed_follows
    WHERE feed_id = $1 AND user_id = $2
) AS found;

-- name: DeleteFeedFollow :exec
DELETE FROM feed_follows
WHERE feed_id = $1 AND user_id = $2;
