package controller

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func Welcome(c *gin.Context) {
	c.String(http.StatusOK, "Welcome")
}
