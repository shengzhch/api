package channel

import (
	"api/log"
	"api/pkg/simulator/base"
	"api/pkg/simulator/util"
	"bufio"
	"errors"
	"net"
	"os/exec"
	"strings"
)

// TCPChannel factory
type TCPChannelFactory struct{}

func (this *TCPChannelFactory) New(c *util.Config) (base.Channel, error) {
	ch := &TCPChannel{
		config:            c,
		name:              c.Get("channelname").MustString(),
		sip:               c.Get("sip").MustString(),
		dip:               c.Get("dip").MustString(),
		port:              c.Get("tcpport").MustString(),
		autostop:          c.Get("autostop").MustBool(),
		observers:         make(map[string]base.ChannelObserver),
		stopflag:          make(chan bool, 1),
		newConnectionChan: make(chan net.Conn, 1),
	}
	if base.IsProtocolSupport(c.Get("protocol").MustString()) == false {
		return nil, errors.New("unsuported protocol\n")
	}
	return ch, nil
}

func init() {
	base.RegisterChannelFactory("tcp-channle", &TCPChannelFactory{})
}

type TCPChannel struct {
	config            *util.Config
	name              string
	sip               string
	dip               string
	port              string
	isbindIp          bool
	isstoptelnet      bool
	autostop          bool
	proto             base.Protocol
	observers         map[string]base.ChannelObserver
	status            string
	stopflag          chan bool
	newConnectionChan chan net.Conn
	conn              net.Conn
	faceName          string
}

func (this *TCPChannel) Configuration() *util.Config {
	return this.config
}

func (this *TCPChannel) GetProtocol() base.Protocol {
	return this.proto
}

func (this *TCPChannel) SetProtocol(proto base.Protocol) {
	this.proto = proto
}

func (this *TCPChannel) RegisterObserver(obs base.ChannelObserver) {
	if this.observers[obs.Name()] == nil {
		this.observers[obs.Name()] = obs
	}
}

func (this *TCPChannel) UnRegisterObserver(name string) {
	if this.observers[name] != nil {
		delete(this.observers, name)
	}
}

// notified when bottom layer's data is available
func (this *TCPChannel) DataAvailable(proto base.Protocol, pkt base.Packet, isFirstItem bool, err error) {
	for _, observer := range this.observers {
		observer.DataAvailable(proto, this, pkt, isFirstItem, err)
	}
}

func (this *TCPChannel) Start() error {
	select {
	case <-this.stopflag:
	default:
	}
	errIp := this.checkIp()
	if errIp != nil {
		log.Info("Check Ip has error")
		return errIp
	}

	log.Info("simulation infomation", this.sip, this.dip, this.port)
	err := base.GetServiceProxy().AddService(base.ServiceInfo{Sip: this.sip, Dip: this.dip, Port: this.port, ConnCh: this.newConnectionChan})
	if err != nil {
		log.Error("add serviceinfo failded during start simulate ", err)
		return err
	}

	if !base.GetServiceProxy().HastcpServer(this.port) {
		err := base.GetServiceProxy().NewtcpServerWithStart(this.config.Get("tcpport").MustString())
		if err != nil {
			log.Info("tcp server listen and add tcp server into tcpserverList")
			return err
		}
	} else {
		if !base.GetServiceProxy().IsUp(this.port) {
			errLis := base.GetServiceProxy().StartListen(this.port)
			if errLis != nil {
				log.Error("start listen faild", errLis)
				return errLis
			}
		}
	}
	if err = this.Run(); err != nil {
		log.Error("Run fail!!", err)
		return err
	}
	return nil
}

func (this *TCPChannel) Run() error {
	defer util.Run()
	protocol := this.GetProtocol()
	for {
		select {
		case conn := <-this.newConnectionChan:
			go func(conn net.Conn) {
				defer func() {
					_ = conn.Close()
				}()
				if this.conn != nil {
					//服务端主动关闭该连接,防止多连接接入
					_ = conn.Close()
					log.Info("multipyle connection exist for tcp channel:", this.port)
					return
				}

				//需要选项协商的过程
				//this.ProcessOP(conn)
				//log.Info("ProcessOP have Finished")

				this.conn = conn

				br := bufio.NewReader(this.conn)
				log.Info("bufio.NewReade have done")
				for {
					data, err := br.ReadBytes(protocol.GetDelim())
					log.Info("br.ReadBytes hava done")
					if err != nil {
						log.Info("remoteaddr disconnect")
						if this.autostop {
							_ = this.Stop()
						}
						this.conn = nil
						break
					}
					if len(data) == 0 {
						if this.autostop {
							_ = this.Stop()
						}
						log.Info("channel simulate stop")
						break
					}
					log.Info("date from conn ", string(data))
					if string(data) == "connection exit\r\n" {
						log.Info("remoteaddr ask to quit")
						_ = conn.Close()
						this.conn = nil
						break
					}
					if data[0] == byte(255) {
						log.Info("this is OP")
					} else {
						pkt, err := protocol.Decode(data, 1, false)
						if err != nil {
							log.Error("Decode err.", err)
						} else if pkt == nil {
							// no data
						} else {
							this.DataAvailable(protocol, pkt, false, err)
						}
					}
				}
			}(conn)
		case <-this.stopflag:
			if this.conn != nil {
				_ = this.conn.Close()
				this.conn = nil
			}
			log.Info("Run stop")
			return nil
		}
	}
}

func (this *TCPChannel) Stop() error {
	select {
	case <-this.stopflag:
	default:
	}
	this.stopflag <- true
	return nil
}

func (this *TCPChannel) GetFaceName() string {
	intf, _ := net.Interfaces()
	for _, v := range intf {
		ips, err := v.Addrs()
		if err != nil {
			continue
		}
		for _, ip := range ips {
			if strings.Contains(ip.String(), this.sip) {
				return v.Name
			}
		}
	}
	return intf[0].Name
}

func (this *TCPChannel) ProcessOP(conn net.Conn) {
	this.proto.ProcessOP(conn)
}

func (this *TCPChannel) checkIp() error {
	addrs, _ := net.InterfaceAddrs()
	for _, address := range addrs {
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				if ipnet.IP.String() == this.dip {
					log.Info("IP exist")
					return nil
				}
			}
		}
	}

	a := this.GetFaceName()
	com := "netsh interface ip add address name=\"" + a + "\" " + this.dip + " 255.255.255.0"
	c := exec.Command("cmd", "/C", com)
	if err := c.Run(); err != nil {
		log.Error("Bind IP fail :", err)
		return err
	} else {
		this.faceName = a
		this.isbindIp = true
		return nil
	}
}
