// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: users.sql

package database

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, name)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING id, created_at, updated_at, name
`

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      sql.NullString
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
	)
	return i, err
}

const getUser = `-- name: GetUser :one
SELECT id, created_at, updated_at, name
FROM users
WHERE name = $1
`

func (q *Queries) GetUser(ctx context.Context, name sql.NullString) (User, error) {
	row := q.db.QueryRowContext(ctx, getUser, name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
	)
	return i, err
}

const getUserByName = `-- name: GetUserByName :one
SELECT id, name 
FROM users
WHERE name = $1
`

type GetUserByNameRow struct {
	ID   uuid.UUID
	Name sql.NullString
}

func (q *Queries) GetUserByName(ctx context.Context, name sql.NullString) (GetUserByNameRow, error) {
	row := q.db.QueryRowContext(ctx, getUserByName, name)
	var i GetUserByNameRow
	err := row.Scan(&i.ID, &i.Name)
	return i, err
}

const getUsers = `-- name: GetUsers :many
SELECT id, name
FROM users
WHERE name <> '_g_invalid'
`

type GetUsersRow struct {
	ID   uuid.UUID
	Name sql.NullString
}

func (q *Queries) GetUsers(ctx context.Context) ([]GetUsersRow, error) {
	rows, err := q.db.QueryContext(ctx, getUsers)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []GetUsersRow
	for rows.Next() {
		var i GetUsersRow
		if err := rows.Scan(&i.ID, &i.Name); err != nil {
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

const resetUsers = `-- name: ResetUsers :exec
DELETE FROM users
WHERE name <> '_g_invalid'
`

func (q *Queries) ResetUsers(ctx context.Context) error {
	_, err := q.db.ExecContext(ctx, resetUsers)
	return err
}
