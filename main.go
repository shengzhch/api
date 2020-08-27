package main

import (
	"api/conf"
	"api/log"
	"api/middleware"
	"api/service"

	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	flag.Parse()
	if err := conf.Init(); err != nil {
		log.Error("conf.Init() error(%v)", err)
		panic(err)
	}
	middleware.InitTracer(conf.Conf.Jaeger)

	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 10)
	// service init
	svc := service.NewHttpServer(conf.Conf)
	service.HttpRun(ctx, svc, errCh)

	grpc := service.NewGrpcServer(conf.Conf)
	service.GrpcRun(ctx, grpc, errCh)

	finishCh := make(chan struct{})
	go func() {
		for {
			select {
			case e, ok := <-errCh:
				if e != nil {
					log.Error(e)
				}
				if !ok {
					log.Info("error channel close")
					finishCh <- struct{}{}
					return
				}
			}
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		s := <-c
		log.Infof("Received system signal[%v]", s)
		switch s {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGSTOP, syscall.SIGINT, syscall.SIGHUP:
			cancel()
			time.Sleep(100 * time.Millisecond)
			close(errCh)
			<-finishCh
			return
		}
	}
}
