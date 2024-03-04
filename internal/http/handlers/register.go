package handlers

import (
	"encoding/json"
	"github.com/Nchezhegova/market/internal/config"
	"github.com/Nchezhegova/market/internal/models"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

func Registration(c *gin.Context, addr string) {
	var user models.UserModel
	b, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(b, &user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if err := user.Add(c); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	token, err := user.Login(c)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.SetCookie(config.NAMETOKEN, token, 3600, "/", addr, false, true)
	c.String(http.StatusOK, "Success adding")
}
