package handlers

import (
	"bytes"
	"github.com/Nchezhegova/market/internal/models"
	"github.com/gin-gonic/gin"
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
	orderUser, err := orders.CheckOrder(c)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
	}
	if orderUser != 0 {
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

	orders, err := models.GetOrders(c, uid.(int))
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.JSON(http.StatusOK, orders)
}
