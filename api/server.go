package api

import (
	"github.com/gin-gonic/gin"
	db "parus.i234.me/paultabaco/simplebank/db/sqlc"
)

// / Server serves HTTP requests for banking service.
type Server struct {
	store  db.Store
	router *gin.Engine
}

// NewSever creates new HTTP server and setup routing.
func NewServer(store db.Store) Server {
	server := &Server{store: store}
	router := gin.Default()

	router.POST("/accounts", server.createAccount)
	router.GET("/accounts/:id", server.getAccount)
	router.GET("/accounts", server.listAccount)

	router.POST("/transfers", server.createTransfer)

	server.router = router
	return *server
}

// / Start runs http server with address
func (server *Server) Start(address string) error {
	return server.router.Run(address)
}

func errorResponse(err error) gin.H {
	return gin.H{
		"error": err.Error(),
	}
}
