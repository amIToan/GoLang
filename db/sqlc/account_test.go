package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sgithub.com/techschool/simplebank/util"
)

func createNewAccount(t *testing.T) (CreateAccountParams, Account, error) {
	_, user, _ := createNewUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	account, err := testQueries.CreateAccount(context.Background(), arg)
	return arg, account, err
}
func TestCreateAccount(t *testing.T) {
	arg, account, err := createNewAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	require.Equal(t, arg.Owner, account.Owner)
	require.Equal(t, arg.Balance, account.Balance)
	require.Equal(t, arg.Currency, account.Currency)
	require.NotZero(t, account.ID)
	require.NotZero(t, account.CreatedAt)
}
func TestAddAccountBalance(t *testing.T) {
	_, account, _ := createNewAccount(t)
	balanceParams := AddAccountBalanceParams{
		Amount: util.RandomMoney(),
		ID:     account.ID,
	}
	i, err := testQueries.AddAccountBalance(context.Background(), balanceParams)
	assert.NotNil(t, i)
	assert.NoError(t, err)
	assert.NotEqual(t, account.Balance, i.Balance)
}
func TestGetAccount(t *testing.T) {
	_, account, _ := createNewAccount(t)

	gotAccount, err := testQueries.GetAccount(context.Background(), account.ID)

	require.NoError(t, err)
	require.NotEmpty(t, gotAccount)
	require.Equal(t, account.Owner, gotAccount.Owner)
	require.Equal(t, account.Balance, gotAccount.Balance)
	require.Equal(t, account.Currency, gotAccount.Currency)
	require.Equal(t, account.ID, gotAccount.ID)
}
func TestAccountForUpdate(t *testing.T) {
	_, account, _ := createNewAccount(t)
	gotAccount, err := testQueries.GetAccountForUpdate(context.Background(), account.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotAccount)
	require.Equal(t, account.Owner, gotAccount.Owner)
	require.Equal(t, account.Balance, gotAccount.Balance)
	require.Equal(t, account.Currency, gotAccount.Currency)
	require.Equal(t, account.ID, gotAccount.ID)
}
func TestDeleteAccount(t *testing.T) {
	_, account, _ := createNewAccount(t)
	err := testQueries.DeleteAccount(context.Background(), account.ID)
	require.NoError(t, err)

	_, gotError := testQueries.GetAccount(context.Background(), account.ID)
	assert.Error(t, gotError)
}

func TestUpdateAccount(t *testing.T) {
	_, account, _ := createNewAccount(t)
	updateParams := UpdateAccountParams{
		ID:      account.ID,
		Balance: util.RandomMoney(),
	}
	i, err := testQueries.UpdateAccount(context.Background(), updateParams)
	require.NoError(t, err)

	assert.NotEqual(t, account.Balance, i.Balance)
	assert.Equal(t, i.ID, account.ID)
}
func TestGetListAccounts(t *testing.T) {
	_, user, _ := createNewUser(t)
	arg := CreateAccountParams{
		Owner:    user.Username,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
	ListAccountsParams1 := ListAccountsByNameParams{
		Owner:  user.Username,
		Limit:  10,
		Offset: 0,
	}
	for i := 0; i < 5; i++ {
		testQueries.CreateAccount(context.Background(), arg)
	}
	accounts, err := testQueries.ListAccountsByName(context.Background(), ListAccountsParams1)
	assert.NoError(t, err)
	assert.Greater(t, len(accounts), 0)
	_, user2, _ := createNewUser(t)
	ListAccountsParams2 := ListAccountsByNameParams{
		Owner:  user2.Username,
		Limit:  10,
		Offset: 0,
	}
	accounts2, err2 := testQueries.ListAccountsByName(context.Background(), ListAccountsParams2)
	assert.Nil(t, err2)
	assert.Equal(t, len(accounts2), 0)
}
