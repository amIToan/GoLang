package db

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"sgithub.com/techschool/simplebank/util"
)

func createNewUser(t *testing.T) (CreateUserParams, User, error) {
	Username := util.RandomOwner()
	hashedPass, err := util.HashPassword(util.RandomStr(6))
	require.NoError(t, err)
	Email := fmt.Sprintf("%s@email.com", Username)
	arg := CreateUserParams{
		Username:       Username,
		HashedPassword: hashedPass,
		FullName:       util.RandomOwner(),
		Email:          Email,
	}
	account, err := testQueries.CreateUser(context.Background(), arg)
	return arg, account, err
}
func TestCreateUser(t *testing.T) {
	arg, account, err := createNewUser(t)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Username, account.Username)
	require.Equal(t, arg.HashedPassword, account.HashedPassword)
	require.Equal(t, arg.Email, account.Email)
	require.NotZero(t, account.CreatedAt)
	require.Zero(t, account.PasswordChangedAt)
}

func TestGetUser(t *testing.T) {
	_, user, err := createNewUser(t)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	gotUser, err := testQueries.GetUser(context.Background(), user.Username)
	require.NoError(t, err)
	require.NotEmpty(t, gotUser)
	require.WithinDuration(t, user.CreatedAt, gotUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, gotUser.PasswordChangedAt, time.Second)
	require.Equal(t, user.Username, gotUser.Username)
	require.Equal(t, user.Email, gotUser.Email)
	require.Equal(t, user.FullName, gotUser.FullName)
}
func TestUpdateUser(t *testing.T) {
	_, user, err := createNewUser(t)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	hashedPass, err := util.HashPassword(util.RandomStr(6))
	require.NoError(t, err)
	arg := UpdateUserParams{
		HashedPassword: sql.NullString{
			String: hashedPass,
			Valid:  true,
		},
		FullName: sql.NullString{
			String: "Ta Quang Toan",
			Valid:  true,
		},
		Email: sql.NullString{
			String: "taquangotandz@gmail.com",
			Valid:  true,
		},
		Username: user.Username,
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.WithinDuration(t, user.CreatedAt, updatedUser.CreatedAt, time.Second)
	require.NotEqual(t, user.PasswordChangedAt, updatedUser.PasswordChangedAt, time.Second)
	require.Equal(t, arg.Username, updatedUser.Username)
	require.Equal(t, arg.Email.String, updatedUser.Email)
	require.Equal(t, arg.FullName.String, updatedUser.FullName)
}

func TestUpdateUserEmail(t *testing.T) {
	_, user, err := createNewUser(t)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NoError(t, err)
	arg := UpdateUserParams{
		Email: sql.NullString{
			String: "taquangotandz1998@gmail.com",
			Valid:  true,
		},
		Username: user.Username,
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.WithinDuration(t, user.CreatedAt, updatedUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, updatedUser.PasswordChangedAt, time.Second)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, arg.Email.String, updatedUser.Email)
	require.Equal(t, user.FullName, updatedUser.FullName)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
}

func TestUpdateUserHashPassword(t *testing.T) {
	_, user, err := createNewUser(t)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NoError(t, err)
	hashedPass, err := util.HashPassword(util.RandomStr(6))
	require.NoError(t, err)
	arg := UpdateUserParams{
		HashedPassword: sql.NullString{
			String: hashedPass,
			Valid:  true,
		},
		Username: user.Username,
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.WithinDuration(t, user.CreatedAt, updatedUser.CreatedAt, time.Second)
	require.NotEqual(t, user.PasswordChangedAt, updatedUser.PasswordChangedAt, time.Second)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, user.FullName, updatedUser.FullName)
	require.Equal(t, arg.HashedPassword.String, updatedUser.HashedPassword)
}

func TestUpdateUserFullName(t *testing.T) {
	_, user, err := createNewUser(t)
	require.NoError(t, err)
	require.NotEmpty(t, user)
	require.NoError(t, err)
	arg := UpdateUserParams{
		FullName: sql.NullString{
			String: "Toan day ne",
			Valid:  true,
		},
		Username: user.Username,
	}
	updatedUser, err := testQueries.UpdateUser(context.Background(), arg)
	require.NoError(t, err)
	require.NotEmpty(t, updatedUser)
	require.WithinDuration(t, user.CreatedAt, updatedUser.CreatedAt, time.Second)
	require.WithinDuration(t, user.PasswordChangedAt, updatedUser.PasswordChangedAt, time.Second)
	require.Equal(t, user.Username, updatedUser.Username)
	require.Equal(t, user.Email, updatedUser.Email)
	require.Equal(t, arg.FullName.String, updatedUser.FullName)
	require.Equal(t, user.HashedPassword, updatedUser.HashedPassword)
}
