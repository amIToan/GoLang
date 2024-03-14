package api

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	db "sgithub.com/techschool/simplebank/db/sqlc"
	"sgithub.com/techschool/simplebank/token"
	"sgithub.com/techschool/simplebank/util"
)

// server serves HTTP requests for our banking service
type Server struct {
	config     util.Config
	store      db.Store
	router     *gin.Engine
	tokenMaker token.Maker
}

// NewServer creates a new HTTP
func NewServer(config util.Config, store db.Store) (*Server, error) {
	tokenMaker, err := token.NewPasetoMaker(config.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("can not initiate token maker : %w", err)
	}
	r := gin.Default()
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("currency", validCurrency)
	}
	server := &Server{
		config:     config,
		store:      store,
		tokenMaker: tokenMaker,
	}
	//auth Routes
	authRoutes := r.Group("/").Use(authMiddleWare(server.tokenMaker))
	//accounts
	authRoutes.POST("/account", server.createAccount)
	authRoutes.GET("/account/:id", server.getAccountById)
	authRoutes.GET("/accounts", server.getAccountsList)
	authRoutes.PUT("/account/update", server.updateAccountById)
	authRoutes.DELETE("/account/delete/:id", server.DeleteAccountById)
	/////////////////////
	//transfer
	authRoutes.POST("/transfer/createTransfer", server.CreateTransfer)
	////////////////////
	//user
	r.POST("/user", server.createUser)
	authRoutes.GET("/user/:username", server.GetUser)
	r.POST("user/login", server.Login)
	///////////////////
	//token
	r.POST("/token/renew", server.renewAccessToken)
	//////////////////////////////////////////////
	server.router = r
	return server, nil
}
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}
func errRes(err error) gin.H {
	return gin.H{"error": err.Error()}
}
