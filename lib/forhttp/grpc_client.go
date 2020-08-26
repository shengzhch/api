package forhttp

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	grpc_middeware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/opentracing/opentracing-go"
	"google.golang.org/grpc"
	"time"
)

func CreateGrpcConn(serviceAddress string, c *gin.Context) *grpc.ClientConn {
	var conn *grpc.ClientConn
	var err error

	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond*500)
	defer cancel()
	tracer, _ := c.Get("Tracer")
	parentSpanContext, _ := c.Get("ParentSpanContext")

	conn, err = grpc.DialContext(
		ctx,
		serviceAddress,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithUnaryInterceptor(
			grpc_middeware.ChainUnaryClient(
				ClientInterceptor(tracer.(opentracing.Tracer), parentSpanContext.(opentracing.SpanContext)),
			),
		),
	)

	if err != nil {
		fmt.Println(serviceAddress, "grpc conn err:", err)
	}
	return conn
}
