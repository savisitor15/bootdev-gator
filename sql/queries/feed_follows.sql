-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (created_at, updated_at, user_id, feed_id)
    VALUES (
    $1,
    $2,
    $3,
    $4
    )
    RETURNING *
)
SELECT
    inserted_feed_follow.*,
    feeds.name as feed_name,
    users.name as user_name
FROM inserted_feed_follow
INNER JOIN feeds ON feeds.id = inserted_feed_follow.feed_id
INNER JOIN users ON users.id = inserted_feed_follow.user_id;

-- name: GetFeedFollowsForUser :many
SELECT ff.*, u.name as user_name, f.name as feed_name
FROM feed_follows AS ff
JOIN feeds AS f ON f.id = ff.feed_id
JOIN users AS u on u.id = ff.user_id
WHERE u.id = $1;

-- name: DeleteFeedFollowsForUser :exec
delete from feed_follows
where id in (select ff.id as feed_follows_id from feed_follows ff 
join users u on u.id = ff.user_id
join feeds f  on f.id = ff.feed_id
where u.name = $1 and f.url = $2);


