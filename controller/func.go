package controller

import (
	"api/pkg/simulator/util"
	"github.com/bitly/go-simplejson"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

func saveSession(ctx *gin.Context, id int64) {
	session := sessions.Default(ctx)
	session.Set("uid", id)
	_ = session.Save()
}

func newConfig(r *http.Request) (*util.Config, error) {
	cfg, _ := util.NewConfig()
	if r.Method == http.MethodGet {
		cfg.Set("protocol", "http")
		cfg.Set("tcpport", "80")
		cfg.Set("channel_factory", "tcp-reveiver")
		cfg.Set("channelname", "tcp_http_reveiver")
		cfg.Set("sip", "192.168.10.135")

		return cfg, nil
	}
	var js *simplejson.Json
	var err error

	body, _ := ioutil.ReadAll(r.Body)
	if js, err = simplejson.NewJson(body); err != nil {
		return nil, err
	}

	if params, err := js.Map(); err != nil {
		return nil, err

	} else {
		for key, val := range params {
			cfg.Set(key, val)
		}
	}
	return cfg, nil
}
