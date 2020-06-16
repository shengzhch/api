package middleware

import (
	"api/xerror"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func RequireLogin(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("uid");
	if v == nil {
		c.Abort()
		c.JSON(403, xerror.NoLogin)
		return
	}
	c.Next()
}
