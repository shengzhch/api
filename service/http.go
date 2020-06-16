package service

import (
	"context"
	"fmt"
	"net/http"
	"api/conf"
	"api/log"
	"api/router"

	"github.com/gin-gonic/gin"
)

type HttpServer struct {
	Port int
	App  *gin.Engine
}

func NewServer(c *conf.Config) *HttpServer {
	app := newGinEngine(c.Http.Mode)
	return &HttpServer{Port: c.Http.Port, App: app}
}

func newGinEngine(mode string) *gin.Engine {
	gin.SetMode(mode)
	app := gin.Default()
	return app
}

func Run(ctx context.Context, s *HttpServer, errCh chan error) {
	router.Register(s.App)
	srv := &http.Server{
		Addr:    fmt.Sprintf(":%v", s.Port),
		Handler: s.App,
	}
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()
	go func() {
		log.Infof("http server running on port: %v", s.Port)
		<-ctx.Done()
		err := srv.Shutdown(context.Background())
		if err != nil {
			errCh <- err
		} else {
			log.Info("http server has shut down")
		}
	}()
}
