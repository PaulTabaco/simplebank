package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"parus.i234.me/paultabaco/simplebank/util"
)

func createRandomEntry(t *testing.T, accountID int64, accountCreatedAt time.Time) Entry {
	arg := CreateEntryParams{
		AccountID: accountID,
		Amount:    util.RandomMoney(),
	}

	entry, err := testQueries.CreateEntry(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, entry)

	require.Equal(t, accountID, entry.AccountID)
	require.Equal(t, arg.Amount, entry.Amount)
	require.WithinDuration(t, accountCreatedAt, entry.CreatedAt, time.Second)

	return entry
}

func TestCreateEntry(t *testing.T) {
	account1 := createRandomAccount(t)
	createRandomEntry(t, account1.ID, account1.CreatedAt)
}

func TestGetEntry(t *testing.T) {

	account1 := createRandomAccount(t)
	entry1 := createRandomEntry(t, account1.ID, account1.CreatedAt)

	entry2, err := testQueries.GetEntry(context.Background(), entry1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, entry2)

	require.Equal(t, entry1.ID, entry2.ID)
	require.Equal(t, entry1.AccountID, entry2.AccountID)
	require.Equal(t, entry1.Amount, entry2.Amount)
	require.WithinDuration(t, entry1.CreatedAt, entry2.CreatedAt, time.Second)
}

// func TestUpdateEntry(t *testing.T) {
// 	/// this method not done in entry.sql yet !
// }

// func TestDeleteEntry(t *testing.T) {
// 	/// this method not done in entry.sql yet !
// }

func TestListEntry(t *testing.T) {

	account1 := createRandomAccount(t)

	for i := 0; i < 4; i++ {
		createRandomEntry(t, account1.ID, account1.CreatedAt)
	}

	arg := ListEntriesParams{
		AccountID: account1.ID,
		Limit:     3,
		Offset:    1,
	}

	entries, err := testQueries.ListEntries(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, entries, 3)

	// for _, entry := range entries {
	// 	require.NotEmpty(t, entry)
	// }

}
