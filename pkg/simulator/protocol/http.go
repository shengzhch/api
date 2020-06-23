package protocol

import (
	"api/log"
	"api/pkg/simulator/base"
	"api/pkg/simulator/util"
	"bufio"
	"bytes"
	"io"
	"net"
	"net/http"
)

// http protocol
type httpProtocol struct{}

var Req *http.Request

func (this *httpProtocol) Decode(buf []byte, datatype int, isnewdata bool) (base.Packet, error) {
	r := bytes.NewReader(buf)
	if datatype == 1 {
		req, err := http.ReadRequest(bufio.NewReader(r))
		if err != nil {
			log.Error("err", err)
			return nil, err
		}
		Req = req
		return &httpPacket{req: req, DataType: 1}, nil
	} else {
		resp, err1 := http.ReadResponse(bufio.NewReader(r), Req)
		if err1 != nil {
			return nil, io.EOF
		}
		return &httpPacket{resp: resp, DataType: 2}, nil
	}
}

func (this *httpProtocol) DealData(bytes []byte) (base.Packet, error) {
	return &httpPacket{}, nil
}

func (this *httpProtocol) Encode([]byte) ([]byte, error) { return nil, nil }
func (this *httpProtocol) Name() string                  { return "http" }
func (this *httpProtocol) GetDelim() byte                { return '\n' }
func (this *httpProtocol) ProcessOP(conn net.Conn)       {}

type httpProtocolFactory struct{}

func (this *httpProtocolFactory) New(c *util.Config) (base.Protocol, error) {
	return &httpProtocol{}, nil
}

// http packet
type httpPacket struct {
	req      *http.Request
	resp     *http.Response
	DataType int
}

func (this *httpPacket) String() string {
	if this.DataType == 0 {
		return this.req.URL.String()
	} else if this.DataType == 1 {
		var buffer bytes.Buffer
		_, err := buffer.ReadFrom(this.req.Body)
		if err != nil {
			log.Error("err", err)
			return ""
		}
		return buffer.String()
	}

	var buffer bytes.Buffer
	_, err := buffer.ReadFrom(this.resp.Body)
	if err != nil {
		log.Error("err", err)
		return ""
	}
	return buffer.String()
}

func (this *httpPacket) Bytes() []byte {
	return ([]byte)(this.String())
}

func (this *httpPacket) Empty() bool {
	return this.String() == ""
}

func (this *httpPacket) Type() int {
	return this.DataType
}

func init() {
	base.RegisterProtocolFactory("http", &httpProtocolFactory{})
}
