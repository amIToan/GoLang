package api

import (
	"database/sql"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/token"
)

type createAccountBody struct {
	Owner    string `json:"owner" binding:"required"`
	Currency string `json:"currency" binding:"required,currency"`
}

func (server *Server) createAccount(ctx *gin.Context) {
	var req createAccountBody
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}
	arg := db.CreateAccountParams{
		Owner:    authPayload.Username,
		Currency: req.Currency,
		Balance:  0,
	}

	//create account
	account, err := server.store.CreateAccount(ctx, arg)

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
	ctx.JSON(http.StatusOK, account)
}

type accountUri struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getAccountById(ctx *gin.Context) {
	var req accountUri
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	account, err := server.store.GetAccountForUpdate(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}

	ctx.JSON(http.StatusOK, account)
}

type accountListBody struct {
	Limit int32 `form:"limit" binding:"required,min=10,max=50"`
	Page  int32 `form:"page" binding:"required,min=1"`
}

func (server *Server) getAccountsList(ctx *gin.Context) {
	var req accountListBody

	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	arg := db.ListAccountsByNameParams{
		Owner:  authPayload.Username,
		Limit:  req.Limit,
		Offset: (req.Page - 1) * req.Limit,
	}
	accounts, err := server.store.ListAccountsByName(ctx, arg)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	ctx.JSON(http.StatusOK, accounts)
}

type UpdateAccountForm struct {
	ID      int64 `form:"id" binding:"required,min=1"`
	Balance int64 `form:"balance" binding:"required,min=1000"`
}

func (server *Server) updateAccountById(ctx *gin.Context) {
	var form UpdateAccountForm
	// This will infer what binder to use depending on the content-type header.
	if err := ctx.ShouldBind(&form); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	arg := db.UpdateAccountParams{
		ID:      form.ID,
		Balance: form.Balance,
	}

	account, _ := server.store.GetAccountForUpdate(ctx, form.ID)
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}
	updatedAccount, err := server.store.UpdateAccount(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	ctx.JSON(http.StatusOK, updatedAccount)
}

type DeleteAccountParam struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) DeleteAccountById(ctx *gin.Context) {
	var req DeleteAccountParam
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}
	account, err := server.store.GetAccountForUpdate(ctx, req.ID)
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	if account.Owner != authPayload.Username {
		err := errors.New("account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}
	err = server.store.DeleteAccount(ctx, req.ID)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return
	}
	ctx.JSON(http.StatusOK, struct {
		Status  string
		Message string
	}{
		Status:  "successful",
		Message: "Delete successfully",
	})
}
