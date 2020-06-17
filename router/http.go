package router

import (
	"api/conf"
	"api/controller"
	"api/log"
	"api/middleware"
	"api/xerror"
	"reflect"
	"runtime"
	"strings"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
)

func wrapper(f interface{}) func(*gin.Context) {
	fc := reflect.ValueOf(f)
	typ := fc.Type()
	if typ.Kind() != reflect.Func {
		log.Panicf("not function")
	}
	if typ.NumIn() != 2 {
		log.Panicf("number of params not equels to 2")
	}
	if typ.In(0).String() != "*gin.Context" {
		log.Panicf("first parameter should be of type *gin.Context")
	}
	if typ.NumOut() != 2 {
		log.Panicf("number of return values not equels to 2")
	}
	fullName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	strs := strings.Split(fullName, ".")
	viewName := strs[len(strs)-1]

	return func(c *gin.Context) {
		c.Set("view_name", viewName)

		_ = c.Request.ParseForm()
		for _, p := range c.Params {
			c.Request.Form[p.Key] = []string{p.Value}
		}
		req := reflect.New(typ.In(1).Elem()).Interface()
		c.Set("request", req)
		if err := c.ShouldBind(req); err != nil {
			log.Error(err)
			c.Abort()
			c.JSON(400, xerror.Unknown)
		}
		in := []reflect.Value{
			reflect.ValueOf(c),
			reflect.ValueOf(req),
		}
		res := fc.Call(in)
		if !res[1].IsNil() {
			err, ok := res[1].Interface().(*xerror.Error)
			if ok {
				c.Abort()
				c.JSON(err.HttpStatus, err)
			} else {
				c.Abort()
				c.JSON(400, xerror.Unknown)
			}
			return
		}
		c.JSON(200, res[0].Interface())
	}
}

func Register(app *gin.Engine) {
	app.StaticFile("/api/swagger.json", "dto/api.swagger.json")
	url := ginSwagger.URL("/api/swagger.json")
	app.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	store, err := redis.NewStore(10, "tcp", conf.Conf.Redis.Dsn, conf.Conf.Redis.Password, []byte("secret"))
	if err != nil {
		panic(err)
	}
	store.Options(sessions.Options{MaxAge: 84200, Path: "/"})

	app.Use(sessions.Sessions("api", store))
	//sso 单点登录
	sso := app.Group("/sso")
	{
		sso.GET("/login", wrapper(controller.Login))
		sso.GET("/welcome", controller.Welcome)
		sso.GET("/welcome2", middleware.RequireLogin, controller.Welcome2)
	}
}
