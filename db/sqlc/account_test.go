package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
	"parus.i234.me/paultabaco/simplebank/util"
)

func TestCreateAccount(t *testing.T) {
	arg := CreateAccountParams{
		Owner:    util.RandomOwner(),
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
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
