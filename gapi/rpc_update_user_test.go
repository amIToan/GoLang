package gapi

import (
	"context"
	"database/sql"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	mockdb "sgithub.com/techschool/simplebank/db/mock"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/pb"
	"sgithub.com/techschool/simplebank/token"
	"sgithub.com/techschool/simplebank/util"
	mockwk "sgithub.com/techschool/simplebank/worker/mock"
)

type originUpdateUserTxParamsMatcher struct {
	OriginPassword string
	Fullname       string
	Email          string
	Username       string
}

func (expected originUpdateUserTxParamsMatcher) Matches(x interface{}) bool {
	passedArg, ok := x.(db.UpdateUserParams)
	if !ok {
		return false
	}

	err := util.CheckPassword(expected.OriginPassword, passedArg.HashedPassword.String)
	return err == nil
}

func (e originUpdateUserTxParamsMatcher) String() string {
	return fmt.Sprintf("matches hashed password and original password %s", e.OriginPassword)
}

func CheckRealUpdateUserParamsToExpectedParams(updatedParams db.UpdateUserParams, password string) gomock.Matcher {
	return originUpdateUserTxParamsMatcher{
		OriginPassword: password,
		Fullname:       updatedParams.FullName.String,
		Email:          updatedParams.Email.String,
		Username:       updatedParams.Username,
	}
}

func TestUpdateUserAPI(t *testing.T) {
	user, _ := randomUser(t, util.DepositorRole)
	newPassword := "12345678"
	hashedNewPass, err := util.HashPassword(newPassword)
	require.NoError(t, err)
	newFullName := util.RandomOwner()
	newEmail := util.RandomEmail()
	testCases := []struct {
		name          string
		req           *pb.UpdateUserRequest
		buildContext  func(tokenMaker token.Maker, username string, role string, duration time.Duration) (context.Context, *token.Payload)
		buildStubs    func(store *mockdb.MockStore, session *token.Payload)
		checkResponse func(t *testing.T, res *pb.UpdateUserResponse, err error)
	}{
		{
			name: "OK",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				Password: &newPassword,
				Fullname: &newFullName,
				Email:    &newEmail,
			},
			buildContext: func(tokenMaker token.Maker, username string, role string, duration time.Duration) (context.Context, *token.Payload) {
				return newContextWithBearerToken(t, tokenMaker, username, role, duration)
			},
			buildStubs: func(store *mockdb.MockStore, session *token.Payload) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(session.ID)).Times(1).Return(db.Session{
					ID:        session.ID,
					IsBlocked: false,
				}, nil)
				updateSessionParam := db.UpdateSessionParams{
					ID:        session.ID,
					IsBlocked: true,
				}
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(1).Return(user, nil)
				store.EXPECT().UpdateSession(gomock.Any(), gomock.Eq(updateSessionParam)).Times(1).Return(db.Session{
					ID:        session.ID,
					Username:  user.Username,
					IsBlocked: true,
				}, nil)
				updatedUserParams := db.UpdateUserParams{
					HashedPassword: sql.NullString{
						String: hashedNewPass,
						Valid:  true,
					},
					FullName: sql.NullString{
						String: newFullName,
						Valid:  true,
					},
					Email: sql.NullString{
						String: newEmail,
						Valid:  true,
					},
					Username: user.Username,
				}
				store.EXPECT().UpdateUser(gomock.Any(), CheckRealUpdateUserParamsToExpectedParams(updatedUserParams, newPassword)).Times(1).Return(db.User{Username: user.Username, HashedPassword: newPassword, FullName: newFullName, Email: newEmail, CreatedAt: user.CreatedAt}, nil)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.NoError(t, err)
				require.NotNil(t, res)
				st, ok := status.FromError(err)
				require.Equal(t, codes.OK, st.Code())
				require.True(t, ok)
				updatedUser := res.GetUser()
				require.Equal(t, user.Username, updatedUser.Username)
				require.Equal(t, newFullName, updatedUser.FullName)
				require.Equal(t, newEmail, updatedUser.Email)
			},
		},
		{
			name: "Unauthenticated",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				Password: &newPassword,
				Fullname: &newFullName,
				Email:    &newEmail,
			},
			buildContext: func(tokenMaker token.Maker, username string, role string, duration time.Duration) (context.Context, *token.Payload) {
				return newContextWithBearerToken(t, tokenMaker, "wrong user name", role, duration)
			},
			buildStubs: func(store *mockdb.MockStore, session *token.Payload) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(session.ID)).Times(1).Return(db.Session{
					ID:        session.ID,
					IsBlocked: false,
				}, nil)
				updateSessionParam := db.UpdateSessionParams{
					ID:        session.ID,
					IsBlocked: true,
				}
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(0)
				store.EXPECT().UpdateSession(gomock.Any(), gomock.Eq(updateSessionParam)).Times(0)
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "SessionBlocked",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				Password: &newPassword,
				Fullname: &newFullName,
				Email:    &newEmail,
			},
			buildContext: func(tokenMaker token.Maker, username string, role string, duration time.Duration) (context.Context, *token.Payload) {
				return newContextWithBearerToken(t, tokenMaker, username, role, duration)
			},
			buildStubs: func(store *mockdb.MockStore, session *token.Payload) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(session.ID)).Times(1).Return(db.Session{
					ID:        session.ID,
					IsBlocked: true,
				}, nil)
				updateSessionParam := db.UpdateSessionParams{
					ID:        session.ID,
					IsBlocked: true,
				}
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(0)
				store.EXPECT().UpdateSession(gomock.Any(), gomock.Eq(updateSessionParam)).Times(0)
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
		{
			name: "InvalidArgument",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				Password: StringPtr("123"),
				Fullname: StringPtr("ahhihihi"),
				Email:    StringPtr("23456789"),
			},
			buildContext: func(tokenMaker token.Maker, username string, role string, duration time.Duration) (context.Context, *token.Payload) {
				return newContextWithBearerToken(t, tokenMaker, username, role, duration)
			},
			buildStubs: func(store *mockdb.MockStore, session *token.Payload) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(session.ID)).Times(1).Return(db.Session{
					ID:        session.ID,
					IsBlocked: false,
				}, nil)
				updateSessionParam := db.UpdateSessionParams{
					ID:        session.ID,
					IsBlocked: true,
				}
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(0)
				store.EXPECT().UpdateSession(gomock.Any(), gomock.Eq(updateSessionParam)).Times(0)
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.InvalidArgument, st.Code())
			},
		},
		{
			name: "Expired Token",
			req: &pb.UpdateUserRequest{
				Username: user.Username,
				Password: &newPassword,
				Fullname: &newFullName,
				Email:    &newEmail,
			},
			buildContext: func(tokenMaker token.Maker, username string, role string, duration time.Duration) (context.Context, *token.Payload) {
				return newContextWithBearerToken(t, tokenMaker, username, role, 0*time.Microsecond)
			},
			buildStubs: func(store *mockdb.MockStore, session *token.Payload) {
				store.EXPECT().GetSession(gomock.Any(), gomock.Eq(session.ID)).Times(1).Return(db.Session{
					ID:        session.ID,
					IsBlocked: false,
				}, nil)
				store.EXPECT().GetUser(gomock.Any(), gomock.Eq(user.Username)).Times(0)
				store.EXPECT().UpdateSession(gomock.Any(), gomock.Any()).Times(0)
				store.EXPECT().UpdateUser(gomock.Any(), gomock.Any()).Times(0)
			},
			checkResponse: func(t *testing.T, res *pb.UpdateUserResponse, err error) {
				require.Error(t, err)
				require.Empty(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, codes.Unauthenticated, st.Code())
			},
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			storeCtrl := gomock.NewController(t)
			defer storeCtrl.Finish()
			store := mockdb.NewMockStore(storeCtrl)
			taskCtrl := gomock.NewController(t)
			defer taskCtrl.Finish()
			taskDistributor := mockwk.NewMockTaskDistributor(taskCtrl)
			server := newTestServer(t, store, taskDistributor)
			metadata, refreshPayload := tc.buildContext(server.tokenMaker, user.Username, util.BankerRole, server.config.ValidDurationTime)
			tc.buildStubs(store, refreshPayload)
			updatedUser, err := server.UpdateUser(metadata, tc.req)
			tc.checkResponse(t, updatedUser, err)
		})
	}
}

// Helper function to create a pointer to string
func StringPtr(s string) *string {
	return &s
}

// Helper function to create a pointer to bool
func BoolPtr(b bool) *bool {
	return &b
}
