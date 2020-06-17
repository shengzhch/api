package main

import (
	"api/conf"
	"api/log"
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
	// service init
	svc := service.NewHttpServer(conf.Conf)
	ctx, cancel := context.WithCancel(context.Background())
	errCh := make(chan error, 10)
	service.HttpRun(ctx, svc, errCh)

	go func() {
		for {
			select {
			case e, ok := <-errCh:
				if e != nil {
					log.Error(e)
				}
				if !ok {
					log.Info("error channel close")
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
			close(errCh)
			time.Sleep(100 * time.Millisecond)
			return
		default:
			//todo
			return
		}
	}
}
