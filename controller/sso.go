package controller

import (
	"api/dto"
	"api/xerror"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context, req *dto.LoginReq) (*dto.LoginResp, error) {
	if req.LoginType == dto.LoginType_NamePd {
		if req.Name == "root" && req.Passwd == "demo" {
			saveSession(c, 1)
			return &dto.LoginResp{}, nil
		}
	}
	return nil, xerror.NoLogin
}
