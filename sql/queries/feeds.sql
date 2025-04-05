-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, user_id, name, url)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING *;

-- name: GetFeeds :many
SELECT id, name, url, user_id
FROM feeds
WHERE name <> '_g_invalid';

-- name: GetFeedByUser :one
SELECT id, name, url, user_id
FROM feeds
WHERE user_id = $1;

-- name: ResetFeeds :exec
DELETE FROM feeds;
