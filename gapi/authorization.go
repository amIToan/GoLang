package gapi

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/google/uuid"
	"google.golang.org/grpc/metadata"
	"sgithub.com/techschool/simplebank/token"
)

const (
	authorizationHeaderKey = "authorization"
	session                = "session_id"
)

var authorizationHeaderTypes = [...]string{"bearer"}

func (server *Server) authorization(ctx context.Context, accessibleRoles []string) (*token.Payload, uuid.UUID, error) {
	//gRPC client
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, uuid.UUID{}, fmt.Errorf("missing metadata")
	}
	session := md.Get(session)
	if len(session) == 0 {
		return nil, uuid.UUID{}, fmt.Errorf("missing session_id")
	}
	sessionString := strings.ToLower(session[0])
	// Parse the UUID string into a uuid.UUID
	parsedUUID, err := uuid.Parse(sessionString)
	if err != nil {
		return nil, uuid.UUID{}, fmt.Errorf("session is wrong")
	}
	userSession, err := server.store.GetSession(ctx, parsedUUID)
	if err != nil {
		return nil, uuid.UUID{}, fmt.Errorf("session isn't existed")
	}
	if userSession.IsBlocked {
		return nil, uuid.UUID{}, fmt.Errorf("session is blocked")
	}
	authorizationHeader := md.Get(authorizationHeaderKey)
	if len(authorizationHeader) == 0 {
		return nil, uuid.UUID{}, fmt.Errorf("missing authorization header")
	}
	authString := authorizationHeader[0]
	fields := strings.Fields(authString)
	if len(fields) < 2 {
		err := errors.New("authorization header format is not supported")
		return nil, uuid.UUID{}, err
	}
	authorizationHeaderType := strings.ToLower(fields[0])
	// check is supported or not.
	i := sort.Search(len(authorizationHeaderTypes), func(i int) bool { return authorizationHeaderType == authorizationHeaderTypes[i] })
	if i < len(authorizationHeaderTypes) && authorizationHeaderTypes[i] == authorizationHeaderType {
		fmt.Printf("Found %s at index %d in %v.\n", authorizationHeaderType, i, authorizationHeaderTypes)
	} else {
		fmt.Printf("Did not find %s in %v.\n", authorizationHeaderType, authorizationHeaderTypes)
		err := errors.New("authorization header type is not supported")
		return nil, uuid.UUID{}, err
	}

	accessToken := fields[1]
	payloadFromToken, err := server.tokenMaker.VerifyToken(accessToken)
	if err != nil {
		return nil, uuid.UUID{}, err
	}
	if !hasPermissionOrNot(payloadFromToken.Role, accessibleRoles) {
		return nil, uuid.UUID{}, fmt.Errorf("not allowed")
	}
	return payloadFromToken, parsedUUID, nil
}
func hasPermissionOrNot(Role string, accessibleRoles []string) bool {
	for _, role := range accessibleRoles {
		if role == Role {
			return true
		}
	}
	return false
}
