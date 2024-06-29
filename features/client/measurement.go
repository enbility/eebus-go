package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Measurement struct {
	*Feature

	*internal.MeasurementCommon
}

var _ api.MeasurementClientInterface = (*Measurement)(nil)

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
func (m *Measurement) RequestDescriptions(
	selector *model.MeasurementDescriptionListDataSelectorsType,
	elements *model.MeasurementDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return m.requestData(model.FunctionTypeMeasurementDescriptionListData, selector, elements)
}

// request FunctionTypeMeasurementConstraintsListData from a remote entity
func (m *Measurement) RequestConstraints(
	selector *model.MeasurementConstraintsListDataSelectorsType,
	elements *model.MeasurementConstraintsDataElementsType,
) (*model.MsgCounterType, error) {
	return m.requestData(model.FunctionTypeMeasurementConstraintsListData, selector, elements)
}

// request FunctionTypeMeasurementListData from a remote entity
func (m *Measurement) RequestData(
	selector *model.MeasurementListDataSelectorsType,
	elements *model.MeasurementDataElementsType,
) (*model.MsgCounterType, error) {
	return m.requestData(model.FunctionTypeMeasurementListData, selector, elements)
}
