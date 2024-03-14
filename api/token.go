package api

import (
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"sgithub.com/techschool/simplebank/util"
)

// /login
type createRefreshTokenBody struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}
type renewedTokenRes struct {
	Token                string    `json:"access_token"`
	AccessTokenExpiredAt time.Time `json:"access_token_expired_at"`
}

func (server *Server) renewAccessToken(ctx *gin.Context) {
	var req createRefreshTokenBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}
	refreshTokenPayload, err := server.tokenMaker.VerifyToken(req.RefreshToken)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}
	dbRefreshToken, err := server.store.GetSession(ctx, refreshTokenPayload.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	if dbRefreshToken.IsBlocked {
		err := fmt.Errorf("session is blocked")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}
	if (dbRefreshToken.RefreshToken != req.RefreshToken) || (dbRefreshToken.Username != refreshTokenPayload.Username) {
		err := fmt.Errorf("wrong token")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}
	if time.Now().After(dbRefreshToken.ExpiresAt) {
		err := fmt.Errorf("token is expired")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}
	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(dbRefreshToken.Username, util.BankerRole, server.config.ValidDurationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
	}

	res := renewedTokenRes{
		Token:                accessToken,
		AccessTokenExpiredAt: accessTokenPayload.ExpiredAt,
	}
	ctx.JSON(http.StatusOK, res)
}
