package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/oschwald/geoip2-golang"
	"log"
	"scoop-order/repository/transactions"
)

type Server struct {
	transaction   transactions.Transaction
	router        *gin.Engine
	dbGeo         *geoip2.Reader
	Logger        *log.Logger
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
}

func NewServer(router *gin.Engine, transaction transactions.Transaction, log *log.Logger, WarningLogger *log.Logger, InfoLogger *log.Logger, ErrorLogger *log.Logger) *Server {
	return &Server{
		router:      router,
		transaction: transaction,
		Logger:      log,
		WarningLogger: WarningLogger,
		InfoLogger: InfoLogger,
		ErrorLogger: ErrorLogger,
	}
}

func (server *Server) Start(address string) error {
	// run function that initializes the routes
	r := server.Routes()

	// run the server through the router
	err := r.Run(address)

	if err != nil {
		log.Fatal(fmt.Errorf("server - there was an error calling Run on router: %v", err), nil)
		return err
	}
	return nil
}

func errorResponse(err error) gin.H {
	return gin.H{
		"status":  false,
		"message": err.Error()}
}

func successResponse(msg string, data interface{}) gin.H {
	return gin.H{
		"status":  true,
		"message": msg,
		"data":    data,
	}
}
