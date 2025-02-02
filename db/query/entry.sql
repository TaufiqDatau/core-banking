-- name: InsertNewEntry :one
INSERT INTO entries(
  account_id,
  amount
) VALUES (
  $1,
  $2
)
RETURNING *;

-- name: UpdateEntry :one
UPDATE entries
SET amount = $2
WHERE id = $1
RETURNING *;

-- name: DeleteEntryById :one
DELETE FROM entries
WHERE id = $1
RETURNING *;

-- name: GetEntriesByAccountId :many
SELECT * FROM entries
WHERE account_id = $1
AND ($2::timestamp IS NULL OR created_at>= $2::timestamp)
AND ($3::timestamp IS NULL OR created_at<= $3::timestamp)
ORDER BY id
LIMIT $4;

