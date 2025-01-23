// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: entry.sql

package db

import (
	"context"
)

const updateEntry = `-- name: UpdateEntry :one
INSERT INTO entries(
  account_id,
  amount
) VALUES (
  $1,
  $2
)
RETURNING id, account_id, amount, created_at
`

type UpdateEntryParams struct {
	AccountID int64  `json:"account_id"`
	Amount    string `json:"amount"`
}

func (q *Queries) UpdateEntry(ctx context.Context, arg UpdateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, updateEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
