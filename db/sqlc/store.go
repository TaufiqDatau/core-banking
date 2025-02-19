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

		_, trxErr = q.LockAccountForTransfer(ctx, LockAccountForTransferParams{
			Fromaccounid: arg.SenderAccountId,
			Toaccountid:  arg.ReceiverAccountId,
		})

		if trxErr != nil {
			return trxErr
		}

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

		result.SenderAccount, result.ReceiverAccount, trxErr = UpdateAccountBalanceAfterTransaction(ctx, q, arg.SenderAccountId, arg.ReceiverAccountId, arg.Amount)
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

func UpdateAccountBalanceAfterTransaction(ctx context.Context, q *Queries, account1Id, account2Id int64, amount string) (Account, Account, error) {
	var err error
	var result1, result2 Account
	result1, err = q.AddAccountBalanceById(ctx, AddAccountBalanceByIdParams{
		ID:     account1Id,
		Amount: fmt.Sprintf("-%s", amount),
	})

	if err != nil {
		return result1, result2, err
	}

	result2, err = q.AddAccountBalanceById(ctx, AddAccountBalanceByIdParams{
		ID:     account2Id,
		Amount: amount,
	})

	if err != nil {
		return result1, result2, err
	}

	return result1, result2, nil
}
