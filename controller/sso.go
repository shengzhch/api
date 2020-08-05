package controller

import (
	"api/dto"
	"api/log"
	"api/xerror"
	"errors"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, req *dto.LoginReq) (*dto.LoginResp, error) {
	if req.LoginType == dto.LoginType_NamePd {
		if req.Name == "root" && req.Passwd == "demo" {
			saveSession(c, 1)
			return &dto.LoginResp{}, nil
		}
	}

	err := MakeErr()

	log.Info("info ", err)
	if err != nil {
		c.Set("STACK", true)
		log.Error("err", err)
	}
	return nil, xerror.NoLogin
}

func MakeErr() error {
	return errors.New("this is a error")
}
