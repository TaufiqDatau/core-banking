// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.28.0
// source: entry.sql

package db

import (
	"context"
	"time"
)

const deleteEntryById = `-- name: DeleteEntryById :one
DELETE FROM entries
WHERE id = $1
RETURNING id, account_id, amount, created_at
`

func (q *Queries) DeleteEntryById(ctx context.Context, id int64) (Entry, error) {
	row := q.db.QueryRowContext(ctx, deleteEntryById, id)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const getEntriesByAccountId = `-- name: GetEntriesByAccountId :many
SELECT id, account_id, amount, created_at FROM entries
WHERE account_id = $1
AND ($2::timestamp IS NULL OR created_time >= $2::timestamp)
AND ($3::timestamp IS NULL OR created_time <= $3::timestamp)
ORDER BY id
LIMIT $4
`

type GetEntriesByAccountIdParams struct {
	AccountID int64     `json:"account_id"`
	FromTime   time.Time `json:"from_time"`
	ToTime   time.Time `json:"to_time"`
	Limit     int32     `json:"limit"`
}

func (q *Queries) GetEntriesByAccountId(ctx context.Context, arg GetEntriesByAccountIdParams) ([]Entry, error) {
	rows, err := q.db.QueryContext(ctx, getEntriesByAccountId,
		arg.AccountID,
		arg.FromTime,
		arg.ToTime,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []Entry
	for rows.Next() {
		var i Entry
		if err := rows.Scan(
			&i.ID,
			&i.AccountID,
			&i.Amount,
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

const insertNewEntry = `-- name: InsertNewEntry :one
INSERT INTO entries(
  account_id,
  amount
) VALUES (
  $1,
  $2
)
RETURNING id, account_id, amount, created_at
`

type InsertNewEntryParams struct {
	AccountID int64  `json:"account_id"`
	Amount    string `json:"amount"`
}

func (q *Queries) InsertNewEntry(ctx context.Context, arg InsertNewEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, insertNewEntry, arg.AccountID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}

const updateEntry = `-- name: UpdateEntry :one
UPDATE entries
SET amount = $2
WHERE id = $1
RETURNING id, account_id, amount, created_at
`

type UpdateEntryParams struct {
	ID     int64  `json:"id"`
	Amount string `json:"amount"`
}

func (q *Queries) UpdateEntry(ctx context.Context, arg UpdateEntryParams) (Entry, error) {
	row := q.db.QueryRowContext(ctx, updateEntry, arg.ID, arg.Amount)
	var i Entry
	err := row.Scan(
		&i.ID,
		&i.AccountID,
		&i.Amount,
		&i.CreatedAt,
	)
	return i, err
}
