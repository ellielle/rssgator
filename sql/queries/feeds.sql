-- name: CreateFeed :one
INSERT INTO feeds(id, created_at, updated_at, name, url, user_id)
VALUES(
	$1,
	$2,
	$3,
	$4,
	$5,
	$6)
RETURNING *;

-- name: GetFeeds :many
SELECT * FROM feeds;

-- name: GetUserCreatedFeed :one
SELECT name FROM users 
WHERE id = $1;

-- name: GetFeedByUrl :one
SELECT * FROM feeds
WHERE url = $1; 

-- name: GetFeedById :one
SELECT * FROM feeds
WHERE id = $1;
