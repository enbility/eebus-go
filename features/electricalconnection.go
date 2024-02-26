package features

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type ElectricalConnection struct {
	*Feature
}

// Get a new ElectricalConnection features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewElectricalConnection(
	localEntity api.EntityLocalInterface,
	remoteEntity api.EntityRemoteInterface) (*ElectricalConnection, error) {
	feature, err := NewFeature(model.FeatureTypeTypeElectricalConnection, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	e := &ElectricalConnection{
		Feature: feature,
	}

	return e, nil
}

// request ElectricalConnectionDescriptionListDataType from a remote entity
func (e *ElectricalConnection) RequestDescriptions() (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionDescriptionListData, nil, nil)
}

// request FunctionTypeElectricalConnectionParameterDescriptionListData from a remote entity
func (e *ElectricalConnection) RequestParameterDescriptions() (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, nil, nil)
}

// request FunctionTypeElectricalConnectionPermittedValueSetListData from a remote entity
func (e *ElectricalConnection) RequestPermittedValueSets() (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, nil, nil)
}

// request FunctionTypeElectricalConnectionCharacteristicListData from a remote entity
func (e *ElectricalConnection) RequestCharacteristics() (*model.MsgCounterType, error) {
	return e.requestData(model.FunctionTypeElectricalConnectionCharacteristicListData, nil, nil)
}

// return list of description for Electrical Connection
func (e *ElectricalConnection) GetDescriptions() ([]model.ElectricalConnectionDescriptionDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.ElectricalConnectionDescriptionListDataType](e.featureRemote, model.FunctionTypeElectricalConnectionDescriptionListData)
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	return data.ElectricalConnectionDescriptionData, nil
}

// return current electrical description for a given measurementId
func (e *ElectricalConnection) GetDescriptionForMeasurementId(measurementId model.MeasurementIdType) (*model.ElectricalConnectionDescriptionDataType, error) {
	param, err := e.GetParameterDescriptionForMeasurementId(measurementId)
	if err != nil {
		return nil, err
	}

	descriptions, err := e.GetDescriptions()
	if err != nil {
		return nil, err
	}

	for _, item := range descriptions {
		if item.ElectricalConnectionId == nil ||
			param.ElectricalConnectionId == nil ||
			*item.ElectricalConnectionId != *param.ElectricalConnectionId {
			continue
		}

		return &item, nil
	}

	return nil, ErrMetadataNotAvailable
}

// return parameter descriptions for all Electrical Connections
func (e *ElectricalConnection) GetParameterDescriptions() ([]model.ElectricalConnectionParameterDescriptionDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.ElectricalConnectionParameterDescriptionListDataType](e.featureRemote, model.FunctionTypeElectricalConnectionParameterDescriptionListData)
	if err != nil {
		return nil, ErrDataNotAvailable
	}

	return data.ElectricalConnectionParameterDescriptionData, nil
}

// return parameter description for a specific scope
func (e *ElectricalConnection) GetParameterDescriptionForScopeType(scope model.ScopeTypeType) (*model.ElectricalConnectionParameterDescriptionDataType, error) {
	desc, err := e.GetParameterDescriptions()
	if err != nil {
		return nil, err
	}

	for _, element := range desc {
		if element.ScopeType == nil || *element.ScopeType != scope {
			continue
		}

		return &element, nil
	}

	return nil, ErrDataNotAvailable
}

// return parameter description for a specific parameterId
func (e *ElectricalConnection) GetParameterDescriptionForParameterId(parameterId model.ElectricalConnectionParameterIdType) (*model.ElectricalConnectionParameterDescriptionDataType, error) {
	desc, err := e.GetParameterDescriptions()
	if err != nil {
		return nil, err
	}

	for _, element := range desc {
		if element.ParameterId == nil || *element.ParameterId != parameterId {
			continue
		}

		return &element, nil
	}

	return nil, ErrDataNotAvailable
}

// return parameter description for a specific measurementId
func (e *ElectricalConnection) GetParameterDescriptionForMeasurementId(measurementId model.MeasurementIdType) (*model.ElectricalConnectionParameterDescriptionDataType, error) {
	desc, err := e.GetParameterDescriptions()
	if err != nil {
		return nil, err
	}

	for _, element := range desc {
		if element.MeasurementId == nil || *element.MeasurementId != measurementId {
			continue
		}

		return &element, nil
	}

	return nil, ErrDataNotAvailable
}

// return parameter description for a specific measurementId
func (e *ElectricalConnection) GetParameterDescriptionForMeasuredPhase(phase model.ElectricalConnectionPhaseNameType) (*model.ElectricalConnectionParameterDescriptionDataType, error) {
	desc, err := e.GetParameterDescriptions()
	if err != nil {
		return nil, err
	}

	for _, element := range desc {
		if element.AcMeasuredPhases == nil || *element.AcMeasuredPhases != phase {
			continue
		}

		return &element, nil
	}

	return nil, ErrDataNotAvailable
}

// return permitted values for all Electrical Connections
func (e *ElectricalConnection) GetPermittedValueSets() ([]model.ElectricalConnectionPermittedValueSetDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.ElectricalConnectionPermittedValueSetListDataType](e.featureRemote, model.FunctionTypeElectricalConnectionPermittedValueSetListData)
	if err != nil {
		return nil, ErrDataNotAvailable
	}

	return data.ElectricalConnectionPermittedValueSetData, nil
}

// return permitted valueset for a provided measuremnetId
func (e *ElectricalConnection) GetPermittedValueSetForParameterId(parameterId model.ElectricalConnectionParameterIdType) (*model.ElectricalConnectionPermittedValueSetDataType, error) {
	values, err := e.GetPermittedValueSets()
	if err != nil {
		return nil, err
	}

	for _, element := range values {
		if element.ParameterId == nil || *element.ParameterId != parameterId {
			continue
		}

		return &element, nil
	}

	return nil, ErrDataNotAvailable
}

// return permitted valueset for a provided measuremnetId
func (e *ElectricalConnection) GetPermittedValueSetForMeasurementId(measurementId model.MeasurementIdType) (*model.ElectricalConnectionPermittedValueSetDataType, error) {
	param, err := e.GetParameterDescriptionForMeasurementId(measurementId)
	if err != nil {
		return nil, err
	}

	values, err := e.GetPermittedValueSets()
	if err != nil {
		return nil, err
	}

	for _, element := range values {
		if element.ParameterId == nil || *element.ParameterId != *param.ParameterId {
			continue
		}

		return &element, nil
	}

	return nil, ErrDataNotAvailable
}

// returns minimum, maximum, default/pause limit values
func (e *ElectricalConnection) GetLimitsForParameterId(parameterId model.ElectricalConnectionParameterIdType) (float64, float64, float64, error) {
	data, err := e.GetPermittedValueSetForParameterId(parameterId)
	if err != nil || data.ElectricalConnectionId == nil || data.PermittedValueSet == nil {
		return 0, 0, 0, err
	}

	var resultMin, resultMax, resultDefault float64

	for _, set := range data.PermittedValueSet {
		if set.Value != nil && len(set.Value) > 0 {
			resultDefault = set.Value[0].GetValue()
		}
		if set.Range != nil {
			for _, rangeItem := range set.Range {
				if rangeItem.Min != nil {
					resultMin = rangeItem.Min.GetValue()
				}
				if rangeItem.Max != nil {
					resultMax = rangeItem.Max.GetValue()
				}
			}
		}
	}

	return resultMin, resultMax, resultDefault, nil
}

// Adjust a value to be within the permitted value range
func (e *ElectricalConnection) AdjustValueToBeWithinPermittedValuesForParameter(value float64, parameterId model.ElectricalConnectionParameterIdType) float64 {
	permittedValues, err := e.GetPermittedValueSetForParameterId(parameterId)
	if err != nil {
		return value
	}

	data := permittedValues.PermittedValueSet

	var defaultValue, minValue, maxValue float64
	var hasDefaultValue, hasRange bool

	for _, element := range data {
		// is a value set
		if element.Value != nil && len(element.Value) > 0 {
			defaultValue = element.Value[0].GetValue()
			hasDefaultValue = true
		}
		// is a range set
		if element.Range != nil && len(element.Range) > 0 {
			if element.Range[0].Min != nil {
				minValue = element.Range[0].Min.GetValue()
			}
			if element.Range[0].Max != nil {
				maxValue = element.Range[0].Max.GetValue()
			}
			hasRange = true
		}
	}

	if hasRange {
		if hasDefaultValue && value < minValue {
			value = defaultValue
		}
		if value > maxValue {
			value = maxValue
		}
	}

	return value
}

// return characteristics for a Electrical Connections
func (e *ElectricalConnection) GetCharacteristics() ([]model.ElectricalConnectionCharacteristicDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.ElectricalConnectionCharacteristicListDataType](e.featureRemote, model.FunctionTypeElectricalConnectionCharacteristicListData)
	if err != nil {
		return nil, ErrDataNotAvailable
	}

	return data.ElectricalConnectionCharacteristicListData, nil
}

// return characteristics for a Electrical Connections
func (e *ElectricalConnection) GetCharacteristicForContextType(
	context model.ElectricalConnectionCharacteristicContextType,
	cType model.ElectricalConnectionCharacteristicTypeType,
) (*model.ElectricalConnectionCharacteristicDataType, error) {
	data, err := e.GetCharacteristics()
	if err != nil || data == nil || len(data) == 0 {
		return nil, ErrDataNotAvailable
	}

	for _, item := range data {
		if item.CharacteristicId != nil &&
			item.CharacteristicContext != nil &&
			*item.CharacteristicContext == context &&
			item.CharacteristicType != nil &&
			*item.CharacteristicType == cType {
			return &item, nil
		}
	}

	return nil, ErrDataNotAvailable
}
