package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/token"
)

type transferBody struct {
	FromAccountID int64  `json:"from_account_id" binding:"required,min=1"`
	ToAccountID   int64  `json:"to_account_id" binding:"required,min=1"`
	Amount        int64  `json:"amount" binding:"required,gt=0"`
	Currency      string `json:"currency" binding:"required,currency"`
}

func (server *Server) CreateTransfer(ctx *gin.Context) {
	var req transferBody
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	fromAccount, err := server.checkIsValidCurrencyOfAccount(ctx, req.FromAccountID, req.Currency)
	if !err {
		return
	}
	if fromAccount.Owner != authPayload.Username {
		err := errors.New("from account doesn't belong to the authenticated user")
		ctx.JSON(http.StatusUnauthorized, errRes(err))
		return
	}
	toAccount, newErr := server.checkIsValidCurrencyOfAccount(ctx, req.ToAccountID, req.Currency)
	if !newErr {
		return
	}

	if fromAccount.Currency != toAccount.Currency {
		err := errors.New("Currency doesn't match")
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return
	}
	arg := db.CreateTransferParams{
		FromAccountID: req.FromAccountID,
		ToAccountID:   req.ToAccountID,
		Amount:        req.Amount,
	}
	result, newErr2 := server.store.TransferTx(ctx, arg)
	if newErr2 != nil {
		ctx.JSON(http.StatusInternalServerError, errRes(newErr2))
		return
	}
	ctx.JSON(http.StatusOK, result)
}
func (server *Server) checkIsValidCurrencyOfAccount(ctx *gin.Context, accountID int64, currency string) (db.Account, bool) {
	account, err := server.store.GetAccountForUpdate(ctx, accountID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errRes(err))
			return account, false
		}
		ctx.JSON(http.StatusInternalServerError, errRes(err))
		return account, false
	}
	if account.Currency != currency {
		err := fmt.Errorf("currency is mismatched")
		ctx.JSON(http.StatusBadRequest, errRes(err))
		return account, false
	}
	return account, true
}
