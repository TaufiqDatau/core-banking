package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomEntry(t *testing.T) Entry {
	account := createRandomAccount(t)
	arg := InsertNewEntryParams{
		AccountID: account.ID,
		Amount:    account.Balance,
	}

	newEntry, err := testQueries.InsertNewEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, newEntry)

	return newEntry
}

func TestCreateEntry(t *testing.T) {
	account := createRandomAccount(t)
	arg := InsertNewEntryParams{
		AccountID: account.ID,
		Amount:    account.Balance,
	}

	newEntry, err := testQueries.InsertNewEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, newEntry.ID)
	require.Equal(t, account.Balance, newEntry.Amount)
	require.Equal(t, account.ID, newEntry.AccountID)

	testQueries.DeleteEntryById(context.Background(), newEntry.ID)
	testQueries.DeleteAccount(context.Background(), account.ID)
}
