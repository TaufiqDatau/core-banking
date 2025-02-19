package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
)

// Store provide all functions to execute db queries and transaction
type Store struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		db:      db,
		Queries: New(db),
	}
}

var txKey = struct{}{}

func (store *Store) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}

	q := New(tx)

	err = fn(q)

	if err != nil {
		if rbError := tx.Rollback(); rbError != nil {
			return fmt.Errorf("tx error: %v, rb error: %v", err, rbError)
		}
		return err
	}

	return tx.Commit()
}

type TransferTxParams struct {
	SenderAccountId   int64  `json:"sender_account_id"`
	ReceiverAccountId int64  `json:"receiver_account_id"`
	Amount            string `json:"amount"`
}
type TransferTxResult struct {
	Transfer        Transfer `json:"transfer"`
	SenderAccount   Account  `json:"sender_account"`
	ReceiverAccount Account  `json:"receiver_account"`
	SenderEntry     Entry    `json:"sender_entry"`
	ReceiverEntry   Entry    `json:"receiver_entry"`
}

func (store *Store) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {

	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var trxErr error

		result.Transfer, trxErr = q.InsertNewTransfer(ctx,
			InsertNewTransferParams{
				FromAccountID: arg.SenderAccountId,
				ToAccountID:   arg.ReceiverAccountId,
				Amount:        arg.Amount,
			})

		if trxErr != nil {
			return trxErr
		}

		result.SenderEntry, trxErr = q.InsertNewEntry(ctx, InsertNewEntryParams{
			AccountID: arg.SenderAccountId,
			Amount:    fmt.Sprintf("-%s", arg.Amount),
		})

		if trxErr != nil {
			return trxErr
		}

		result.ReceiverEntry, trxErr = q.InsertNewEntry(ctx, InsertNewEntryParams{
			AccountID: arg.ReceiverAccountId,
			Amount:    arg.Amount,
		})

		if trxErr != nil {
			return trxErr
		}

		result.SenderAccount, trxErr = q.AddAccountBalanceById(ctx, AddAccountBalanceByIdParams{
			Amount: fmt.Sprintf("-%s", arg.Amount),
			ID:     arg.SenderAccountId,
		})

		if trxErr != nil {
			return trxErr
		}

		result.ReceiverAccount, trxErr = q.AddAccountBalanceById(ctx, AddAccountBalanceByIdParams{
			ID:     arg.ReceiverAccountId,
			Amount: arg.Amount,
		})

		if trxErr != nil {
			return trxErr
		}

		log.Printf("Transfer with ID %d succesfully created", result.Transfer.ID)
		return nil
	})

	if err != nil {
		log.Printf("transaction failed with this error: %s", err.Error())
		return result, err
	}

	return result, nil
}

func UpdateAccountBalanceAfterTransaction(ctx context.Context, q *Queries, account1, account2 Account, amount string) (Account, Account, error) {
	var err error
	account1, err = q.AddAccountBalanceById(ctx, AddAccountBalanceByIdParams{
		ID:     account1.ID,
		Amount: fmt.Sprintf("-%s", amount),
	})

	if err != nil {
		return account1, account2, err
	}

	account2, err = q.AddAccountBalanceById(ctx, AddAccountBalanceByIdParams{
		ID:     account2.ID,
		Amount: amount,
	})

	if err != nil {
		return account1, account2, err
	}

	return account1, account2, nil
}
