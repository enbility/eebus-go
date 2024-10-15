package server

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

func NewDeviceDiagnosis(localEntity spineapi.EntityLocalInterface) (*DeviceDiagnosis, error) {
	feature, err := NewFeature(model.FeatureTypeTypeDeviceDiagnosis, localEntity)
	if err != nil {
		return nil, err
	}

	dc := &DeviceDiagnosis{
		Feature:               feature,
		DeviceDiagnosisCommon: internal.NewLocalDeviceDiagnosis(feature.featureLocal),
	}

	return dc, nil
}

var _ api.DeviceDiagnosisServerInterface = (*DeviceDiagnosis)(nil)

// set the local device diagnosis state of the device
func (d *DeviceDiagnosis) SetLocalState(state *model.DeviceDiagnosisStateDataType) {
	d.featureLocal.SetData(model.FunctionTypeDeviceDiagnosisStateData, state)
}

// set the local device diagnosis operating state
func (d *DeviceDiagnosis) SetLocalOperatingState(operatingState model.DeviceDiagnosisOperatingStateType) {
	stateData := &model.DeviceDiagnosisStateDataType{
		OperatingState: &operatingState,
	}
	d.SetLocalState(stateData)
}
