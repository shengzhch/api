package controller

import (
	"api/controller/validator"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func init() {
	binding.Validator = new(validator.DefaultValidator)
}

func Welcome(c *gin.Context) {
	c.String(http.StatusOK, "Welcome", )
}

func Welcome2(c *gin.Context) {
	c.String(http.StatusOK, "Welcome2 %v", c.MustGet("uid"))
}
