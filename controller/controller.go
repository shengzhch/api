package controller

import (
	"api/conf"
	"api/controller/validator"
	"api/dto"
	"api/lib/forhttp"
	"api/log"
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func init() {
	binding.Validator = new(validator.DefaultValidator)
}

func Welcome(c *gin.Context) {
	c.String(http.StatusOK, "Welcome")
}

func Welcome2(c *gin.Context) {
	c.String(http.StatusOK, "Welcome2 %v", c.MustGet("uid"))
}

func WelcomeWithJae(c *gin.Context) {
	conn := forhttp.CreateGrpcConn("localhost:"+conf.Conf.Grpc.Port, c)
	grpcClient := dto.NewApiServiceClient(conn)
	grpcResp, err := grpcClient.Welcome(context.Background(), &dto.WelcomeReq{})

	log.Info("grpcResp ", grpcResp, "err ", err)
	c.String(http.StatusOK, "WelcomeJae")
}

func WelcomeWithJae2(c *gin.Context) {
	c.String(http.StatusOK, "WelcomeJae2")
}
