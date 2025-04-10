// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: feeds.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createFeed = `-- name: CreateFeed :one
INSERT INTO feeds (id, created_at, updated_at, user_id, name, url)
VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6
)
RETURNING id, created_at, updated_at, user_id, name, url, last_fetched_at
`

type CreateFeedParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.NullUUID
	Name      sql.NullString
	Url       sql.NullString
}

func (q *Queries) CreateFeed(ctx context.Context, arg CreateFeedParams) (Feed, error) {
	row := q.db.QueryRowContext(ctx, createFeed,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.UserID,
		arg.Name,
		arg.Url,
	)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Name,
		&i.Url,
		&i.LastFetchedAt,
	)
	return i, err
}

const getFeeds = `-- name: GetFeeds :many
SELECT feeds.id, feeds.name, feeds.url, users.name as username
FROM feeds
JOIN users
ON users.id = feeds.user_id
WHERE feeds.name <> '_g_invalid'
`

type GetFeedsRow struct {
	ID       uuid.UUID
	Name     sql.NullString
	Url      sql.NullString
	Username sql.NullString
}

func (q *Queries) GetFeeds(ctx context.Context) ([]GetFeedsRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeeds)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsRow
	for rows.Next() {
		var i GetFeedsRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.Username,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getFeedsByUrl = `-- name: GetFeedsByUrl :one
SELECT id, name, url, user_id
FROM feeds
WHERE url = $1
`

type GetFeedsByUrlRow struct {
	ID     uuid.UUID
	Name   sql.NullString
	Url    sql.NullString
	UserID uuid.NullUUID
}

func (q *Queries) GetFeedsByUrl(ctx context.Context, url sql.NullString) (GetFeedsByUrlRow, error) {
	row := q.db.QueryRowContext(ctx, getFeedsByUrl, url)
	var i GetFeedsByUrlRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.Url,
		&i.UserID,
	)
	return i, err
}

const getFeedsByUser = `-- name: GetFeedsByUser :many
SELECT id, name, url, user_id
FROM feeds
WHERE user_id = $1
`

type GetFeedsByUserRow struct {
	ID     uuid.UUID
	Name   sql.NullString
	Url    sql.NullString
	UserID uuid.NullUUID
}

func (q *Queries) GetFeedsByUser(ctx context.Context, userID uuid.NullUUID) ([]GetFeedsByUserRow, error) {
	rows, err := q.db.QueryContext(ctx, getFeedsByUser, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetFeedsByUserRow
	for rows.Next() {
		var i GetFeedsByUserRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.Url,
			&i.UserID,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const getNextFeedToFetch = `-- name: GetNextFeedToFetch :one
select id, created_at, updated_at, user_id, name, url, last_fetched_at
from feeds f 
where f."name" <> '_g_invalid'
order by f.last_fetched_at asc nulls first
limit 1
`

func (q *Queries) GetNextFeedToFetch(ctx context.Context) (Feed, error) {
	row := q.db.QueryRowContext(ctx, getNextFeedToFetch)
	var i Feed
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.UserID,
		&i.Name,
		&i.Url,
		&i.LastFetchedAt,
	)
	return i, err
}

const markFeedFetched = `-- name: MarkFeedFetched :exec
UPDATE feeds
SET updated_at = $2,
last_fetched_at = $2
WHERE id = $1
`

type MarkFeedFetchedParams struct {
	ID        uuid.UUID
	UpdatedAt time.Time
}

func (q *Queries) MarkFeedFetched(ctx context.Context, arg MarkFeedFetchedParams) error {
	_, err := q.db.ExecContext(ctx, markFeedFetched, arg.ID, arg.UpdatedAt)
	return err
}

const resetFeeds = `-- name: ResetFeeds :exec
DELETE FROM feeds
`

func (q *Queries) ResetFeeds(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetFeeds)
	return err
}
