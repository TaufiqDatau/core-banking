package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"math/big"
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

		txName := ctx.Value(txKey)

		log.Println(txName, "create transfer")
		result.Transfer, trxErr = q.InsertNewTransfer(ctx,
			InsertNewTransferParams{
				FromAccountID: arg.SenderAccountId,
				ToAccountID:   arg.ReceiverAccountId,
				Amount:        arg.Amount,
			})

		if trxErr != nil {
			return trxErr
		}

		log.Println(txName, "create entry 1")
		result.SenderEntry, trxErr = q.InsertNewEntry(ctx, InsertNewEntryParams{
			AccountID: arg.SenderAccountId,
			Amount:    fmt.Sprintf("-%s", arg.Amount),
		})

		if trxErr != nil {
			return trxErr
		}

		log.Println(txName, "create entry 2")
		result.ReceiverEntry, trxErr = q.InsertNewEntry(ctx, InsertNewEntryParams{
			AccountID: arg.ReceiverAccountId,
			Amount:    arg.Amount,
		})

		if trxErr != nil {
			return trxErr
		}

		amount, _ := new(big.Float).SetString(arg.Amount)

		log.Println(txName, "get account 1")
		result.SenderAccount, trxErr = q.GetAccountByIdForUpdate(ctx, arg.SenderAccountId)

		if trxErr != nil {
			return trxErr
		}
		balance, _ := new(big.Float).SetString(result.SenderAccount.Balance)

		log.Println(txName, "update balance 1")
		result.SenderAccount, trxErr = q.UpdateBalanceByAccountId(ctx, UpdateBalanceByAccountIdParams{
			ID:      arg.SenderAccountId,
			Balance: new(big.Float).Sub(balance, amount).Text('f', 2),
		})

		if trxErr != nil {
			return trxErr
		}

		log.Println(txName, "get account 2")
		result.ReceiverAccount, trxErr = q.GetAccountByIdForUpdate(ctx, arg.ReceiverAccountId)

		if trxErr != nil {
			return trxErr
		}
		balance, _ = new(big.Float).SetString(result.ReceiverAccount.Balance)

		log.Println(txName, "update balance 2")
		result.ReceiverAccount, trxErr = q.UpdateBalanceByAccountId(ctx, UpdateBalanceByAccountIdParams{
			ID:      arg.ReceiverAccountId,
			Balance: new(big.Float).Add(balance, amount).Text('f', 2),
		})

		if trxErr != nil {
			return trxErr
		}

		return nil
	})

	if err != nil {
		return result, err
	}

	return result, nil
}
