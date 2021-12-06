package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (server *Server) Routes() *gin.Engine {
	router := server.router
	// Setup middleware
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.Default())

	// group all routes under /v1
	v1 := router.Group("/v1")
	{
		v1.GET("/status", server.ApiStatus())
		// prefix the order routes
		orders := v1.Group("/orders")
		{
			orders.GET("", server.GetOrderList())
			orders.GET("/id/:id", server.GetOrderByID)
			orders.GET("/order_number/:id", server.GetOrderByOrderNumber)
			orders.Use(server.authenticate()).GET("/user", server.GetOrderByUser)
		}

		pricing := v1.Group("/price")
		{
			pricing.GET("", server.GetValidPrice)
		}

		checkout := v1.Group("/checkout")
		{
			checkout.Use(server.authenticate()).POST("", server.Checkout)
		}

		complete := v1.Group("/complete-order")
		{
			complete.POST("", server.Complete)
		}

		paymentGateways := v1.Group("/payment-gateways")
		{
			paymentGateways.GET("", server.GetPaymentGateway())
		}
	}

	return router
}