package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    "Vasya",
		Balance:  100,
		Currency: "USD",
	}

	accont, err := testQueries.CreateAccount(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, accont)

	require.Equal(t, arg.Owner, accont.Owner)
	require.Equal(t, arg.Balance, accont.Balance)
	require.Equal(t, arg.Currency, accont.Currency)

	require.NotZero(t, accont.ID)
	require.NotZero(t, accont.CreatedAt)
}
