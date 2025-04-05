-- +goose Up
CREATE TABLE feeds (
  id UUID PRIMARY KEY,
  created_at TIMESTAMP NOT NULL,
  updated_at TIMESTAMP NOT NULL,
  user_id UUID REFERENCES users (id) ON DELETE CASCADE,
  name TEXT UNIQUE DEFAULT '_g_invalid',
  url TEXT UNIQUE
);

-- +goose Down
DROP TABLE feeds;
