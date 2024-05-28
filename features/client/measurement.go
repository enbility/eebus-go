package client

import (
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Measurement struct {
	*Feature

	*internal.MeasurementCommon
}

// Get a new Measurement features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewMeasurement(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*Measurement, error) {
	feature, err := NewFeature(model.FeatureTypeTypeMeasurement, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	m := &Measurement{
		Feature:           feature,
		MeasurementCommon: internal.NewRemoteMeasurement(feature.featureRemote),
	}

	return m, nil
}

// request FunctionTypeMeasurementDescriptionListData from a remote device
func (m *Measurement) RequestDescriptions() (*model.MsgCounterType, error) {
	return m.requestData(model.FunctionTypeMeasurementDescriptionListData, nil, nil)
}

// request FunctionTypeMeasurementConstraintsListData from a remote entity
func (m *Measurement) RequestConstraints() (*model.MsgCounterType, error) {
	return m.requestData(model.FunctionTypeMeasurementConstraintsListData, nil, nil)
}

// request FunctionTypeMeasurementListData from a remote entity
func (m *Measurement) RequestData() (*model.MsgCounterType, error) {
	return m.requestData(model.FunctionTypeMeasurementListData, nil, nil)
}
