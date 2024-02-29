package server

import (
	"github.com/Nchezhegova/market/internal/http/handlers"
	"github.com/gin-gonic/gin"
)

func StartServer(addr string) {
	r := gin.Default()
	r.ContextWithFallback = true

	r.POST("/api/user/register", func(c *gin.Context) {
		handlers.Registration(c)
	})
	r.POST("/api/user/login", func(c *gin.Context) {
		handlers.Login(c)
	})
	r.POST("/api/user/orders", func(c *gin.Context) {
		handlers.LoadOrders(c)
	})
	r.GET("/api/user/orders", func(c *gin.Context) {
		handlers.GetOrders(c)
	})
	r.GET("/api/user/balance", func(c *gin.Context) {
		handlers.GetBalance(c)
	})
	r.POST("/api/user/balance/withdraw", func(c *gin.Context) {
		handlers.AddWithdrawal(c)
	})
	r.GET("/api/user/withdrawals", func(c *gin.Context) {
		handlers.Withdrawals(c)
	})

	r.Run(addr)
}
