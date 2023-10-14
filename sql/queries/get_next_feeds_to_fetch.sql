-- name: GetNextFeedsToFetch :many
SELECT *
FROM feeds
ORDER BY last_fetched_at ASC
LIMIT 3;

-- name: MarkFeedFetched :one
UPDATE feeds
SET last_fetched_at = CURRENT_TIMESTAMP,
    updated_at = CURRENT_TIMESTAMP
WHERE id = $1
RETURNING *;