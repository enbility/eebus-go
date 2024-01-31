package features

import (
	"time"

	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type DeviceDiagnosis struct {
	*FeatureImpl
}

func NewDeviceDiagnosis(
	localRole, remoteRole model.RoleType,
	localEntity api.EntityLocalInterface,
	remoteEntity api.EntityRemoteInterface) (*DeviceDiagnosis, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeDeviceDiagnosis, localRole, remoteRole, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	dd := &DeviceDiagnosis{
		FeatureImpl: feature,
	}

	return dd, nil
}

// request DeviceDiagnosisStateData from a remote entity
func (d *DeviceDiagnosis) RequestState() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceDiagnosisStateData, nil, nil)
}

// get the current diagnosis state for an device entity
func (d *DeviceDiagnosis) GetState() (*model.DeviceDiagnosisStateDataType, error) {
	rData := d.featureRemote.DataCopy(model.FunctionTypeDeviceDiagnosisStateData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.DeviceDiagnosisStateDataType)
	if data == nil {
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
	rData := d.featureRemote.DataCopy(model.FunctionTypeDeviceDiagnosisHeartbeatData)
	if rData == nil {
		return false
	}

	data := rData.(*model.DeviceDiagnosisHeartbeatDataType)
	if data == nil || data.Timestamp == nil {
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
