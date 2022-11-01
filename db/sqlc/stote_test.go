package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	// to check concurrency should
	// run n concurrent transfer transactions
	n := 5
	ammount := int64(10)

	errs := make(chan error)               // chanes is designet to connect  concurrent go-routins !!!
	results := make(chan TransferTxResult) // chanes result used becouse we not whure exact order

	//start  n go-routins
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        ammount,
			})
			errs <- err       // send an arror to errors chane
			results <- result // send resalr to resals chane
		}()
	}

	// check results
	for i := 0; i < n; i++ {

		err := <-errs
		require.NoError(t, err)

		result := <-results
		require.NotNil(t, result)

		/// Check transfers
		tranfer := result.Transfer
		require.NotEmpty(t, tranfer)
		require.Equal(t, account1.ID, tranfer.FromAccountID)
		require.Equal(t, account2.ID, tranfer.ToAccountID)
		require.Equal(t, ammount, tranfer.Amount)
		require.NotZero(t, tranfer.ID)
		require.NotZero(t, tranfer.CreatedAt)

		_, err = store.GetTransfer(context.Background(), tranfer.ID) /// check if record realy exists in db
		require.NoError(t, err)

		/// Check results
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -ammount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		/// Check entry
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, ammount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		/// TODO: - check account balance

	}

}
