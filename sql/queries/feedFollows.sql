-- name: CreateFeedFollow :one
WITH created_feed_follow AS (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5)
    RETURNING *
    )
SELECT created_feed_follow.*, users.name AS user_name, feeds.name AS feed_name
FROM created_feed_follow
JOIN users ON users.id = created_feed_follow.user_id
JOIN feeds ON feeds.id = created_feed_follow.feed_id;