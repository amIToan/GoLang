package db

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestTransferTx(t *testing.T) {
	store := NewStore(testDB)
	_, account1, _ := createNewAccount()
	_, account2, _ := createNewAccount()
	fmt.Println(">>>>before : ", account1.Balance, account2.Balance)
	n := 2
	amount := int64(100)
	errs := make(chan error)
	results := make(chan TransferResultType)
	txNames := make(chan string)
	for i := 0; i < n; i++ {
		txName := fmt.Sprintf("tx %d", i+1)
		fmt.Printf("khoi tao %s \n", txName)
		go func() {
			ctx := context.WithValue(context.Background(), txKey, txName)
			result, err := store.transferTx(ctx, CreateTransferParams{
				FromAccountID: account1.ID,
				ToAccountID:   account2.ID,
				Amount:        amount,
			})
			errs <- err
			results <- result
			txNames <- txName
		}()
	}
	//check err
	existed := make(map[int]bool)
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)
		result := <-results
		require.NotEmpty(t, result)
		txName := <-txNames
		require.NotEmpty(t, txName)
		//check transfer
		transfer := result.Transfer
		require.Equal(t, transfer.FromAccountID, account1.ID)
		require.Equal(t, transfer.ToAccountID, account2.ID)
		require.Equal(t, transfer.Amount, amount)
		require.NotNil(t, transfer.CreatedAt)
		_, err = store.GetTransfer(context.Background(), transfer.ID)
		require.NoError(t, err)
		//check from and to entries
		fromEntry := result.FromEntry
		require.NotEmpty(t, fromEntry)
		require.NotZero(t, fromEntry.ID)
		require.NotZero(t, fromEntry.CreatedAt)
		require.Equal(t, fromEntry.AccountID, transfer.FromAccountID)
		require.Equal(t, fromEntry.Amount, -amount)
		_, err = store.GetEntry(context.Background(), fromEntry.ID)
		require.NoError(t, err)
		toEntry := result.ToEntry
		require.NotEmpty(t, toEntry)
		require.NotZero(t, toEntry.ID)
		require.NotZero(t, toEntry.CreatedAt)
		require.Equal(t, toEntry.AccountID, transfer.ToAccountID)
		require.Equal(t, toEntry.Amount, amount)
		_, err = store.GetEntry(context.Background(), toEntry.ID)
		require.NoError(t, err)
		//todo : update balance
		fromAccount := result.FromAccount
		require.NotEmpty(t, fromAccount)
		require.Equal(t, fromAccount.ID, account1.ID)

		toAccount := result.ToAccount
		require.NotEmpty(t, toAccount)
		require.Equal(t, toAccount.ID, account2.ID)
		//check balance
		fmt.Println(">>>>transaction : ", txName, fromAccount.Balance, toAccount.Balance)
		diff1 := account1.Balance - fromAccount.Balance // money is banked
		fmt.Println(">>>>diff1 : ", txName, diff1)
		diff2 := toAccount.Balance - account2.Balance // money is added
		fmt.Println(">>>>diff2 : ", txName, diff2)
		require.Equal(t, diff1, diff2)
		require.True(t, diff1 > 0)
		require.True(t, diff1%amount == 0) // for each turn has to be divisible

		k := int(diff1 / amount)
		require.True(t, k >= 1 && k <= n)
		require.NotContains(t, existed, k)
		existed[k] = true
	}
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount1)
	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount2)
	require.Equal(t, account1.Balance-int64(n)*amount, updateAccount1.Balance)
	require.Equal(t, account2.Balance+int64(n)*amount, updateAccount2.Balance)
}

func TestTransferDeadLockTx(t *testing.T) {
	store := NewStore(testDB)
	_, account1, _ := createNewAccount()
	_, account2, _ := createNewAccount()
	fmt.Println(">>>>before : ", account1.Balance, account2.Balance)
	n := 10
	amount := int64(100)
	errs := make(chan error)
	for i := 0; i < n; i++ {
		idFrom := account1.ID
		idTo := account2.ID
		if i%2 == 1 {
			idFrom = account2.ID
			idTo = account1.ID
		}
		go func() {
			_, err := store.transferTx(context.Background(), CreateTransferParams{
				FromAccountID: idFrom,
				ToAccountID:   idTo,
				Amount:        amount,
			})
			errs <- err
		}()
	}
	//check err
	for i := 0; i < n; i++ {
		err := <-errs
		require.NoError(t, err)

	}
	updateAccount1, err := testQueries.GetAccount(context.Background(), account1.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount1)
	updateAccount2, err := testQueries.GetAccount(context.Background(), account2.ID)
	require.NoError(t, err)
	require.NotEmpty(t, updateAccount2)
	require.Equal(t, account1.Balance, updateAccount1.Balance)
	require.Equal(t, account2.Balance, updateAccount2.Balance)
}
