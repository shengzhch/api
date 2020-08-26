package router

import (
	"api/dto"
	"api/handler"
	"google.golang.org/grpc"
)

func RegisterGrpc(svc *grpc.Server) {
	dto.RegisterApiServiceServer(svc, handler.NewWelcomeServiceHandler())
	dto.RegisterApiService2Server(svc, handler.NewWelcomeServiceHandler())
}
