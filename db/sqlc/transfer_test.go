package db

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/techschool/simplebank/util"
)

func CreateNewTransfer(t *testing.T) Transfer {
	_, account1, _ := createNewAccount()
	_, account2, _ := createNewAccount()
	require.NotEmpty(t, account1)
	require.NotEmpty(t, account2)
	p := CreateTransferParams{
		FromAccountID: account1.ID,
		ToAccountID:   account2.ID,
		Amount:        util.RandomMoney(),
	}
	i, tranErr := testQueries.CreateTransfer(context.Background(), p)
	require.NoError(t, tranErr)
	require.NotEmpty(t, i)
	return i
}

func TestCreateTransfer(t *testing.T) {
	CreateNewTransfer(t)
}
func TestGetTransfer(t *testing.T) {
	i := CreateNewTransfer(t)
	gotTransfer, err := testQueries.GetTransfer(context.Background(), i.ID)
	require.NoError(t, err)
	require.NotEmpty(t, gotTransfer)
	assert.Equal(t, i.ID, gotTransfer.ID)
	assert.Equal(t, i.Amount, gotTransfer.Amount)
}

func TestGetListTransfer(t *testing.T) {
	_, account1, _ := createNewAccount()
	_, account2, _ := createNewAccount()
	require.NotEmpty(t, account1)
	require.NotEmpty(t, account2)
	p := CreateTransferParams{
		FromAccountID: account2.ID,
		ToAccountID:   account1.ID,
		Amount:        util.RandomMoney(),
	}
	for i := 0; i < 6; i++ {
		i, tranErr := testQueries.CreateTransfer(context.Background(), p)
		require.NoError(t, tranErr)
		require.NotEmpty(t, i)
	}
	pList := ListTransfersParams{
		FromAccountID: account2.ID,
		ToAccountID:   account1.ID,
		Limit:         100,
		Offset:        0,
	}
	list, listErr := testQueries.ListTransfers(context.Background(), pList)
	require.NoError(t, listErr)
	require.NotEmpty(t, list)
	assert.Greater(t, len(list), 0)
	for _, v := range list {
		assert.Equal(t, v.FromAccountID, account2.ID)
		assert.Equal(t, v.ToAccountID, account1.ID)
	}
}
