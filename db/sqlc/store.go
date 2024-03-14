package db

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/rs/zerolog/log"
)

type Store interface {
	Querier
	TransferTx(ctx context.Context, arg CreateTransferParams) (TransferResultType, error)
	CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error)
	VerifyEmailTx(ctx context.Context, arg UpdateVerifyEmailParams) (*User, *VerifyEmail, error)
}
type SQLStore struct {
	*Queries
	db *sql.DB
}

func NewStore(db *sql.DB) Store {
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
	if err != nil {
		return err
	}
	q := New(tx)
	err = callback(q)
	if err != nil {
		log.Error().Err(err).Msg("Error in transaction")
		if rbErr := tx.Rollback(); rbErr != nil {
			return fmt.Errorf("tx: err: %v, rollback error: %v", err, rbErr)
		}
	}
	return tx.Commit()
}
