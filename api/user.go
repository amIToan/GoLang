package api

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/lib/pq"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/token"
	"sgithub.com/techschool/simplebank/util"
)

type createUserBody struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
	FullName string `json:"full_name"`
	Email    string `json:"email" binding:"required,email"`
}
type userResponse struct {
	Username          string    `json:"username"`
	FullName          string    `json:"full_name"`
	Email             string    `json:"email"`
	HashedPassword    string    `json:"hashed_password"`
	PasswordChangedAt time.Time `json:"password_changed_at"`
	CreatedAt         time.Time `json:"created_at"`
}

func newUserResponse(user db.User) userResponse {
	return userResponse{
		Username:          user.Username,
		FullName:          user.FullName,
		Email:             user.Email,
		HashedPassword:    user.HashedPassword,
		PasswordChangedAt: user.PasswordChangedAt,
		CreatedAt:         user.CreatedAt,
	}
}
func (server *Server) createUser(ctx *gin.Context) {
	var req createUserBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}
	hashedPass, err := util.HashPassword(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	arg := db.CreateUserParams{
		Username:       req.Username,
		HashedPassword: hashedPass,
		FullName:       req.FullName,
		Email:          req.Email,
	}

	//create account
	user, err := server.store.CreateUser(ctx, arg)

	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errRes(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	rsp := newUserResponse(user)
	ctx.JSON(http.StatusOK, rsp)
}

type userNameUri struct {
	UserName string `uri:"username" binding:"required,alphanum"`
}

func (server *Server) GetUser(ctx *gin.Context) {
	var req userNameUri
	authPayload, ok := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if err := ctx.ShouldBindUri(&req); err != nil {
		log.Println("error in get user", err)
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}
	if !ok || authPayload.Username != req.UserName {
		err := errors.New("Unauthorized")
		ctx.JSON(http.StatusBadRequest, errRes(err))
	}
	user, err := server.store.GetUser(ctx, req.UserName)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	fmt.Println(user)
	ctx.JSON(http.StatusOK, user)
}

// /login
type createLoginBody struct {
	Username string `json:"username" binding:"required,alphanum"`
	Password string `json:"password" binding:"required,min=8"`
}
type loginRes struct {
	SessionID             uuid.UUID    `json:"session_id"`
	Token                 string       `json:"access_token"`
	User                  userResponse `json:"user_info"`
	AccessTokenExpiredAt  time.Time    `json:"access_token_expired_at"`
	RefreshToken          string       `json:"refresh_token"`
	RefreshTokenExpiredAt time.Time    `json:"refresh_token_expired_at"`
}

func (server *Server) Login(ctx *gin.Context) {
	var req createLoginBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}
	user, err := server.store.GetUser(ctx, req.Username)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	err = util.CheckPassword(req.Password, user.HashedPassword)
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(user.Username, util.BankerRole, server.config.ValidDurationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
	}
	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(user.Username, util.BankerRole, server.config.RefreshTokenDurationTime)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
	}
	session, err := server.store.CreateSession(ctx, db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     user.Username,
		RefreshToken: refreshToken,
		UserAgent:    "",
		ClientIp:     "",
		IsBlocked:    false,
		ExpiresAt:    refreshTokenPayload.ExpiredAt,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
	}
	res := loginRes{
		SessionID:             session.ID,
		Token:                 accessToken,
		AccessTokenExpiredAt:  accessTokenPayload.ExpiredAt,
		RefreshToken:          refreshToken,
		RefreshTokenExpiredAt: refreshTokenPayload.ExpiredAt,
		User:                  newUserResponse(user),
	}
	ctx.JSON(http.StatusOK, res)
}

//////
