package db

import (
	"context"
	"database/sql"
	"fmt"
)

// Store provides all function to execute db queries and transactions individualy
type Store interface {
	Querier
	TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error)
}

// SQLStore provides all function to execute SQL queries and transactions individualy
type SQLStore struct {
	//*Queries 				- we embed qeries in store and use 	-  store.GetTransfer(...)
	// queries *Queries  	- otherwise 						-  use store.Queries.Get...()
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	return &SQLStore{
		db:      db,
		Queries: New(db),
	}
}

// execTx executes a function within a db transaction - if success save this state or rollback
// to make safe taransaction
// it internal common for all spesific operations
func (store *SQLStore) execTx(ctx context.Context, fn func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil) // later will set isolation level -  tore.db.BeginTx(ctx, &sql.TxOptions{Isolation: sql.IsolationLevel(...)})
	if err != nil {
		return err
	}

	q := New(tx)    // another then store.Queries !!!
	err = fn(q)     // we have queries than will be used in transaction
	if err != nil { // if error - rollback transaction
		if rbErr := tx.Rollback(); rbErr != nil { // if rollback get error too - combine witn main query err
			return fmt.Errorf("tx err: %v, rb err: %v", err, rbErr)
		}
		return err // if rollbeck success - return just main query err
	}
	// if all success
	return tx.Commit()
}

// TransferTxParams contains all nessesary data for transfer money
type TransferTxParams struct {
	FromAccountID int64 `json:"from_accont_id"`
	ToAccountID   int64 `json:"to_accont_id"`
	Amount        int64 `json:"amount"`
}

// TransferTxResult is the result ot transfer transaction
type TransferTxResult struct {
	Transfer    Transfer `json:"transfer"`
	FromAccount Account  `json:"from_account"`
	ToAccount   Account  `json:"to_account"`
	FromEntry   Entry    `json:"from_entry"`
	ToEntry     Entry    `json:"to_entry"`
}

// TransferTx performs a money transfer from one account to another
// It creates a new transfer record, add new account entries, and updates account balanse within a single db transaction
func (store *SQLStore) TransferTx(ctx context.Context, arg TransferTxParams) (TransferTxResult, error) {
	var result TransferTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.Transfer, err = q.CreateTransfer(ctx, CreateTransferParams{
			FromAccountID: arg.FromAccountID,
			ToAccountID:   arg.ToAccountID,
			Amount:        arg.Amount,
		})
		if err != nil {
			return err
		}

		result.FromEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.FromAccountID,
			Amount:    -arg.Amount, // minus amount from
		})
		if err != nil {
			return err
		}

		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount, // plus amont to
		})
		if err != nil {
			return err
		}

		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoney(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
		} else {
			result.ToAccount, result.FromAccount, err = addMoney(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
		}

		return nil
	})

	return result, err
}

func addMoney(
	ctx context.Context,
	q *Queries,
	accountID1 int64,
	amount1 int64,
	accountID2 int64,
	amount2 int64,
) (account1 Account, account2 Account, err error) {
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID1,
		Amount: amount1,
	})
	if err != nil {
		return // same as - return account1 , account2, err
	}

	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     accountID2,
		Amount: amount2,
	})

	return // same as - return account1 , account2, err

}
