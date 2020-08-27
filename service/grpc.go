package service

import (
	"api/conf"
	"api/lib/forgrpc"
	"api/log"
	"api/middleware"
	"api/router"
	"context"
	"fmt"
	"net"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"google.golang.org/grpc"
)

type GrpcServer struct {
	Port string
	App  *grpc.Server
}

func NewGrpcServer(conf *conf.Config) *GrpcServer {
	app := newGRPC(conf)
	return &GrpcServer{Port: conf.Grpc.Port, App: app}
}

func p(pp interface{}) error {
	return fmt.Errorf("%v", pp)
}

func newGRPC(c *conf.Config) *grpc.Server {
	_, _, err := forgrpc.NewJaegerTracer(c.Grpc.ServiceName, c.Jaeger.Port)
	if err != nil {
		log.Panic(err)
	}
	_, usi := forgrpc.ServerOption()
	s := grpc.NewServer(
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				usi,
				middleware.LogUnary,
				grpc_recovery.UnaryServerInterceptor(grpc_recovery.WithRecoveryHandler(p)),
			),
		),
	)
	return s
}

func GrpcRun(ctx context.Context, s *GrpcServer, errCh chan error) {
	//var listenCh = make(chan net.Listener, 1)
	go func() {
		lis, err := net.Listen("tcp", "0.0.0.0:"+s.Port)
		if err != nil {
			log.Fatalf("Failed to listen: %v", err)
		}
		//listenCh <- lis
		log.Infof("Starting grpc on 0.0.0.0:%v", s.Port)
		router.RegisterGrpc(s.App)
		if err := s.App.Serve(lis); err != nil {
			log.Fatalf("Failed to serve: %v", err)
		}
	}()

	go func() {
		<-ctx.Done()
		s.App.GracefulStop()
		log.Info("grpc server has shut down")
		//listener has been close

		//lis, ok := <-listenCh
		//if ok {
		//	err := lis.Close()
		//	if err != nil {
		//		log.Error(err)
		//		select {
		//		case errCh <- err:
		//		default:
		//			log.Error("send err to channel failed, err", err)
		//		}
		//	} else {
		//		log.Info("grpc list stop")
		//	}
		//}
	}()
}
