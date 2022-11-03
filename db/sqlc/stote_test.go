package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)

	account1 := createRandomAccount(t)
	account2 := createRandomAccount(t)

	/// TRY LOOK FOR EXPECTED DEADLOOK
	fmt.Println(">> before transaction: ", account1.Balance, account2.Balance)

	/// to check concurrency should
	/// run n concurrent transfer transactions
	n := 5
	ammount := int64(10)

	errs := make(chan error)               /// chanes is designet to connect  concurrent go-routins !!!
	results := make(chan TransferTxResult) /// chanes result used becouse we not whure exact order

	/// start  n go-routins
	for i := 0; i < n; i++ {
		go func() {
			result, err := store.TransferTx(context.Background(), TransferTxParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        ammount,
			})
			errs <- err       /// send an arror to errors chane
			results <- result /// send resalr to resals chane
		}()
	}

	/// check results
	existed := make(map[int]bool)

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

		/// Check entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.Equal(t, account1.ID, fromEntry.AccountID)
		require.Equal(t, -ammount, fromEntry.Amount)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)

		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.Equal(t, account2.ID, toEntry.AccountID)
		require.Equal(t, ammount, toEntry.Amount)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)

		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)

		/// Check accounts
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, account1.ID, fromAccount.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, account2.ID, toAccount.ID)

		/// Check accounts' balance
		fmt.Println(">> each tx: ", fromAccount.Balance, toAccount.Balance)

		diff1 := account1.Balance - fromAccount.Balance
		diff2 := toAccount.Balance - account2.Balance
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)

		/// ???? next ??? so if we tranfer few non equal amounts?
		require.True(t, diff1%ammount == 0) /// amount, 2 * amount, 3 * amount
		k := int(diff1 / ammount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}

	/// Check the final updated balance
	updatedAccount1, err := testQueries.GetAccount(context.Background(), account1.ID) /// ? why not - store.GetAccount(
	require.NoError(t, err)

	updatedAccount2, err := testQueries.GetAccount(context.Background(), account2.ID) /// ? why not - store.GetAccount(
	require.NoError(t, err)

	fmt.Println(">> before transaction: ", account1.Balance, account2.Balance)

	require.Equal(t, account1.Balance-int64(n)*ammount, updatedAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*ammount, updatedAccount2.Balance)

}
