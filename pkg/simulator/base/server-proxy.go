package base

import (
	"errors"
	"net"
	"strings"
	"sync"

	"api/log"
)

type ServiceInfo struct {
	Sip    string
	Dip    string
	Port   string
	ConnCh chan net.Conn
}

type tcpServer struct {
	Listener net.Listener //socket
	Port     string
	IsUp     bool
	ch       chan int
}

type ServiceInfoList []ServiceInfo

type ServiceProxy struct {
	serviceInfoListLock sync.Mutex
	serviceInfoList     map[string]ServiceInfoList
	tcpServerListLock   sync.Mutex
	tcpServerList       map[string]*tcpServer
}

var (
	serviceProxy *ServiceProxy
	once         sync.Once
)

func GetServiceProxy() *ServiceProxy {
	once.Do(func() {
		serviceProxy = &ServiceProxy{serviceInfoList: make(map[string]ServiceInfoList), tcpServerList: make(map[string]*tcpServer)}
	})
	return serviceProxy
}

//todo
func (this *ServiceProxy) Destory() {

}

func (this *ServiceProxy) HastcpServer(port string) bool {
	return serviceProxy.tcpServerList[port] != nil

}

func (this *ServiceProxy) AddService(info ServiceInfo) error {
	// check wether the specified channel exist
	this.serviceInfoListLock.Lock()
	if this.serviceInfoList[info.Port] != nil {
		for _, ser := range this.serviceInfoList[info.Port] {
			if ser.Sip == info.Sip && ser.Dip == info.Dip && ser.Port == info.Port {
				this.serviceInfoListLock.Unlock()
				return errors.New("same channel exist")
			}
		}
		if info.ConnCh == nil {
			this.serviceInfoListLock.Unlock()
			return errors.New("no notification created for the channel")
		}
	}
	// add the services
	log.Info("add serviceInfo into serviceInfoList")
	this.serviceInfoList[info.Port] = append(this.serviceInfoList[info.Port], info)
	this.serviceInfoListLock.Unlock()
	return nil
}

func (this *ServiceProxy) UpdateService(bsip string, bdip string, bport string, nsip string, ndip string, nport string) {
	for k, v := range this.serviceInfoList {
		if k == bport {
			for in, ser := range v {
				if ser.Sip == bsip && ser.Dip == bdip {
					this.serviceInfoList[k][in].Sip = nsip
					this.serviceInfoList[k][in].Dip = ndip
					this.serviceInfoList[k][in].Dip = nport
					log.Info("services info has update")
					return
				}
			}
		}
	}
}

func (this *ServiceProxy) NewtcpServerWithStart(port string) error {
	addr := ":" + port
	netListen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("net.Listen failed", err)
		return nil
	}
	tcpserver := &tcpServer{Port: port, IsUp: true, Listener: netListen}
	this.addtcpServer(tcpserver)
	go Listen(this, port, netListen)
	return nil
}

func (this *ServiceProxy) addtcpServer(t *tcpServer) {
	this.tcpServerListLock.Lock()
	this.tcpServerList[t.Port] = t
	this.tcpServerListLock.Unlock()
}
func (this *ServiceProxy) IsUp(port string) bool {
	return this.tcpServerList[port].IsUp
}

//关闭监听，并删除tcpserver。
func (this *ServiceProxy) DeleteTcpServer(port string) error {
	this.tcpServerListLock.Lock()
	if this.tcpServerList[port] != nil {
		_ = this.tcpServerList[port].Listener.Close()
		delete(this.tcpServerList, port)
		this.tcpServerListLock.Unlock()
		return nil
	} else {
		log.Error("no tcpserver to this port")
	}
	this.tcpServerListLock.Unlock()
	return nil
}

//只关闭tcpserver的监听
func (this *ServiceProxy) StopListen(port string) error {
	this.tcpServerListLock.Lock()
	if this.tcpServerList[port] != nil {
		_ = this.tcpServerList[port].Listener.Close()
		this.tcpServerList[port].IsUp = false
		this.tcpServerListLock.Unlock()
		return nil
	} else {
		log.Error("no tcpserver listen to this port")
	}
	this.tcpServerListLock.Unlock()
	return nil
}

//重新启监听
func (this *ServiceProxy) StartListen(port string) error {
	addr := ":" + port
	netListen, err := net.Listen("tcp", addr)
	if err != nil {
		log.Error("net.Listen failed", err)
		return errors.New("net Listen failed")
	}
	if this.tcpServerList[port].IsUp {
		log.Info("netlistener still up")
		return nil
	}
	this.tcpServerList[port].Listener = netListen
	this.tcpServerList[port].IsUp = true
	go Listen(this, port, netListen)
	return nil
}

func Listen(this *ServiceProxy, port string, netListen net.Listener) {
	var err error
	log.Infof("listen [port: %v]start：waiting for clients", port)
	defer func() {
		err = netListen.Close()
		if err != nil {
			log.Error("Listener close failed ", err)
		}
	}()
	for {
	L:
		conn, err := netListen.Accept()
		if err != nil {
			//表示监听已经关闭，应退出for循环
			log.Error("tcp server accept failed for channel: ", port)
			this.tcpServerList[port].IsUp = false
			return
		}
		log.Info("new connection has establish")
		localAddr := conn.LocalAddr().String()
		remoteAddr := conn.RemoteAddr().String()
		log.Info("conn LocalAddr", localAddr)
		log.Info("conn RemoteAddr", remoteAddr)
		for _, ser := range this.serviceInfoList[port] {
			if strings.Contains(localAddr, ser.Dip) && strings.Contains(remoteAddr, ser.Sip) {
				log.Info("write conn to serverinfo'conn")
				ser.ConnCh <- conn
				goto L
			}
		}
		log.Info("no service for this.conn,closing ...")
		err = conn.Close()
		if err != nil {
			log.Error("connect close failed ", err)
		}
	}
}
