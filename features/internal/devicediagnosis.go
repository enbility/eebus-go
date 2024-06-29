package internal

import (
	"time"

	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type DeviceDiagnosisCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalDeviceDiagnosis(featureLocal spineapi.FeatureLocalInterface) *DeviceDiagnosisCommon {
	return &DeviceDiagnosisCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteDeviceDiagnosis(featureRemote spineapi.FeatureRemoteInterface) *DeviceDiagnosisCommon {
	return &DeviceDiagnosisCommon{
		featureRemote: featureRemote,
	}
}

var _ api.DeviceDiagnosisCommonInterface = (*DeviceDiagnosisCommon)(nil)

// get the current diagnosis state for an device entity
func (d *DeviceDiagnosisCommon) GetState() (*model.DeviceDiagnosisStateDataType, error) {
	function := model.FunctionTypeDeviceDiagnosisStateData

	data, err := featureDataCopyOfType[model.DeviceDiagnosisStateDataType](d.featureLocal, d.featureRemote, function)
	if err != nil {
		return nil, api.ErrDataNotAvailable
	}

	return data, nil
}

// check if the currently available heartbeat data is within a time duration
func (d *DeviceDiagnosisCommon) IsHeartbeatWithinDuration(duration time.Duration) bool {
	function := model.FunctionTypeDeviceDiagnosisHeartbeatData

	data, err := featureDataCopyOfType[model.DeviceDiagnosisHeartbeatDataType](d.featureLocal, d.featureRemote, function)
	if err != nil || data == nil || data.Timestamp == nil {
		return false
	}

	timeValue, err := data.Timestamp.GetTime()
	if err != nil {
		return false
	}

	diff := time.Now().UTC().Add(-1 * duration)

	return diff.Compare(timeValue.Local()) <= 0
}
