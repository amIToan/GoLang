package db

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func createNewUser(t *testing.T) (CreateUserParams, User, error) {
	Username := util.RandomOwner()
	hashedPass, err := util.HashPassword(util.RandomString(6))
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
