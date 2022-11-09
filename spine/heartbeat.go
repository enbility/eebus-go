package spine

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type HeartbeatSender struct {
	heartBeatNum                uint64 // see https://github.com/golang/go/issues/11891
	stopHeartbeatC              chan struct{}
	stopMux                     sync.Mutex
	senderAddr, destinationAddr *model.FeatureAddressType
	sender                      Sender
	heartBeatTimeout            *model.DurationType
}

func NewHeartbeatSender(sender Sender) *HeartbeatSender {
	h := &HeartbeatSender{
		sender: sender,
	}
	// default to 4 seconds timeout
	h.heartBeatTimeout = model.NewDurationType(time.Second * 4)

	return h
}

func (c *HeartbeatSender) StartHeartbeatSend(senderAddr, destinationAddr *model.FeatureAddressType) {
	// stop a already running heartbeat
	c.StopHeartbeat()

	c.senderAddr = senderAddr
	c.destinationAddr = destinationAddr

	c.stopHeartbeatC = make(chan struct{})

	go func() {
		c.sendHearbeat(c.stopHeartbeatC, 800*time.Millisecond)
	}()
}

func (c *HeartbeatSender) StopHeartbeat() {
	c.stopMux.Lock()
	defer c.stopMux.Unlock()

	if c.stopHeartbeatC != nil && !c.isHeartbeatClosed() {
		close(c.stopHeartbeatC)
	}
}

func (c *HeartbeatSender) heartbeatCmd(t time.Time) model.CmdType {
	timestamp := t.UTC().Format(time.RFC3339)
	cmd := model.CmdType{
		DeviceDiagnosisHeartbeatData: &model.DeviceDiagnosisHeartbeatDataType{
			Timestamp:        &timestamp,
			HeartbeatCounter: c.heartBeatCounter(),
			HeartbeatTimeout: c.heartBeatTimeout,
		},
	}

	return cmd
}

func (c *HeartbeatSender) SendHeartBeatData(requestHeader *model.HeaderType) error {
	// TODO is this all we need here?

	cmd := c.heartbeatCmd(time.Now())

	return c.sender.Reply(requestHeader, c.senderAddr, cmd)
}

func (c *HeartbeatSender) sendHearbeat(stopC chan struct{}, d time.Duration) {
	ticker := time.NewTicker(d)
	for {
		select {
		case <-ticker.C:

			if c.senderAddr == nil || c.destinationAddr == nil {
				break
			}

			cmd := c.heartbeatCmd(time.Now())

			err := c.sender.Notify(c.senderAddr, c.destinationAddr, cmd)
			if err != nil {
				logging.Log.Error("ERROR sending heartbeat: ", err)
			}
		case <-stopC:
			return
		}
	}
}

func (c *HeartbeatSender) isHeartbeatClosed() bool {
	select {
	case <-c.stopHeartbeatC:
		return true
	default:
	}

	return false
}

// TODO heartBeatCounter should be global on CEM level, not on connection level
func (c *HeartbeatSender) heartBeatCounter() *uint64 {
	i := atomic.AddUint64(&c.heartBeatNum, 1)
	return &i
}
