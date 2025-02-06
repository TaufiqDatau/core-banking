package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/taufiqDatau/core-banking/util"
)

func createRandomAccount(t *testing.T) Account {
	arg := CreateAccountParams{
		Owner:   util.RandomString(8),
		Balance: util.RandomBalance(50, 100),
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
	for i := 0; i < 5; i++ {
		createRandomAccount(t)
	}

	// Fetch accounts
	listAccountData, err := testQueries.GetListAccount(context.Background(), 5)

	// Assertions
	require.NoError(t, err)
	require.Equal(t, 5, len(listAccountData))

	// Verify that fetched accounts are valid
	for _, acc := range listAccountData {
		require.NotEmpty(t, acc)
		testQueries.DeleteAccount(context.Background(), acc.ID)
	}
}

func TestGetAccountByName(t *testing.T) {
	createdAccount := createRandomAccount(t)

	account, err := testQueries.GetAccountFromName(context.Background(), createdAccount.Owner)

	require.NoError(t, err)
	require.Equal(t, createdAccount, account)

	testQueries.DeleteAccount(context.Background(), createdAccount.ID)
}
