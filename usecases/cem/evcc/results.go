package evcc

import (
	"fmt"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

func (e *EVCC) HandleResponse(responseMsg api.ResponseMessage) {
	// before SPINE 1.3 the heartbeats are on the EVSE entity
	if responseMsg.EntityRemote == nil ||
		(responseMsg.EntityRemote.EntityType() != model.EntityTypeTypeEV &&
			responseMsg.EntityRemote.EntityType() != model.EntityTypeTypeEVSE) {
		return
	}

	// handle errors coming from the remote EVSE entity
	if responseMsg.FeatureLocal.Type() == model.FeatureTypeTypeDeviceDiagnosis {
		e.handleResultDeviceDiagnosis(responseMsg)
	}
}

// Handle DeviceDiagnosis Results
func (e *EVCC) handleResultDeviceDiagnosis(responseMsg api.ResponseMessage) {
	// is this an error for a heartbeat message?
	if responseMsg.DeviceRemote == nil ||
		responseMsg.Data == nil {
		return
	}

	result, ok := responseMsg.Data.(*model.ResultDataType)
	if !ok {
		return
	}

	if result.ErrorNumber == nil ||
		*result.ErrorNumber == model.ErrorNumberTypeNoError {
		return
	}

	// check if this is for a cached notify message
	datagram, err := responseMsg.DeviceRemote.Sender().DatagramForMsgCounter(responseMsg.MsgCounterReference)
	if err != nil {
		return
	}

	if len(datagram.Payload.Cmd) > 0 &&
		datagram.Payload.Cmd[0].DeviceDiagnosisHeartbeatData != nil {
		// something is horribly wrong, disconnect and hope a new connection will fix it
		errorText := fmt.Sprintf("Error Code: %d", result.ErrorNumber)
		if result.Description != nil {
			errorText = fmt.Sprintf("%s - %s", errorText, string(*result.Description))
		}
		e.service.DisconnectSKI(responseMsg.DeviceRemote.Ski(), errorText)
	}
}
