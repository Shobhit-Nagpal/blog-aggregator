-- name: CreateFeed :one

INSERT INTO feeds (id, created_at, updated_at, name, url, user_id)
VALUES (
  $1, 
  $2, 
  $3, 
  $4, 
  $5, 
  $6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: CreateFeedFollow :many
WITH inserted_feed_follow AS (
INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
VALUES (
  $1,
  $2,
  $3,
  $4,
  $5
  )
RETURNING *
) 
SELECT inserted_feed_follow.*,
  feeds.name as feed_name,
  users.name as user_name
  FROM inserted_feed_follow
  INNER JOIN feeds ON inserted_feed_follow.feed_id = feeds.id
  INNER JOIN users ON inserted_feed_follow.user_id = users.id;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE $1 = url;

-- name: GetFeedById :one
SELECT * FROM feeds
WHERE $1 = id;

-- name: GetFeedFollowsForUser :many
SELECT * FROM feed_follows
WHERE feed_follows.user_id = $1;
