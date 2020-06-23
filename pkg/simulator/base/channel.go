package base

import (
	"api/log"
	"api/pkg/simulator/util"
	"fmt"
	"github.com/gorilla/websocket"
)

var (
	channelFactories  = make(map[string]ChannelFactory)
	observerFactories = make(map[string]ChannelObserverFactory)
)

type ChannelFactory interface {
	New(*util.Config) (Channel, error)
}

type ChannelObserver interface {
	Name() string
	DataAvailable(Protocol, Channel, Packet, bool, error)
}

type ChannelObserverFactory interface {
	New(*util.Config, Channel, *websocket.Conn) ChannelObserver
}

type Channel interface {
	GetProtocol() Protocol
	SetProtocol(Protocol)

	RegisterObserver(ChannelObserver)
	UnRegisterObserver(string)

	DataAvailable(Protocol, Packet, bool, error)
	Start() error
	Stop() error

	Configuration() *util.Config
}

func RegisterChannelFactory(s string, creator ChannelFactory) {
	channelFactories[s] = creator
}

func RegisterObserverFactory(name string, factory ChannelObserverFactory) {
	observerFactories[name] = factory
}

func NewChannel(c *util.Config) (Channel, error) {
	if c.Get("protocol").MustString() == "" || c.Get("tcpport").MustString() == "" {
		return nil, fmt.Errorf("invalid parameters for ethernet\n")
	}

	cf := c.Get("channel_factory").MustString()
	if channelFactories[cf] == nil {
		log.Error("channel factory is null")
		return nil, fmt.Errorf("channel is error")
	}
	proto := c.Get("protocol").MustString()
	// create protocol
	p, err := NewProtocol(proto, c)
	if p == nil {
		return nil, err
	}
	ch, err := channelFactories[cf].New(c)
	if err != nil {
		return nil, err
	}

	ch.SetProtocol(p)
	return ch, nil
}

// create observers from config
func NewObserver(name string, c *util.Config, ch Channel, conn *websocket.Conn) ChannelObserver {
	if of := observerFactories[name]; of != nil {
		log.Println("has create a new observer")
		return of.New(c, ch, conn)
	}
	return nil
}
