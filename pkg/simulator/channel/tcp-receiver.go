package channel

import (
	"errors"

	"api/log"
	"time"

	"api/pkg/simulator/base"
	"api/pkg/simulator/util"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"github.com/google/gopacket/tcpassembly"
)

var (
	dev, _ = pcap.FindAllDevs()
)

type TCPReceiver struct {
	config    *util.Config
	name      string
	ip        string
	port      int
	proto     base.Protocol
	observers map[string]base.ChannelObserver
	status    string
	statusDB  string
	handler   *pcap.Handle
	stopflag  chan bool
}

// TCPReceiver factory
type TCPReceiverFactory struct{}

func (this *TCPReceiverFactory) New(c *util.Config) (base.Channel, error) {
	var ch = TCPReceiver{
		config:    c,
		name:      c.Get("channelname").MustString(),
		observers: make(map[string]base.ChannelObserver),
		handler:   nil,
		stopflag:  make(chan bool, 1),
	}

	if base.IsProtocolSupport(c.Get("protocol").MustString()) == false {
		return nil, errors.New("unsuported protocol\n")
	}
	return &ch, nil
}

func (this *TCPReceiver) GetProtocol() base.Protocol {
	return this.proto
}

func (this *TCPReceiver) SetProtocol(proto base.Protocol) {
	this.proto = proto
}

func (this *TCPReceiver) RegisterObserver(obs base.ChannelObserver) {
	if this.observers[obs.Name()] == nil {
		this.observers[obs.Name()] = obs
	}
}

func (this *TCPReceiver) UnRegisterObserver(name string) {
	if this.observers[name] != nil {
		delete(this.observers, name)
	}
}

func (this *TCPReceiver) Start() error {
	select {
	case <-this.stopflag:
	default:
	}
	log.Info("channel start")
	var iface = this.config.Get("device").MustString() //no use
	log.Info("db_iface:", iface)
	var flag bool
	for _, devVal := range dev {
		for _, ipVal := range devVal.Addresses {
			if ipVal.IP.String() == this.config.Get("sip").MustString() || ipVal.IP.String() == this.config.Get("dip").MustString() {
				flag = true
				iface = devVal.Name
				break
			}
		}
		if flag {
			break
		}
	}

	log.Info("use_iface:", iface)
	//var iface = "eth0"
	// Set up pcap packet capture
	var err error
	if iface != "" {
		log.Info("Starting capture on interface %q", iface)
		this.handler, err = pcap.OpenLive(iface, 1600, true, pcap.BlockForever)
		if err != nil {
			log.Error(err)
			return err
		}
		log.Info("handler link type ", this.handler.LinkType())
	}

	defer this.handler.Close()
	// set bp filer
	tcpport := this.config.Get("tcpport").MustString()
	sip := this.config.Get("sip").MustString()
	dip := this.config.Get("dip").MustString()

	var filter = "host " + sip + " and host " + dip + " and port " + tcpport
	log.Info(filter)

	if err := this.handler.SetBPFFilter(filter); err != nil {
		log.Error(err)
		return err
	}

	if err := this.Run(); err != nil {
		log.Info("Run fail!!")
		return err
	}
	log.Info("start finished")
	return nil
}

func (this *TCPReceiver) Stop() error {
	// should we stop the pcap
	select {
	case <-this.stopflag:
	default:
	}
	this.stopflag <- true
	return nil
}

func (this *TCPReceiver) Run() error {
	defer util.Run()

	streamFactory := &TcpStreamFactory{channel: this}
	streamPool := tcpassembly.NewStreamPool(streamFactory)
	assembler := tcpassembly.NewAssembler(streamPool)

	log.Info("reading  packets")

	packetSource := gopacket.NewPacketSource(this.handler, this.handler.LinkType())
	packets := packetSource.Packets()
	ticker := time.Tick(time.Minute)
	for {
		select {
		case packet := <-packets:
			if packet == nil {
				log.Info("packet nil")
				return nil
			}
			if packet.NetworkLayer() == nil || packet.TransportLayer() == nil || packet.TransportLayer().LayerType() != layers.LayerTypeTCP {
				continue
			}
			tcp, ok := packet.TransportLayer().(*layers.TCP)
			if ok {
				assembler.AssembleWithTimestamp(packet.NetworkLayer().NetworkFlow(), tcp, packet.Metadata().Timestamp)
			} else {
				log.Info("packet.TransportLayer() is not tcp", packet.TransportLayer().LayerType())
			}
		case <-ticker:
			assembler.FlushOlderThan(time.Now().Add(time.Minute * -2))
		case <-this.stopflag:
			log.Info("receiver stop signal")
			return nil
		}
	}
}

// notified when bottowm layer's data is available
func (this *TCPReceiver) DataAvailable(proto base.Protocol, pkt base.Packet, isFirstItem bool, err error) {
	for _, observer := range this.observers {
		observer.DataAvailable(proto, this, pkt, isFirstItem, err)
	}
}

// helpers
func (this *TCPReceiver) Name() string                { return this.name }
func (this *TCPReceiver) Configuration() *util.Config { return this.config }
