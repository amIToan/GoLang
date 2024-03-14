package db

import (
	"context"
	"reflect"
)

// Define the CreateUserTxParams interface
type CreateUserTxParamsInterface interface {
	CreateUserParams
	CallbackAfterCreate(user User) error
}

type CreateUserTxParams struct {
	CreateUserParams
	CallbackAfterCreate func(user User) error
}

type CreateUserTxResult struct {
	User
}

func (store *SQLStore) CreateUserTx(ctx context.Context, arg CreateUserTxParams) (CreateUserTxResult, error) {
	var result CreateUserTxResult

	err := store.execTx(ctx, func(q *Queries) error {
		var err error

		result.User, err = q.CreateUser(ctx, arg.CreateUserParams)
		if err != nil {
			return err
		}
		txType := reflect.TypeOf(arg)
		// Iterate over the fields of the struct
		for i := 0; i < txType.NumField(); i++ {
			field := txType.Field(i)
			if field.Name == "CallbackAfterCreate" {
				err = arg.CallbackAfterCreate(result.User)
				break
			}
		}
		return err
	})

	return result, err
}
