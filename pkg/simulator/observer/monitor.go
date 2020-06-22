package observer

import (
	"api/log"
	"api/pkg/simulator/base"
	"api/pkg/simulator/util"
	"encoding/json"
	"time"

	"github.com/gorilla/websocket"
)

type Data struct {
	Type    int    `json:"type"`
	Content string `json:"content"`
	Systime string `json:"systime"`
}

type Monitor struct {
	conn     *websocket.Conn
	lastTime time.Time
	//firstCmdCome bool
	//firstRspCome bool
}

func (this *Monitor) Name() string {
	return "monitor"
}

func (this *Monitor) DataAvailable(proto base.Protocol, ch base.Channel, pkt base.Packet, isFistItem bool, err error) {
	if this.conn != nil && pkt != nil {
		if !pkt.Empty() {
			var itemData = pkt.String()
			this.lastTime = time.Now()
			data := Data{
				Type:    pkt.Type(),
				Content: itemData,
				Systime: time.Now().Format("2006-01-02 15:04:05"),
			}
			d, _ := json.Marshal(data)
			err := this.conn.WriteMessage(1, d)
			if err != nil {

			}
		}
	} else {
		log.Info(" conn is nil or pkt is nil")
	}
}

type monitorFactory struct{}

func (this *monitorFactory) New(c *util.Config, ch base.Channel, conn *websocket.Conn) base.ChannelObserver {
	return &Monitor{conn: conn}
}

func init() {
	base.RegisterObserverFactory("monitor", &monitorFactory{})
}
