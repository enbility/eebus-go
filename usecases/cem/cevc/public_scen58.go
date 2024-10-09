package cevc

import (
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/spine-go/model"
)

// Scenario 5 & 6

// start sending heartbeat from the local CEM entity
//
// the heartbeat is started by default when a non 0 timeout is set in the service configuration
func (e *CEVC) StartHeartbeat() {
	if hm := e.LocalEntity.HeartbeatManager(); hm != nil {
		_ = hm.StartHeartbeat()
	}
}

// stop sending heartbeat from the local CEM entity
func (e *CEVC) StopHeartbeat() {
	if hm := e.LocalEntity.HeartbeatManager(); hm != nil {
		hm.StopHeartbeat()
	}
}

// Scenario 7 & 8

// set the local operating state of the local cem entity
//
// parameters:
//   - failureState: if true, the operating state is set to failure, otherwise to normal
func (e *CEVC) SetOperatingState(failureState bool) error {
	lf, err := server.NewDeviceDiagnosis(e.LocalEntity)
	if err != nil {
		return err
	}

	state := model.DeviceDiagnosisOperatingStateTypeNormalOperation
	if failureState {
		state = model.DeviceDiagnosisOperatingStateTypeFailure
	}
	lf.SetLocalOperatingState(state)

	return nil
}
