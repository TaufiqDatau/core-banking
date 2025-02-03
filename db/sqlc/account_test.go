package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/taufiqDatau/core-banking/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:   util.RandomString(8),
		Balance: util.RandomBalance(1, 100),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)

	return account
}

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:   util.RandomString(8),
		Balance: util.RandomBalance(1, 100),
	}

	account, err := testQueries.CreateAccount(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, account)

	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)

	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
	testQueries.DeleteAccount(context.Background(), account.ID)
}

func TestDeleteAccount(t *testing.T) {
	account := createRandomAccount(t)

	deletedAccount, err := testQueries.DeleteAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.Equal(t, deletedAccount, account)
}

func TestUpdateAccount(t *testing.T) {
	account := createRandomAccount(t)

	arg := UpdateBalanceByAccountIdParams{
		ID:      account.ID,
		Balance: util.RandomBalance(1, 100),
	}

	updatedAccount, err := testQueries.UpdateBalanceByAccountId(context.Background(), arg)

	require.NoError(t, err)
	require.Equal(t, updatedAccount.Balance, arg.Balance)

	testQueries.DeleteAccount(context.Background(), updatedAccount.ID)
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

	getNewlyCreatedAccount, err := testQueries.GetAccountById(context.Background(), account.ID)

	require.NoError(t, err)
	require.Equal(t, account, getNewlyCreatedAccount)

	testQueries.DeleteAccount(context.Background(), getNewlyCreatedAccount.ID)
}

func TestListAccount(t *testing.T) {
	// Insert test accounts
	var createdAccounts []Account
	for i := 0; i < 5; i++ {
		account := createRandomAccount(t)
		createdAccounts = append(createdAccounts, account)
	}

	// Fetch accounts
	listAccountData, err := testQueries.GetListAccount(context.Background(), 5)

	// Assertions
	require.NoError(t, err)
	require.Equal(t, 5, len(listAccountData))

	// Verify that fetched accounts are valid
	for i, acc := range listAccountData {
		require.NotEmpty(t, acc)
		require.Equal(t, createdAccounts[i].Owner, acc.Owner)
		require.Equal(t, createdAccounts[i].Balance, acc.Balance)
		require.WithinDuration(t, createdAccounts[i].CreatedAt, acc.CreatedAt, time.Second)
		testQueries.DeleteAccount(context.Background(), acc.ID)
	}
}

func TestListAccount_LimitExceedsRecords(t *testing.T) {
	// Ensure there are only 3 accounts in the DB
	for i := 0; i < 3; i++ {
		createRandomAccount(t)
	}

	// Fetch more than available
	listAccountData, err := testQueries.GetListAccount(context.Background(), 10)

	require.NoError(t, err)
	require.LessOrEqual(t, len(listAccountData), 3) // Should return only available records

	for _, acc := range listAccountData {
		testQueries.DeleteAccount(context.Background(), acc.ID)
	}
}

func TestGetAccountByName(t *testing.T) {
	createdAccount := createRandomAccount(t)

	account, err := testQueries.GetAccountFromName(context.Background(), createdAccount.Owner)

	require.NoError(t, err)
	require.Equal(t, createdAccount, account)

}
