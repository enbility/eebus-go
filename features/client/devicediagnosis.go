package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type DeviceDiagnosis struct {
	*Feature

	*internal.DeviceDiagnosisCommon
}

// Get a new DeviceDiagnosis features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewDeviceDiagnosis(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*DeviceDiagnosis, error) {
	feature, err := NewFeature(model.FeatureTypeTypeDeviceDiagnosis, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	dd := &DeviceDiagnosis{
		Feature:               feature,
		DeviceDiagnosisCommon: internal.NewRemoteDeviceDiagnosis(feature.featureRemote),
	}

	return dd, nil
}

var _ api.DeviceDiagnosisClientInterface = (*DeviceDiagnosis)(nil)

// request DeviceDiagnosisStateData from a remote entity
func (d *DeviceDiagnosis) RequestState() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceDiagnosisStateData, nil, nil)
}

// request FunctionTypeDeviceDiagnosisHeartbeatData from a remote device
func (d *DeviceDiagnosis) RequestHeartbeat() (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceDiagnosisHeartbeatData, nil, nil)
}
