package channel

import (
	"api/log"
	"api/pkg/simulator/base"
	"bufio"
	"github.com/google/gopacket"
	"github.com/google/gopacket/tcpassembly"
	"github.com/google/gopacket/tcpassembly/tcpreader"
	"io"
)

// tcpStreamFactory implements tcpassembly.StreamFactory
type TcpStreamFactory struct {
	channel base.Channel
}

// tcpStream will handle the actual decoding of tcp requests.
type TcpStream struct {
	channel   base.Channel
	net       gopacket.Flow
	transport gopacket.Flow
	r         tcpreader.ReaderStream
	newdata   bool
}

func (h *TcpStreamFactory) New(net, transport gopacket.Flow) tcpassembly.Stream {
	hstream := &TcpStream{
		channel:   h.channel,
		net:       net,
		transport: transport,
		r:         tcpreader.NewReaderStream(),
		newdata:   true,
	}
	go hstream.run() // Important... we must guarantee that data from the reader stream is read.

	// ReaderStream implements tcpassembly.Stream, so we can return a pointer to it.
	return &hstream.r
}

func (h *TcpStream) run() {
	proto := h.channel.GetProtocol()
	src := h.channel.Configuration().Get("sip").MustString()
	r := bufio.NewReader(&h.r)

	for {
		var dataType int
		Sip := h.net.Src().String()
		if src == Sip {
			dataType = 1
		} else {
			dataType = 2
		}
		var n = 0
		var err error = nil
		buf := make([]byte, 2048)
		if n, err = r.Read(buf); err != nil {
			log.Error("telnetProtocol:reading data failed: ", err)
			return
		}

		pkt, err := proto.Decode(buf[0:n], dataType, h.newdata)
		if err != nil {
			tcpreader.DiscardBytesToEOF(r)
			if err == io.EOF {
				return
			}
		} else if pkt == nil {
			log.Info("packet is nil")
		} else {
			h.channel.DataAvailable(proto, pkt, h.newdata, err)
			if !pkt.Empty() {
				h.newdata = false
			}
		}
	}
}
