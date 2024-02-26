package features

import (
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type DeviceDiagnosis struct {
	*Feature
}

// Get a new DeviceDiagnosis features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewDeviceDiagnosis(
	localEntity api.EntityLocalInterface,
	remoteEntity api.EntityRemoteInterface) (*DeviceDiagnosis, error) {
	feature, err := NewFeature(model.FeatureTypeTypeDeviceDiagnosis, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	dd := &DeviceDiagnosis{
		Feature: feature,
	}

	return dd, nil
}

// request DeviceDiagnosisStateData from a remote entity
func (d *DeviceDiagnosis) RequestState() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceDiagnosisStateData, nil, nil)
}

// get the current diagnosis state for an device entity
func (d *DeviceDiagnosis) GetState() (*model.DeviceDiagnosisStateDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.DeviceDiagnosisStateDataType](d.featureRemote, model.FunctionTypeDeviceDiagnosisStateData)
	if err != nil {
		return nil, ErrDataNotAvailable
	}

	return data, nil
}

func (d *DeviceDiagnosis) SetLocalState(operatingState *model.DeviceDiagnosisStateDataType) {
	d.featureLocal.SetData(model.FunctionTypeDeviceDiagnosisStateData, operatingState)
}

// request FunctionTypeDeviceDiagnosisHeartbeatData from a remote device
func (d *DeviceDiagnosis) RequestHeartbeat() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceDiagnosisHeartbeatData, nil, nil)
}

// check if the currently available heartbeat data is within a time duration
func (d *DeviceDiagnosis) IsHeartbeatWithinDuration(duration time.Duration) bool {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.DeviceDiagnosisHeartbeatDataType](d.featureRemote, model.FunctionTypeDeviceDiagnosisHeartbeatData)
	if err != nil || data == nil || data.Timestamp == nil {
		return false
	}

	timeValue, err := data.Timestamp.GetTime()
	if err != nil {
		return false
	}

	now := time.Now()
	diff := now.Sub(timeValue)

	return diff < duration
}
