package forgrpc

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func CreateGrpcClient(serviceAddress string, ctx context.Context) *grpc.ClientConn {
	conn, err := grpc.Dial(serviceAddress, grpc.WithInsecure(), grpc.WithUnaryInterceptor(ClientInterceptor(Tracer, ctx)))
	if err != nil {
		fmt.Println(serviceAddress, "grpc conn err:", err)
	}
	return conn
}
