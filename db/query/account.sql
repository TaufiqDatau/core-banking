-- name: GetAccountFromName :one
SELECT * FROM accounts 
WHERE owner= $1 
LIMIT 1;

-- name: CreateAccount :one
INSERT INTO accounts(
  owner,
  balance
) VALUES (
  $1, 
  $2
)
RETURNING *;

-- name: UpdateBalanceByAccountId :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING *;

-- name: DeleteAccount :one
DELETE FROM accounts
WHERE id = $1
RETURNING *;
