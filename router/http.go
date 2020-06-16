package router

import (
	"api/controller"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/swaggo/gin-swagger/swaggerFiles"

	"github.com/gin-gonic/gin"
)

func Register(app *gin.Engine) {
	app.StaticFile("/api/swagger.json", "dto/api.swagger.json")
	url := ginSwagger.URL("/api/swagger.json")
	app.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	v1 := app.Group("/api/v1")
	{
		v1.GET("/welcome", controller.Welcome)
	}
}
