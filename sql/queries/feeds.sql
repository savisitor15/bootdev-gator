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
SELECT feeds.id, feeds.name, feeds.url, users.name as username
FROM feeds
JOIN users
ON users.id = feeds.user_id
WHERE feeds.name <> '_g_invalid';

-- name: GetFeedsByUser :many
SELECT id, name, url, user_id
FROM feeds
WHERE user_id = $1;

-- name: GetFeedsByUrl :one
SELECT id, name, url, user_id
FROM feeds
WHERE url = $1;

-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $2,
last_fetched_at = $2
WHERE id = $1;

-- name: GetNextFeedToFetch :one
select *
from feeds f 
where f."name" <> '_g_invalid'
order by f.last_fetched_at asc nulls first
limit 1;

-- name: ResetFeeds :exec
DELETE FROM feeds;
