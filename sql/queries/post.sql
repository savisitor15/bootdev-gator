-- name: CreatePost :one
INSERT INTO posts (id, created_at, updated_at, title, url, description, published_at, feed_id)
VALUES (
$1,
$2,
$2,
$3,
$4,
$5,
$6,
$7
)
RETURNING *;

-- name: GetPostsForUser :many
select p.*, f."name" as feed_name, f.url as feed_url, f.last_fetched_at as feed_last_update
from posts p 
left join feed_follows ff on ff.feed_id  = p.feed_id
left join feeds f on f.id = p.feed_id
where ff.user_id = $1
group by f."name", f.url, f.last_fetched_at, p.id
order by p.updated_at desc
limit $2;