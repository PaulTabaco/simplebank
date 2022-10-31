package db

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"parus.i234.me/paultabaco/simplebank/util"
)

func createRandomTransfer(t *testing.T, fromAccountID int64, toAccountID int64, accountCreatedAt time.Time) Transfer {
	arg := CreateTransferParams{
		FromAccountID: fromAccountID,
		ToAccountID:   toAccountID,
		Amount:        util.RandomMoney(),
	}

	transfer, err := testQueries.CreateTransfer(context.Background(), arg)

	require.NoError(t, err)
	require.NotEmpty(t, transfer)

	require.Equal(t, arg.FromAccountID, transfer.FromAccountID)
	require.Equal(t, arg.ToAccountID, transfer.ToAccountID)
	require.Equal(t, arg.Amount, transfer.Amount)
	require.WithinDuration(t, accountCreatedAt, transfer.CreatedAt, time.Second)

	return transfer
}

func TestCreateTransfer(t *testing.T) {
	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	createRandomTransfer(t, account1.ID, account2.ID, account2.CreatedAt)
}

func TestGetTransfer(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)
	transfer1 := createRandomTransfer(t, account1.ID, account2.ID, account2.CreatedAt)

	transfer2, err := testQueries.GetTransfer(context.Background(), transfer1.ID)

	require.NoError(t, err)
	require.NotEmpty(t, transfer2)

	require.Equal(t, transfer1.ID, transfer2.ID)
	require.Equal(t, transfer1.FromAccountID, transfer2.FromAccountID)
	require.Equal(t, transfer1.ToAccountID, transfer2.ToAccountID)
	require.Equal(t, transfer1.Amount, transfer2.Amount)
	require.WithinDuration(t, transfer1.CreatedAt, transfer1.CreatedAt, time.Second)
}

// func TestUpdateEntry(t *testing.T) {
// 	/// this method not done in entry.sql yet !
// }

// func TestDeleteEntry(t *testing.T) {
// 	/// this method not done in entry.sql yet !
// }

func TestListTransfer(t *testing.T) {

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	for i := 0; i < 4; i++ {
		createRandomTransfer(t, account1.ID, account2.ID, account2.CreatedAt)
	}

	arg := ListTransfersParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Limit:         3,
		Offset:        1,
	}

	transfers, err := testQueries.ListTransfers(context.Background(), arg)
	require.NoError(t, err)
	require.Len(t, transfers, 3)

	for _, transfer := range transfers {
		require.NotEmpty(t, transfer)
	}

}
