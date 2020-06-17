package controller

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func saveSession(ctx *gin.Context, id int64) {
	session := sessions.Default(ctx)
	session.Set("uid", id)
	_ = session.Save()
}
