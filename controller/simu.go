package controller

import (
	"api/log"
	"api/pkg/simulator/base"
	_ "api/pkg/simulator/channel"
	_ "api/pkg/simulator/observer"
	_ "api/pkg/simulator/protocol"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

var (
	upgrader = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
)

func MoniStart(c *gin.Context) {
	conf, err := newConfig(c.Request)

	if err != nil {
		log.Error("生成配置失败", err)
		return
	}

	ch, err := base.NewChannel(conf)
	if err != nil {
		log.Error("创建通道失败", err)
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("websocket upgrade failed")
		return
	}

	obv := base.NewObserver("monitor", ch.Configuration(), ch, conn)
	ch.RegisterObserver(obv)
	err = ch.Start()
	if err != nil {
		log.Error(err)
		return
	}
	return
}
