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
	uid, exists := c.Get("userID")
	if !exists {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	b := models.BalanceModel{}
	if b.GetBalance(c, uid.(int)) != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
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
	token, err := c.Cookie(config.NAMETOKEN)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	if uid, err = user.CheckToken(c, token); err != nil {
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
	withdrawalSum, _ := w.Sum.Float64()
	if b.Sum < withdrawalSum {
		c.AbortWithStatus(http.StatusPaymentRequired)
		return
	}

	if err := w.AddWithdrawal(c, uid); err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	c.String(http.StatusOK, "Success adding")
}

func Withdrawals(c *gin.Context) {
	uid, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User ID not found in context"})
		return
	}
	var err error
	var withdrawals []models.WithdrawalResp
	withdrawals, err = models.GetWithdrawal(c, uid.(int))
	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	if withdrawals == nil {
		c.AbortWithStatus(http.StatusNoContent)
		return
	}
	c.JSON(http.StatusOK, withdrawals)
}
