package db

import (
	"context"
	"database/sql"
)

func (store *SQLStore) VerifyEmailTx(ctx context.Context, arg UpdateVerifyEmailParams) (*User, *VerifyEmail, error) {
	var user *User
	var mail *VerifyEmail
	err := store.execTx(ctx, func(q *Queries) error {
		verifyEmail, err := store.UpdateVerifyEmail(ctx, arg)
		if err != nil {
			return err
		}
		mail = &verifyEmail
		updatedUser, err := store.UpdateUser(ctx, UpdateUserParams{
			IsEmailVerified: sql.NullBool{
				Bool:  true,
				Valid: true,
			},
			Username: verifyEmail.Username,
		})
		if err != nil {
			return err
		}
		user = &updatedUser
		return nil
	})
	return (user), mail, err
}
