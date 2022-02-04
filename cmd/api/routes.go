package api

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func (server *Server) Routes() *gin.Engine {
	router := server.router
	corsConfigs := cors.DefaultConfig()
	corsConfigs.AllowAllOrigins = true
	corsConfigs.AllowHeaders = append(corsConfigs.AllowHeaders, "User-Agent")
	corsConfigs.AllowHeaders = append(corsConfigs.AllowHeaders, "x-api-key")
	corsConfigs.AllowHeaders = append(corsConfigs.AllowHeaders, "Authorization")
	corsConfigs.AllowHeaders = append(corsConfigs.AllowHeaders, "Signature") // ini di body
	corsConfigs.AllowMethods = append(corsConfigs.AllowMethods, "OPTIONS")
	router.Use(cors.New(corsConfigs))
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

		pricing := v1.Group("/finalPrice")
		{
			//pricing.GET("/v2", server.GetValidPrice)
			pricing.GET("", server.GetValidPrice2)
		}

		checkout := v1.Group("/checkout")
		{
			checkout.Use(server.authenticate()).POST("", server.Checkout)
		}

		complete := v1.Group("/complete-order")
		{
			complete.POST("", server.CompleteOrder)
			complete.POST("/v1", server.Complete)
		}

		paymentGateways := v1.Group("/payment-gateways")
		{
			paymentGateways.GET("", server.GetPaymentGateway())
		}
	}

	return router
}