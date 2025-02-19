// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: account.sql

package db

import (
	"context"
)

const addAccountBalanceById = `-- name: AddAccountBalanceById :one
UPDATE accounts
SET balance = balance + $1
WHERE id = $2
RETURNING id, owner, balance, created_at
`

type AddAccountBalanceByIdParams struct {
	Amount string `json:"amount"`
	ID     int64  `json:"id"`
}

func (q *Queries) AddAccountBalanceById(ctx context.Context, arg AddAccountBalanceByIdParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, addAccountBalanceById, arg.Amount, arg.ID)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

const createAccount = `-- name: CreateAccount :one
INSERT INTO accounts(
  owner,
  balance
) VALUES (
  $1, 
  $2
)
RETURNING id, owner, balance, created_at
`

type CreateAccountParams struct {
	Owner   string `json:"owner"`
	Balance string `json:"balance"`
}

func (q *Queries) CreateAccount(ctx context.Context, arg CreateAccountParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, createAccount, arg.Owner, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

const deleteAccount = `-- name: DeleteAccount :one
DELETE FROM accounts
WHERE id = $1
RETURNING id, owner, balance, created_at
`

func (q *Queries) DeleteAccount(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, deleteAccount, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountById = `-- name: GetAccountById :one
SELECT id, owner, balance, created_at FROM accounts
WHERE id = $1
LIMIT 1
`

func (q *Queries) GetAccountById(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountById, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountByIdForUpdate = `-- name: GetAccountByIdForUpdate :one
SELECT id, owner, balance, created_at FROM accounts
WHERE id = $1
LIMIT 1
FOR NO KEY UPDATE
`

func (q *Queries) GetAccountByIdForUpdate(ctx context.Context, id int64) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountByIdForUpdate, id)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

const getAccountFromName = `-- name: GetAccountFromName :one
SELECT id, owner, balance, created_at FROM accounts 
WHERE owner= $1 
LIMIT 1
`

func (q *Queries) GetAccountFromName(ctx context.Context, owner string) (Account, error) {
	row := q.db.QueryRowContext(ctx, getAccountFromName, owner)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}

const getListAccount = `-- name: GetListAccount :many
SELECT id, owner, balance, created_at FROM accounts 
LIMIT $1
`

func (q *Queries) GetListAccount(ctx context.Context, limit int32) ([]Account, error) {
	rows, err := q.db.QueryContext(ctx, getListAccount, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Account
	for rows.Next() {
		var i Account
		if err := rows.Scan(
			&i.ID,
			&i.Owner,
			&i.Balance,
			&i.CreatedAt,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateBalanceByAccountId = `-- name: UpdateBalanceByAccountId :one
UPDATE accounts
SET balance = $2
WHERE id = $1
RETURNING id, owner, balance, created_at
`

type UpdateBalanceByAccountIdParams struct {
	ID      int64  `json:"id"`
	Balance string `json:"balance"`
}

func (q *Queries) UpdateBalanceByAccountId(ctx context.Context, arg UpdateBalanceByAccountIdParams) (Account, error) {
	row := q.db.QueryRowContext(ctx, updateBalanceByAccountId, arg.ID, arg.Balance)
	var i Account
	err := row.Scan(
		&i.ID,
		&i.Owner,
		&i.Balance,
		&i.CreatedAt,
	)
	return i, err
}
