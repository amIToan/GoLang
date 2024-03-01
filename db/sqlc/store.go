package db

import (
	"context"
	"database/sql"
	"fmt"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferResultType, error)
}
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
	fmt.Println("gia tri tri trong newStore", db)
	return &SQLStore{
		Queries: New(db),
		db:      db,
	}
}

var (
	txKey struct{} = struct{}{}
)

// execTx executes a function within a database transaction
func (store *SQLStore) execTx(ctx context.Context, callback func(*Queries) error) error {
	tx, err := store.db.BeginTx(ctx, nil)
	fmt.Println("dia chi store exec", store)
	if err != nil {
		return err
	}
	q := New(tx)
	err = callback(q)
	if err != nil {
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx: err: %v, rollback error: %v", err, rbErr)
		}
	}
	return tx.Commit()
}

// transferTx performs a money transfer from one account to the other.
// It creates a transfer record. add account entries , and update accounts' balance within a single database transaction
type TransferResultType struct {
	Transfer
	FromAccount Account
	ToAccount   Account
	FromEntry   Entry
	ToEntry     Entry
}

func (store *SQLStore) TransferTx(ctx context.Context, arg CreateTransferParams) (TransferResultType, error) {
	var result TransferResultType
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
			Amount:    -arg.Amount,
		})
		if err != nil {
			return err
		}
		result.ToEntry, err = q.CreateEntry(ctx, CreateEntryParams{
			AccountID: arg.ToAccountID,
			Amount:    arg.Amount,
		})
		if err != nil {
			return err
		}
		//Todo: It's relating to potential deadlock and locking
		if arg.FromAccountID < arg.ToAccountID {
			result.FromAccount, result.ToAccount, err = addMoneyToAccount(ctx, q, arg.FromAccountID, -arg.Amount, arg.ToAccountID, arg.Amount)
			if err != nil {
				return err
			}
		} else {
			result.ToAccount, result.FromAccount, err = addMoneyToAccount(ctx, q, arg.ToAccountID, arg.Amount, arg.FromAccountID, -arg.Amount)
			if err != nil {
				return err
			}
		}
		return nil
	})
	return (result), err
}
func addMoneyToAccount(ctx context.Context, q *Queries, smallId, amount1, greaterId, amount2 int64) (Account, Account, error) {
	var (
		account1 Account
		account2 Account
		err      error
	)
	account1, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     smallId,
		Amount: amount1,
	})
	account2, err = q.AddAccountBalance(ctx, AddAccountBalanceParams{
		ID:     greaterId,
		Amount: amount2,
	})
	return account1, account2, err
}
