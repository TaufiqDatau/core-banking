-- name: InsertNewTransfer :one
INSERT INTO transfers (
  from_account_id,
  to_account_id,
  amount
) VALUES ( 
  $1,
  $2,
  $3
)
RETURNING *;

-- name: UpdateTransferAmount :one
UPDATE transfers
  SET amount = $2 
  WHERE id = $1
  RETURNING *;

-- name: DeleteTransferById :one
DELETE FROM transfers
  WHERE id = $1
  RETURNING *;

-- name: GetTransferFromSenderId :many
SELECT * FROM transfers
WHERE from_account_id = $1
AND ($2::timestamp IS NULL OR created_at>= $2::timestamp)
AND ($3::timestamp IS NULL OR created_at<= $3::timestamp)
ORDER BY created_at
LIMIT $4
;
