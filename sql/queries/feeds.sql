-- name: AddFeed :one
INSERT INTO feed (id, created_at, updated_at, title, url, user_id)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetAllFeeds :many
SELECT * FROM feed;

-- name: CreateFeedFollow :one
WITH inserted_feed_follow AS (
    INSERT INTO feed_follows (
        id,
        created_at,
        updated_at,
        user_id,
        feed_id
    ) VALUES ($1, $2, $3, $4, $5)
    RETURNING *
)
SELECT 
    inserted_feed_follow.*,
    users.name AS user_name,
    feed.title AS feed_name
FROM inserted_feed_follow
INNER JOIN users
    ON inserted_feed_follow.user_id = users.id
INNER JOIN feed
    ON inserted_feed_follow.feed_id = feed.id;

-- name: GetFeedIdByFeedUrl :one
SELECT id FROM feed WHERE url = $1;

-- name: GetFeedFollowsForUser :many
SELECT
    feed.title
FROM feed_follows
INNER JOIN feed ON feed_follows.feed_id = feed.id
INNER JOIN users ON feed_follows.user_id = users.id
WHERE users.name = $1;

-- name: UnfollowFeed :exec
DELETE FROM feed_follows
WHERE user_id = $1 AND feed_id = $2;

-- name: MarkFeedFetched :exec
UPDATE feed 
SET last_fetched_at = $1, updated_at = $2
WHERE id = $3; 

-- name: GetNextFeedToFetch :one
SELECT url
FROM feed
ORDER BY last_fetched_at NULLS FIRST
LIMIT 1;

