package handler

import (
	"api/dto"
	"api/log"
	"context"
)

type WelcomeServiceHandler struct{}

func NewWelcomeServiceHandler() *WelcomeServiceHandler {
	return &WelcomeServiceHandler{}
}

func (w WelcomeServiceHandler) Welcome(context.Context, *dto.WelcomeReq) (*dto.WelcomeResp, error) {
	log.Info("Down")
	return &dto.WelcomeResp{}, nil
}


func (w WelcomeServiceHandler) Welcome2(context.Context, *dto.WelcomeReq) (*dto.WelcomeResp, error) {
	return &dto.WelcomeResp{}, nil
}
