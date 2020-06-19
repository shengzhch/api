package base

import (
	"fmt"
	"net"

	"api/pkg/simulator/util"
)

var (
	protocolFactories = make(map[string]ProtocolFactory)
)

type Protocol interface {
	Name() string
	DealData([]byte) (Packet, error)
	Decode([]byte, int, bool) (Packet, error)
	Encode([]byte) ([]byte, error)
	ProcessOP(conn net.Conn)
	GetDelim() byte
}

// protocol factory
type ProtocolFactory interface {
	New(*util.Config) (Protocol, error)
}

func NewProtocol(name string, c *util.Config) (Protocol, error) {
	if protocolFactories[name] != nil {
		return protocolFactories[name].New(c)
	}
	return nil, fmt.Errorf("unsported protocol")
}

func RegisterProtocolFactory(name string, factory ProtocolFactory) {
	protocolFactories[name] = factory
}

func IsProtocolSupport(name string) bool {
	return protocolFactories[name] != nil
}
