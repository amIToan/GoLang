package gapi

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/metadata"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/token"
	"sgithub.com/techschool/simplebank/util"
	"sgithub.com/techschool/simplebank/worker"
)

func newTestServer(t *testing.T, store db.Store, taskDistributor worker.TaskDistributor) *Server {
	config := util.Config{
		TokenSymmetricKey: util.RandomStr(32),
		ValidDurationTime: time.Minute,
	}

	server, err := NewServer(config, store, taskDistributor)
	require.NoError(t, err)

	return server
}

func newContextWithBearerToken(t *testing.T, tokenMaker token.Maker, username string, role string, duration time.Duration) (context.Context, *token.Payload) {
	accessToken, _, err := tokenMaker.CreateToken(username, role, duration)
	require.NoError(t, err)
	var refresh_duration time.Duration = 24 * time.Hour
	_, refreshTokenPayload, _ := tokenMaker.CreateToken(username, role, refresh_duration)
	bearerToken := fmt.Sprintf("%s %s", authorizationHeaderTypes[0], accessToken)
	md := metadata.MD{
		authorizationHeaderKey: []string{
			bearerToken,
		},
		"session_id": []string{refreshTokenPayload.ID.String()},
	}

	return metadata.NewIncomingContext(context.Background(), md), refreshTokenPayload
}
