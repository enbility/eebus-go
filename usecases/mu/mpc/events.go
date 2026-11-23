package mpc

import (
	spineapi "github.com/enbility/spine-go/api"
)

// handle SPINE events
func (e *MPC) HandleEvent(payload spineapi.EventPayload) {
	// No events expected as remote Monitoring Appliance has no server data, and writes are not supported by the Monitored Unit
}
