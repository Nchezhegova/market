package handlers

import (
	"bytes"
	"encoding/json"
	"github.com/Nchezhegova/market/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

func Registration(c *gin.Context) {
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

	if err := user.Add(c); err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.String(http.StatusOK, "Success adding")

}
