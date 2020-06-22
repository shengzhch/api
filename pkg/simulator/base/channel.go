package base

import (
	"api/pkg/simulator/util"
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
