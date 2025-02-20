package db

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	fmt.Println("initial --> ", account1.Balance, account2.Balance)
	senderBalanceBeforeTransaction, _ := new(big.Float).SetString(account1.Balance)
	receiverBalanceBeforeTransaction, _ := new(big.Float).SetString(account2.Balance)

	n, amount := 5, 10

	errs := make(chan error)
	results := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				SenderAccountId:   account1.ID,
				ReceiverAccountId: account2.ID,
				Amount:            fmt.Sprintf("%d", amount)})
			errs <- err
			results <- result
		}()
	}

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

		zero := new(big.Float).SetInt64(0)
		require.Greater(t, senderBalanceAfterTransaction.Cmp(zero), 0) // Ensures diff1 > 0

		receiverBalanceAfterTransaction, _ := new(big.Float).SetString(result.ReceiverAccount.Balance)

		diff2 := new(big.Float).Sub(receiverBalanceAfterTransaction, receiverBalanceBeforeTransaction)
		require.True(t, diff2.Cmp(zero) > 0)
	}
	senderAccount, err := store.GetAccountById(context.Background(), account1.ID)
	require.NoError(t, err)
	senderAccountBalance, _ := new(big.Float).SetString(senderAccount.Balance)

	fmt.Println("after tx ->", senderAccountBalance)
	require.True(t, fmt.Sprintf("%.2f", senderAccountBalance) == fmt.Sprintf("%.2f", new(big.Float).Sub(senderBalanceBeforeTransaction, new(big.Float).SetInt64(50))))

	receiverAccount, err := store.GetAccountById(context.Background(), account2.ID)
	require.NoError(t, err)
	receiverAccountBalance, _ := new(big.Float).SetString(receiverAccount.Balance)
	expectedReceiverAccountBalance := new(big.Float).Add(new(big.Float).SetInt64(50), receiverBalanceBeforeTransaction)
	require.True(t, fmt.Sprintf("%.2f", receiverAccountBalance) == fmt.Sprintf("%.2f", expectedReceiverAccountBalance))
}

func TestTransferTxForDeadlock(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	account3 := createRandomAccount(t)

	log.Printf("initial value of each account: %s->%s, %s->%s, %s->%s", account1.Owner, account1.Balance, account2.Owner, account2.Balance, account3.Owner, account3.Balance)

	n := 5
	amount := 10

	errs := make(chan error)
	errs2 := make(chan error)
	resultTf1 := make(chan TransferTxResult)
	resultTf2 := make(chan TransferTxResult)

	for i := 0; i < n; i++ {
		go func() {
			ctx := context.Background()
			result, err := store.TransferTx(ctx, TransferTxParams{
				SenderAccountId:   account1.ID,
				ReceiverAccountId: account2.ID,
				Amount:            fmt.Sprint(amount),
			})
			result2, err2 := store.TransferTx(ctx, TransferTxParams{
				SenderAccountId:   account2.ID,
				ReceiverAccountId: account3.ID,
				Amount:            fmt.Sprint(amount),
			})

			errs <- err
			errs2 <- err2
			resultTf2 <- result2
			resultTf1 <- result
		}()
	}

	for i := 0; i < n; i++ {
		err1 := <-errs
		err2 := <-errs2

		require.NoError(t, err1)
		require.NoError(t, err2)
	}

	//Checking is account 1 reduce the right amount
	currentBalanceAcc1, _ := new(big.Float).SetString(account1.Balance)
	expectedBalanceAccount1 := new(big.Float).Sub(currentBalanceAcc1, new(big.Float).SetInt64(int64(n*amount)))
	remainingBalance1, _ := store.GetAccountById(context.Background(), account1.ID)
	require.True(t, remainingBalance1.Balance == fmt.Sprintf("%.2f", expectedBalanceAccount1))

	//Checking account 2 is not reduce or add up
	Account2AfterTransaction, _ := store.GetAccountById(context.Background(), account2.ID)
	require.True(t, account2.Balance == Account2AfterTransaction.Balance)

	//Checking account3 need to be +50 then started
	currentBalanceAcc3, _ := new(big.Float).SetString(account3.Balance)
	expectedBalanceAccount3 := new(big.Float).Add(currentBalanceAcc3, new(big.Float).SetInt64(int64(n*amount)))
	currentAccount3, _ := store.GetAccountById(context.Background(), account3.ID)
	require.True(t, currentAccount3.Balance == fmt.Sprintf("%.2f", expectedBalanceAccount3))
}
