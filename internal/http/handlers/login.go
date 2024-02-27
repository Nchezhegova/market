package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/Nchezhegova/market/internal/config"
	"github.com/Nchezhegova/market/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Login(c *gin.Context) {
	var user models.UserModel
	var buf bytes.Buffer

	_, err := buf.ReadFrom(c.Request.Body)
	if err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if err := json.Unmarshal(buf.Bytes(), &user); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	if token, err := user.Login(c); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	} else {
		c.SetCookie(config.NAME_TOKEN, token, 3600, "/", config.ADDRSERV, false, true)
	}

	c.String(http.StatusOK, "Success login")
}
