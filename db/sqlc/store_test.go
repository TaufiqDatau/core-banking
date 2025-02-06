package db

import (
	"context"
	"fmt"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	fmt.Println("initial --> ", account1.Balance, account2.Balance)

	n, amount := 5, 10

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				SenderAccountId:   account1.ID,
				ReceiverAccountId: account2.ID,
				Amount:            fmt.Sprintf("%d", amount),
			})

			errs <- err
			results <- result
		}()
	}

	amountString := new(big.Float).SetInt64(int64(amount))
	// Check Result
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results

		fmt.Println("trx --> ", result.SenderAccount.Balance, result.ReceiverAccount.Balance)
		require.NotEmpty(t, result)
		transfer := result.Transfer

		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)

		_, err = store.GetTransferFromId(context.Background(), transfer.ID)
		require.NoError(t, err)

		senderBalanceAfterTransaction, _ := new(big.Float).SetString(result.SenderAccount.Balance)
		senderBalanceBeforeTransaction, _ := new(big.Float).SetString(account1.Balance)
		diff1 := new(big.Float).Sub(senderBalanceAfterTransaction, amountString)

		zero := new(big.Float).SetInt64(0)
		require.Greater(t, diff1.Cmp(zero), 0) // Ensures diff1 > 0
		require.Greater(t, senderBalanceBeforeTransaction.Cmp(senderBalanceAfterTransaction), 0)
	}

}
