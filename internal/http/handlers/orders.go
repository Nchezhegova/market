package handlers

import (
	"bytes"
	"github.com/Nchezhegova/market/internal/log"
	"github.com/Nchezhegova/market/internal/models"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
)

func LoadOrders(c *gin.Context) {
	var orders models.OrderModel
	var buf bytes.Buffer

	uid, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	number := buf.String()
	if err := orders.AddOrder(c, number, uid.(int)); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.String(http.StatusAccepted, "Success adding")
}

func GetOrders(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	var orders []models.OrderWithdrawal
	orders = models.GetOrders(c, uid.(int))
	log.Logger.Info("Response orders:", zap.Any("orders", orders))

	//ordersByte, err := json.Marshal(orders)
	//if err != nil {
	//	c.AbortWithError(http.StatusInternalServerError, err)
	//	return
	//}
	//c.Data(http.StatusOK, "application/json", ordersByte)
	c.JSON(http.StatusOK, orders)
}
