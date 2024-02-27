package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sgithub.com/techschool/simplebank/util"
)

func TestCreateEntry(t *testing.T) {
	_, account, err := createNewAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	entryParams := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	i, entryErr := testQueries.CreateEntry(context.Background(), entryParams)
	assert.NoError(t, entryErr)
	require.NotEmpty(t, i)
	assert.Equal(t, i.AccountID, account.ID)
}

func TestGetEntry(t *testing.T) {
	_, account, err := createNewAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	entryParams := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	i, entryErr := testQueries.CreateEntry(context.Background(), entryParams)
	assert.NoError(t, entryErr)
	require.NotEmpty(t, i)
	assert.Equal(t, i.AccountID, account.ID)
	getI, getErr := testQueries.GetEntry(context.Background(), i.ID)
	require.NoError(t, getErr)
	require.NotEmpty(t, getI)
	assert.Equal(t, getI.ID, i.ID)
	assert.Equal(t, getI.Amount, i.Amount)
}
func TestGetListEntries(t *testing.T) {
	_, account, err := createNewAccount(t)
	require.NoError(t, err)
	require.NotEmpty(t, account)
	entryParams := CreateEntryParams{
		AccountID: account.ID,
		Amount:    util.RandomMoney(),
	}
	for i := 0; i < 5; i++ {
		testQueries.CreateEntry(context.Background(), entryParams)
	}
	p := ListEntriesParams{
		AccountID: account.ID,
		Limit:     100,
		Offset:    0,
	}
	items, listEntriesErr := testQueries.ListEntries(context.Background(), p)
	assert.NoError(t, listEntriesErr)
	require.NotEmpty(t, items)
	assert.Greater(t, len(items), 0)
	for _, v := range items {
		assert.Equal(t, account.ID, v.AccountID)
	}

}
