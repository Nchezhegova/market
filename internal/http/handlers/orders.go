package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/Nchezhegova/market/internal/config"
	"github.com/Nchezhegova/market/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func LoadOrders(c *gin.Context) {
	var user models.UserModel
	var orders models.OrderModel
	var buf bytes.Buffer
	var uid int

	token, err := c.Cookie(config.NAME_TOKEN)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err, uid = user.CheckToken(c, token); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	_, err = buf.ReadFrom(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	number := buf.String()
	if err := orders.AddOrder(c, number, uid); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.String(http.StatusOK, "Success adding")
}

func GetOrders(c *gin.Context) {
	var user models.UserModel
	var uid int
	token, err := c.Cookie(config.NAME_TOKEN)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if err, uid = user.CheckToken(c, token); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	orders := models.GetOrders(c, uid)

	ordersByte, err := json.Marshal(orders)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "application/json", ordersByte)
}
