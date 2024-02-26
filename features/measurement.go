package features

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type Measurement struct {
	*Feature
}

// Get a new Measurement features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewMeasurement(
	localEntity api.EntityLocalInterface,
	remoteEntity api.EntityRemoteInterface) (*Measurement, error) {
	feature, err := NewFeature(model.FeatureTypeTypeMeasurement, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	m := &Measurement{
		Feature: feature,
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
func (m *Measurement) RequestValues() (*model.MsgCounterType, error) {
	return m.requestData(model.FunctionTypeMeasurementListData, nil, nil)
}

// return list of descriptions
func (m *Measurement) GetDescriptions() ([]model.MeasurementDescriptionDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.MeasurementDescriptionListDataType](m.featureRemote, model.FunctionTypeMeasurementDescriptionListData)
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	return data.MeasurementDescriptionData, nil
}

// return a list of MeasurementDescriptionDataType for a given scope
func (m *Measurement) GetDescriptionsForScope(scope model.ScopeTypeType) ([]model.MeasurementDescriptionDataType, error) {
	data, err := m.GetDescriptions()
	if err != nil {
		return nil, err
	}

	var result []model.MeasurementDescriptionDataType
	for _, item := range data {
		if item.ScopeType != nil && *item.ScopeType == scope {
			result = append(result, item)
		}
	}

	if len(result) == 0 {
		return nil, ErrDataNotAvailable
	}

	return result, nil
}

// return current electrical description for a given measurementId
func (m *Measurement) GetDescriptionForMeasurementId(measurementId model.MeasurementIdType) (*model.MeasurementDescriptionDataType, error) {
	descriptions, err := m.GetDescriptions()
	if err != nil {
		return nil, err
	}

	for _, item := range descriptions {
		if item.MeasurementId == nil ||
			*item.MeasurementId != measurementId {
			continue
		}

		return &item, nil
	}

	return nil, ErrMetadataNotAvailable
}

// return current values for measurements
func (m *Measurement) GetValues() ([]model.MeasurementDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.MeasurementListDataType](m.featureRemote, model.FunctionTypeMeasurementListData)
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	return data.MeasurementData, nil
}

// return current values of a measurementId
//
// if nothing is found, it will return an error
func (m *Measurement) GetValueForMeasurementId(id model.MeasurementIdType) (float64, error) {
	values, err := m.GetValues()
	if err != nil {
		return 0, err
	}

	for _, item := range values {
		if item.MeasurementId == nil || item.Value == nil {
			continue
		}

		if *item.MeasurementId == id {
			return item.Value.GetValue(), nil
		}
	}

	return 0, ErrDataNotAvailable
}

// return current values of a defined measurementType, commodityType and scopeType
//
// if nothing is found, it will return an error
func (m *Measurement) GetValuesForTypeCommodityScope(measurement model.MeasurementTypeType, commodity model.CommodityTypeType, scope model.ScopeTypeType) ([]model.MeasurementDataType, error) {
	values, err := m.GetValues()
	if err != nil {
		return nil, err
	}

	var resultSet []model.MeasurementDataType
	for _, item := range values {
		if item.MeasurementId == nil || item.Value == nil {
			continue
		}

		desc, err := m.GetDescriptionForMeasurementId(*item.MeasurementId)
		if err != nil ||
			desc.MeasurementType == nil || *desc.MeasurementType != measurement ||
			desc.CommodityType == nil || *desc.CommodityType != commodity ||
			desc.ScopeType == nil || *desc.ScopeType != scope {
			continue
		}

		resultSet = append(resultSet, item)
	}

	if len(resultSet) == 0 {
		return nil, ErrDataNotAvailable
	}

	return resultSet, nil
}

// return measurement constraints
func (m *Measurement) GetConstraints() ([]model.MeasurementConstraintsDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.MeasurementConstraintsListDataType](m.featureRemote, model.FunctionTypeMeasurementConstraintsListData)
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	return data.MeasurementConstraintsData, nil
}
