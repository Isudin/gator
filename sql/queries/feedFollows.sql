-- Name: CreateFeedFollow :one
WITH createdFeedFollow  (
    INSERT INTO feed_follows (id, created_at, updated_at, user_id, feed_id)
    VALUES ($1, $2, $3, $4, $5);
    RETURNING *;
    )
SELECT createdFeedFollow, users.name, feeds.name 
FROM feeds
JOIN users ON users.id = feeds.user_id
WHERE feeds.id = createdFeedFollow.feed_id;