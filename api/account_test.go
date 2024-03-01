package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	mockdb "sgithub.com/techschool/simplebank/db/mock"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/token"
	"sgithub.com/techschool/simplebank/util"
)

func TestGetAccountApi(t *testing.T) {
	user, _ := randomUser(t)
	account := randomAccount(user.Username)
	testCases := []struct {
		name          string
		accountID     int64
		setupAuth     func(t *testing.T, request *http.Request, tokenMaker token.Maker)
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationHeaderTypes[0], user.Username, util.DepositorRole, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyHaveToMatch(t, recorder.Body, account)
			},
		},
		{
			name:      "UnauthorizedUser",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationHeaderTypes[0], "unauthorized_user", util.DepositorRole, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).
					Times(1).
					Return(account, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "NoAuthorization",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().
					GetAccountForUpdate(gomock.Any(), gomock.Any()).
					Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusUnauthorized, recorder.Code)
			},
		},
		{
			name:      "Not Found",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationHeaderTypes[0], user.Username, util.DepositorRole, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
		{
			name:      "InternalError",
			accountID: account.ID,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationHeaderTypes[0], user.Username, util.DepositorRole, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrConnDone)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusInternalServerError, recorder.Code)
			},
		},
		{
			name:      "Invalid ID",
			accountID: 0,
			setupAuth: func(t *testing.T, request *http.Request, tokenMaker token.Maker) {
				addAuthorization(t, request, tokenMaker, authorizationHeaderTypes[0], user.Username, util.DepositorRole, time.Minute)
			},
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountForUpdate(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			v.buildStubs(store)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/account/%d", v.accountID)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			v.setupAuth(t, request, server.tokenMaker)
			server.router.ServeHTTP(recorder, request)
			require.NoError(t, err)
			v.checkResponse(t, recorder)
		})
	}
}

func TestGetAccountList(t *testing.T) {
	user, _ := randomUser(t)
	accountListRes := []db.Account{
		{
			ID:      util.RandomInt(1, 1000),
			Owner:   user.Username,
			Balance: util.RandomMoney(),
		},
		{
			ID:      util.RandomInt(1, 1000),
			Owner:   user.Username,
			Balance: util.RandomMoney(),
		},
		{
			ID:      util.RandomInt(1, 1000),
			Owner:   user.Username,
			Balance: util.RandomMoney(),
		},
		{
			ID:      util.RandomInt(1, 1000),
			Owner:   user.Username,
			Balance: util.RandomMoney(),
		},
		{
			ID:      util.RandomInt(1, 1000),
			Owner:   user.Username,
			Balance: util.RandomMoney(),
		},
		{
			ID:      util.RandomInt(1, 1000),
			Owner:   user.Username,
			Balance: util.RandomMoney(),
		},
		{
			ID:      util.RandomInt(1, 1000),
			Owner:   user.Username,
			Balance: util.RandomMoney(),
		},
		{
			ID:      util.RandomInt(1, 1000),
			Owner:   user.Username,
			Balance: util.RandomMoney(),
		},
		{
			ID:      util.RandomInt(1, 1000),
			Owner:   user.Username,
			Balance: util.RandomMoney(),
		},
	}
	testCases := []struct {
		name          string
		page          int64
		limit         int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:  "OK",
			page:  1,
			limit: 20,
			buildStubs: func(store *mockdb.MockStore) {
				arg := db.ListAccountsByNameParams{
					Owner:  user.Username,
					Limit:  20,
					Offset: 0,
				}
				store.EXPECT().ListAccountsByName(gomock.Any(), gomock.Eq(arg)).Times(1).Return(accountListRes, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				var gotAccount []db.Account
				err = json.Unmarshal(data, &gotAccount)
				require.NoError(t, err)
				require.Equal(t, gotAccount, accountListRes)
			},
		},
		{
			name:  "Invalid",
			page:  0,
			limit: 5,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().ListAccounts(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			v.buildStubs(store)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			url := fmt.Sprintf("/accounts?limit=%d&page=%d", v.limit, v.page)
			request, err := http.NewRequest(http.MethodGet, url, nil)
			token, payload, err := server.tokenMaker.CreateToken(user.Username, "", time.Minute)
			require.NoError(t, err)
			require.NotEmpty(t, payload)
			authorizationHeader := fmt.Sprintf("%s %s", authorizationHeaderTypes[0], token)
			request.Header.Set(authorizationHeaderKey, authorizationHeader)
			server.router.ServeHTTP(recorder, request)
			require.NoError(t, err)
			v.checkResponse(t, recorder)
		})
	}
}

func TestUpdateAccount(t *testing.T) {
	user1, _ := randomUser(t)
	account := randomAccount(user1.Username)
	testCases := []struct {
		name          string
		accountID     int64
		balance       int64
		buildStubs    func(store *mockdb.MockStore, updatedAccount *db.Account)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder, updatedAccount *db.Account)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			balance:   10000,
			buildStubs: func(store *mockdb.MockStore, updatedAccount *db.Account) {
				arg := db.UpdateAccountParams{
					ID:      account.ID,
					Balance: int64(10000),
				}
				*updatedAccount = account
				updatedAccount.Balance = arg.Balance
				store.EXPECT().GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(*updatedAccount, nil)
				store.EXPECT().UpdateAccount(gomock.Any(), gomock.Eq(arg)).Times(1).Return(*updatedAccount, nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, updatedAccount *db.Account) {
				require.Equal(t, http.StatusOK, recorder.Code)
				requireBodyHaveToMatch(t, recorder.Body, *updatedAccount)
			},
		},
		{
			name:      "Invalid",
			accountID: 0,
			balance:   0,
			buildStubs: func(store *mockdb.MockStore, updatedAccount *db.Account) {
				store.EXPECT().UpdateAccount(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder, updatedAccount *db.Account) {
				require.Equal(t, recorder.Code, http.StatusBadRequest)
			},
		},
	}
	for _, v := range testCases {
		t.Run(v.name, func(t *testing.T) {
			var updatedAccount db.Account
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			store := mockdb.NewMockStore(ctrl)
			v.buildStubs(store, &updatedAccount)
			server := newTestServer(t, store)
			recorder := httptest.NewRecorder()
			// Create a JSON string with dynamic id and balance
			jsonBody := fmt.Sprintf(`{"id": %d, "balance": %d}`, v.accountID, v.balance)
			newBufferJsonBody := bytes.NewBufferString(jsonBody)
			request, err := http.NewRequest(http.MethodPut, "/account/update", newBufferJsonBody)
			require.NoError(t, err)
			// Set the content type header
			token, payload, err := server.tokenMaker.CreateToken(user1.Username, "", time.Minute)
			require.NoError(t, err)
			require.NotEmpty(t, payload)
			authorizationHeader := fmt.Sprintf("%s %s", authorizationHeaderTypes[0], token)
			request.Header.Set(authorizationHeaderKey, authorizationHeader)
			request.Header.Set("Content-Type", "application/json")
			require.NoError(t, err)
			server.router.ServeHTTP(recorder, request)
			require.NoError(t, err)
			v.checkResponse(t, recorder, &updatedAccount)
		})
	}
}
func TestDeleteAccountAPI(t *testing.T) {
	user1, _ := randomUser(t)
	account := randomAccount(user1.Username)
	testCases := []struct {
		name          string
		accountID     int64
		buildStubs    func(store *mockdb.MockStore)
		checkResponse func(t *testing.T, recorder *httptest.ResponseRecorder)
	}{
		{
			name:      "OK",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(account, nil)
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(nil)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusOK, recorder.Code)
				data, err := ioutil.ReadAll(recorder.Body)
				require.NoError(t, err)
				var body struct {
					Status  string
					Message string
				}
				err = json.Unmarshal(data, &body)
				require.NoError(t, err)
				require.Equal(t, body, struct {
					Status  string
					Message string
				}{
					Status:  "successful",
					Message: "Delete successfully",
				})
			},
		},
		{
			name:      "Invalid",
			accountID: 0,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).Times(0)
				store.EXPECT().GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).Times(0)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusBadRequest, recorder.Code)
			},
		}, {
			name:      "Not Found",
			accountID: account.ID,
			buildStubs: func(store *mockdb.MockStore) {
				store.EXPECT().DeleteAccount(gomock.Any(), gomock.Eq(account.ID)).Times(0)
				store.EXPECT().GetAccountForUpdate(gomock.Any(), gomock.Eq(account.ID)).Times(1).Return(db.Account{}, sql.ErrNoRows)
			},
			checkResponse: func(t *testing.T, recorder *httptest.ResponseRecorder) {
				require.Equal(t, http.StatusNotFound, recorder.Code)
			},
		},
	}
	for _, v := range testCases {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		store := mockdb.NewMockStore(ctrl)
		v.buildStubs(store)
		server := newTestServer(t, store)
		recorder := httptest.NewRecorder()
		// Create a JSON string with dynamic id and balance
		url := fmt.Sprintf("/account/delete/%d", v.accountID)
		request, err := http.NewRequest(http.MethodDelete, url, nil)
		require.NoError(t, err)
		require.NoError(t, err)
		// Set the content type header
		token, payload, err := server.tokenMaker.CreateToken(user1.Username, "", time.Minute)
		require.NoError(t, err)
		require.NotEmpty(t, payload)
		authorizationHeader := fmt.Sprintf("%s %s", authorizationHeaderTypes[0], token)
		request.Header.Set(authorizationHeaderKey, authorizationHeader)
		server.router.ServeHTTP(recorder, request)
		require.NoError(t, err)
	}
}
func requireBodyHaveToMatch(t *testing.T, body *bytes.Buffer, account db.Account) {
	data, err := ioutil.ReadAll(body)
	require.NoError(t, err)
	var gotAccount db.Account
	err = json.Unmarshal(data, &gotAccount)
	require.NoError(t, err)
	require.Equal(t, gotAccount, account)
}

func randomAccount(owner string) db.Account {
	return db.Account{
		ID:       util.RandomInt(1, 10000),
		Owner:    owner,
		Balance:  util.RandomMoney(),
		Currency: util.RandomCurrency(),
	}
}
