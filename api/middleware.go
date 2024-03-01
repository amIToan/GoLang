package api

import (
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"sgithub.com/techschool/simplebank/token"
)

const (
	authorizationHeaderKey  = "authorization"
	authorizationPayloadKey = "authorization_payload"
)

var authorizationHeaderTypes = [...]string{"bearer"}

func authMiddleWare(tokenMaker token.Maker) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authorizationHeader := ctx.GetHeader(authorizationHeaderKey)
		if len(authorizationHeader) == 0 {
			err := errors.New("authorization header is not supplied")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errRes(err))
			return
		}
		fields := strings.Fields(authorizationHeader)
		if len(fields) < 2 {
			err := errors.New("authorization header format is not supported")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errRes(err))
			return
		}

		authorizationHeaderType := strings.ToLower(fields[0])
		// check is supported or not.
		i := sort.Search(len(authorizationHeaderTypes), func(i int) bool { return authorizationHeaderType == authorizationHeaderTypes[i] })
		if i < len(authorizationHeaderTypes) && authorizationHeaderTypes[i] == authorizationHeaderType {
			fmt.Printf("Found %s at index %d in %v.\n", authorizationHeaderType, i, authorizationHeaderTypes)
		} else {
			fmt.Printf("Did not find %s in %v.\n", authorizationHeaderType, authorizationHeaderTypes)
			err := errors.New("authorization header type is not supported")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errRes(err))
			return
		}
		/* if idx < 0 || (strings.ToLower(authorizationHeaderTypes[idx]) != authorizationHeaderType) {
			err := errors.New("authorization header type is not supported")
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errRes(err))
			return
		} */
		accessToken := fields[1]
		payloadFromToken, err := tokenMaker.VerifyToken(accessToken)
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, errRes(err))
		}
		ctx.Set(authorizationPayloadKey, payloadFromToken)
	}
}
