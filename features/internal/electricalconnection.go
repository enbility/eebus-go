package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type ElectricalConnectionCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalElectricalConnection(featureLocal spineapi.FeatureLocalInterface) *ElectricalConnectionCommon {
	return &ElectricalConnectionCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteElectricalConnection(featureRemote spineapi.FeatureRemoteInterface) *ElectricalConnectionCommon {
	return &ElectricalConnectionCommon{
		featureRemote: featureRemote,
	}
}

var _ api.ElectricalConnectionCommonInterface = (*ElectricalConnectionCommon)(nil)

// check if spine.EventPayload Data contains data for a given filter
//
// data type will be checked for model.ElectricalConnectionPermittedValueSetListDataType,
// filter type will be checked for model.ElectricalConnectionParameterDescriptionDataType
func (e *ElectricalConnectionCommon) CheckEventPayloadDataForFilter(payloadData any, filter any) bool {
	if payloadData == nil {
		return false
	}

	data, ok := payloadData.(*model.ElectricalConnectionPermittedValueSetListDataType)
	filterData, ok2 := filter.(model.ElectricalConnectionParameterDescriptionDataType)
	if !ok || !ok2 {
		return false
	}

	descs, err := e.GetParameterDescriptionsForFilter(filterData)
	if err != nil {
		return false
	}

	for _, desc := range descs {
		for _, item := range data.ElectricalConnectionPermittedValueSetData {
			if item.ParameterId != nil &&
				desc.ParameterId != nil &&
				*item.ParameterId == *desc.ParameterId &&
				len(item.PermittedValueSet) != 0 {
				return true
			}
		}
	}

	return false
}

// Get the description for a given filter
//
// Returns an error if no matching description is found
func (e *ElectricalConnectionCommon) GetDescriptionsForFilter(
	filter model.ElectricalConnectionDescriptionDataType,
) ([]model.ElectricalConnectionDescriptionDataType, error) {
	function := model.FunctionTypeElectricalConnectionDescriptionListData

	data, err := featureDataCopyOfType[model.ElectricalConnectionDescriptionListDataType](e.featureLocal, e.featureRemote, function)
	if err != nil || data == nil || data.ElectricalConnectionDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.ElectricalConnectionDescriptionDataType](
		data.ElectricalConnectionDescriptionData, filter)
	return result, nil
}

// return current electrical description for a given parameter description
func (e *ElectricalConnectionCommon) GetDescriptionForParameterDescriptionFilter(
	filter model.ElectricalConnectionParameterDescriptionDataType) (
	*model.ElectricalConnectionDescriptionDataType, error) {
	param, err := e.GetParameterDescriptionsForFilter(filter)
	if err != nil || len(param) == 0 {
		return nil, err
	}

	descriptions, err := e.GetDescriptionsForFilter(model.ElectricalConnectionDescriptionDataType{})
	if err != nil {
		return nil, err
	}

	for _, item := range descriptions {
		if item.ElectricalConnectionId == nil ||
			param[0].ElectricalConnectionId == nil ||
			*item.ElectricalConnectionId != *param[0].ElectricalConnectionId {
			continue
		}

		return &item, nil
	}

	return nil, api.ErrMetadataNotAvailable
}

// Get the description for a given filter
//
// Returns an error if no matching description is found
func (e *ElectricalConnectionCommon) GetParameterDescriptionsForFilter(
	filter model.ElectricalConnectionParameterDescriptionDataType,
) ([]model.ElectricalConnectionParameterDescriptionDataType, error) {
	function := model.FunctionTypeElectricalConnectionParameterDescriptionListData

	data, err := featureDataCopyOfType[model.ElectricalConnectionParameterDescriptionListDataType](e.featureLocal, e.featureRemote, function)
	if err != nil || data == nil || data.ElectricalConnectionParameterDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.ElectricalConnectionParameterDescriptionDataType](
		data.ElectricalConnectionParameterDescriptionData, filter)
	return result, nil
}

// return permitted values for all Electrical Connections
func (e *ElectricalConnectionCommon) GetPermittedValueSetForFilter(
	filter model.ElectricalConnectionPermittedValueSetDataType) (
	[]model.ElectricalConnectionPermittedValueSetDataType, error) {
	function := model.FunctionTypeElectricalConnectionPermittedValueSetListData

	data, err := featureDataCopyOfType[model.ElectricalConnectionPermittedValueSetListDataType](e.featureLocal, e.featureRemote, function)
	if err != nil || data == nil || data.ElectricalConnectionPermittedValueSetData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.ElectricalConnectionPermittedValueSetDataType](data.ElectricalConnectionPermittedValueSetData, filter)
	return result, nil
}

// returns minimum, maximum, default/pause limit values
func (e *ElectricalConnectionCommon) GetPermittedValueDataForFilter(
	filter model.ElectricalConnectionPermittedValueSetDataType) (
	float64, float64, float64, error) {
	data, err := e.GetPermittedValueSetForFilter(filter)
	if err != nil || len(data) != 1 ||
		data[0].ElectricalConnectionId == nil ||
		data[0].PermittedValueSet == nil {
		return 0, 0, 0, api.ErrDataNotAvailable
	}

	var resultMin, resultMax, resultDefault float64

	for _, set := range data[0].PermittedValueSet {
		if len(set.Value) > 0 {
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

var PhaseNameMapping = []model.ElectricalConnectionPhaseNameType{
	model.ElectricalConnectionPhaseNameTypeA,
	model.ElectricalConnectionPhaseNameTypeB,
	model.ElectricalConnectionPhaseNameTypeC}

// Get the min, max, default current limits for each phase
func (e *ElectricalConnectionCommon) GetPhaseCurrentLimits(measurementDescs []model.MeasurementDescriptionDataType) (
	resultMin []float64, resultMax []float64, resultDefault []float64, resultErr error) {
	for _, phaseName := range PhaseNameMapping {
		// electricalParameterDescription contains the measured phase for each measurementId
		filter := model.ElectricalConnectionParameterDescriptionDataType{
			AcMeasuredPhases: util.Ptr(phaseName),
		}
		elParamDesc, err := e.GetParameterDescriptionsForFilter(filter)
		if err != nil || len(elParamDesc) == 0 {
			continue
		}

		// check all params and assume there are no phase specific power limits
		for _, paramEl := range elParamDesc {
			if paramEl.ParameterId == nil || paramEl.MeasurementId == nil {
				continue
			}

			// check if the measurementId is in measurementDescs
			found := false
			for _, mDesc := range measurementDescs {
				if mDesc.MeasurementId != nil && *mDesc.MeasurementId == *paramEl.MeasurementId {
					found = true
					break
				}
			}
			if !found {
				continue
			}

			filter1 := model.ElectricalConnectionPermittedValueSetDataType{
				ParameterId: paramEl.ParameterId,
			}
			dataMin, dataMax, dataDefault, err := e.GetPermittedValueDataForFilter(filter1)
			if err != nil {
				continue
			}

			// Min current data should be derived from min power data
			// but as this value is only properly provided via VAS the
			// currrent min values can not be trusted.

			resultMin = append(resultMin, dataMin)
			resultMax = append(resultMax, dataMax)
			resultDefault = append(resultDefault, dataDefault)
		}
	}

	if len(resultMin) == 0 {
		return nil, nil, nil, api.ErrDataNotAvailable
	}

	return resultMin, resultMax, resultDefault, nil
}

// Adjust a value to be within the permitted value range
func (e *ElectricalConnectionCommon) AdjustValueToBeWithinPermittedValuesForParameterId(
	value float64, parameterId model.ElectricalConnectionParameterIdType) float64 {
	filter := model.ElectricalConnectionPermittedValueSetDataType{
		ParameterId: &parameterId,
	}
	data, err := e.GetPermittedValueSetForFilter(filter)
	if err != nil || len(data) != 1 {
		return value
	}

	var defaultValue, minValue, maxValue float64
	var hasDefaultValue, hasRange bool

	for _, element := range data[0].PermittedValueSet {
		// is a value set
		if len(element.Value) > 0 {
			defaultValue = element.Value[0].GetValue()
			hasDefaultValue = true
		}
		// is a range set
		if len(element.Range) > 0 {
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

// Get the characteristics for a given filter
//
// Returns an error if no matching description is found
func (e *ElectricalConnectionCommon) GetCharacteristicsForFilter(
	filter model.ElectricalConnectionCharacteristicDataType,
) ([]model.ElectricalConnectionCharacteristicDataType, error) {
	function := model.FunctionTypeElectricalConnectionCharacteristicListData

	data, err := featureDataCopyOfType[model.ElectricalConnectionCharacteristicListDataType](e.featureLocal, e.featureRemote, function)
	if err != nil || data == nil || data.ElectricalConnectionCharacteristicData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.ElectricalConnectionCharacteristicDataType](data.ElectricalConnectionCharacteristicData, filter)
	return result, nil
}
