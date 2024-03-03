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
	if err = orders.CheckNumber(c, number); err != nil {
		c.AbortWithError(http.StatusUnprocessableEntity, err)
		return
	}
	qwe := orders.CheckOrder(c)
	if orderUser := qwe; orderUser != 0 {
		if orderUser == uid {
			c.AbortWithStatus(http.StatusOK)
			return
		}
		c.AbortWithStatus(http.StatusConflict)
		return
	}
	if err := orders.AddOrder(c, uid.(int)); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
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

	orders := models.GetOrders(c, uid.(int))
	log.Logger.Info("Response orders:", zap.Any("orders", orders))

	c.JSON(http.StatusOK, orders)
}
