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
