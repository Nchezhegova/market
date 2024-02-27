package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/Nchezhegova/market/internal/config"
	"github.com/Nchezhegova/market/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func GetBalance(c *gin.Context) {
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

	b := models.BalanceModel{}
	if b.GetBalance(c, uid) != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	balanceByte, err := json.Marshal(b)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "application/json", balanceByte)
}

func AddWithdrawal(c *gin.Context) {
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

	var buf bytes.Buffer
	_, err = buf.ReadFrom(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	var w models.WithdrawalModel
	if err := json.Unmarshal(buf.Bytes(), &w); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	b := models.BalanceModel{}
	if b.GetBalance(c, uid) != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if b.Sum < w.Sum {
		c.AbortWithError(http.StatusPaymentRequired, err)
		return
	}

	if err := w.AddWithdrawal(c, uid); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusOK, "Success adding")
}

func Withdrawals(c *gin.Context) {
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

	var withdrawals []models.WithdrawalModel
	err, withdrawals = models.GetWithdrawal(c, uid)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if withdrawals == nil {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	wByte, err := json.Marshal(withdrawals)
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.Data(http.StatusOK, "application/json", wByte)
}
