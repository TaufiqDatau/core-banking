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
}

func TestGetAccount(t *testing.T) {
	account := createRandomAccount(t)

	getNewlyCreatedAccount, err := testQueries.GetAccountById(context.Background(), account.ID)

	require.NoError(t, err)
	require.Equal(t, account, getNewlyCreatedAccount)
}

func TestListAccount(t *testing.T) {
	for i := 0; i < 5; i++ {
		createRandomAccount(t)
	}

	listAccountData, err := testQueries.GetListAccount(context.Background(), 5)

	require.NoError(t, err)
	require.Equal(t, 5, len(listAccountData))

	for _, acc := range listAccountData {
		require.NotEmpty(t, acc)
	}
}
