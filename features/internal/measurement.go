package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type MeasurementCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalMeasurement(featureLocal spineapi.FeatureLocalInterface) *MeasurementCommon {
	return &MeasurementCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteMeasurement(featureRemote spineapi.FeatureRemoteInterface) *MeasurementCommon {
	return &MeasurementCommon{
		featureRemote: featureRemote,
	}
}

var _ api.MeasurementCommonInterface = (*MeasurementCommon)(nil)

// check if spine.EventPayload Data contains data for a given filter
//
// data type will be checked for model.MeasurementListDataType,
// filter type will be checked for model.MeasurementDescriptionDataType
func (m *MeasurementCommon) CheckEventPayloadDataForFilter(payloadData any, filter any) bool {
	if payloadData == nil {
		return false
	}

	data, ok := payloadData.(*model.MeasurementListDataType)
	filterData, ok2 := filter.(model.MeasurementDescriptionDataType)
	if !ok || !ok2 {
		return false
	}

	descs, err := m.GetDescriptionsForFilter(filterData)
	if err != nil {
		return false
	}
	for _, desc := range descs {
		if desc.MeasurementId == nil {
			continue
		}

		for _, item := range data.MeasurementData {
			if item.MeasurementId != nil &&
				*item.MeasurementId == *desc.MeasurementId &&
				item.Value != nil {
				return true
			}
		}
	}

	return false
}

// Get the description for a given id
//
// Returns an error if no matching description is found
func (m *MeasurementCommon) GetDescriptionForId(
	measurementId model.MeasurementIdType,
) (*model.MeasurementDescriptionDataType, error) {
	data, err := m.GetDescriptionsForFilter(model.MeasurementDescriptionDataType{MeasurementId: &measurementId})

	if err != nil || len(data) != 1 {
		return nil, api.ErrDataNotAvailable
	}

	return &data[0], nil
}

// Get the description for a given filter
//
// Returns an error if no matching description is found
func (m *MeasurementCommon) GetDescriptionsForFilter(
	filter model.MeasurementDescriptionDataType,
) ([]model.MeasurementDescriptionDataType, error) {
	function := model.FunctionTypeMeasurementDescriptionListData

	data, err := featureDataCopyOfType[model.MeasurementDescriptionListDataType](m.featureLocal, m.featureRemote, function)
	if err != nil || data == nil || data.MeasurementDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.MeasurementDescriptionDataType](data.MeasurementDescriptionData, filter)
	return result, nil
}

// Get the constraints for a given filter
//
// Returns an error if no matching constraint is found
func (m *MeasurementCommon) GetConstraintsForFilter(
	filter model.MeasurementConstraintsDataType,
) ([]model.MeasurementConstraintsDataType, error) {
	function := model.FunctionTypeMeasurementConstraintsListData

	data, err := featureDataCopyOfType[model.MeasurementConstraintsListDataType](m.featureLocal, m.featureRemote, function)
	if err != nil || data == nil || data.MeasurementConstraintsData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.MeasurementConstraintsDataType](data.MeasurementConstraintsData, filter)
	return result, nil
}

// Get the measuement data for a given measurementId
//
// Will return nil if no data is available
func (m *MeasurementCommon) GetDataForId(measurementId model.MeasurementIdType) (
	*model.MeasurementDataType, error) {
	result, err := m.GetDataForFilter(model.MeasurementDescriptionDataType{MeasurementId: &measurementId})
	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return &result[0], nil
}

// Get measuement data for a given filter
//
// Will return nil if no data is available
func (m *MeasurementCommon) GetDataForFilter(filter model.MeasurementDescriptionDataType) (
	[]model.MeasurementDataType, error) {
	function := model.FunctionTypeMeasurementListData

	descriptions, err := m.GetDescriptionsForFilter(filter)
	if err != nil || len(descriptions) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	data, err := featureDataCopyOfType[model.MeasurementListDataType](m.featureLocal, m.featureRemote, function)
	if err != nil || data == nil || data.MeasurementData == nil {
		return nil, api.ErrDataNotAvailable
	}

	var result []model.MeasurementDataType

	for _, desc := range descriptions {
		filter2 := model.MeasurementDataType{
			MeasurementId: desc.MeasurementId,
		}

		elements := searchFilterInList[model.MeasurementDataType](data.MeasurementData, filter2)
		result = append(result, elements...)
	}
	return result, nil
}
