package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type Measurement struct {
	*FeatureImpl
}

func NewMeasurement(localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*Measurement, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeMeasurement, localRole, remoteRole, spineLocalDevice, entity)
	if err != nil {
		return nil, err
	}

	m := &Measurement{
		FeatureImpl: feature,
	}

	return m, nil
}

// request FunctionTypeMeasurementDescriptionListData from a remote device
func (m *Measurement) RequestDescriptions() error {
	_, err := m.requestData(model.FunctionTypeMeasurementDescriptionListData, nil, nil)

	return err
}

// request FunctionTypeMeasurementConstraintsListData from a remote entity
func (m *Measurement) RequestConstraints() error {
	_, err := m.requestData(model.FunctionTypeMeasurementConstraintsListData, nil, nil)
	return err
}

// request FunctionTypeMeasurementListData from a remote entity
func (m *Measurement) RequestValues() (*model.MsgCounterType, error) {
	return m.requestData(model.FunctionTypeMeasurementListData, nil, nil)
}

// return list of descriptions
func (m *Measurement) GetDescriptions() ([]model.MeasurementDescriptionDataType, error) {
	rData := m.featureRemote.Data(model.FunctionTypeMeasurementDescriptionListData)
	if rData == nil {
		return nil, ErrMetadataNotAvailable
	}
	data := rData.(*model.MeasurementDescriptionListDataType)
	if data == nil {
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
	rData := m.featureRemote.Data(model.FunctionTypeMeasurementListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.MeasurementListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.MeasurementData, nil
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
	rData := m.featureRemote.Data(model.FunctionTypeMeasurementConstraintsListData)
	if rData == nil {
		return nil, ErrMetadataNotAvailable
	}

	data := rData.(*model.MeasurementConstraintsListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.MeasurementConstraintsData, nil
}
