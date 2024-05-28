package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type LoadControlCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalLoadControl(featureLocal spineapi.FeatureLocalInterface) *LoadControlCommon {
	return &LoadControlCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteLoadControl(featureRemote spineapi.FeatureRemoteInterface) *LoadControlCommon {
	return &LoadControlCommon{
		featureRemote: featureRemote,
	}
}

var _ api.FeatureServerInterface = (*LoadControlCommon)(nil)

// check if spine.EventPayload Data contains data for a given filter
//
// data type will be checked for model.LoadControlLimitListDataType,
// filter type will be checked for model.LoadControlLimitDescriptionDataType
func (l *LoadControlCommon) CheckEventPayloadDataForFilter(payloadData any, filter any) bool {
	if payloadData == nil {
		return false
	}

	data, ok := payloadData.(*model.LoadControlLimitListDataType)
	filterData, ok2 := filter.(model.LoadControlLimitDescriptionDataType)
	if !ok || !ok2 {
		return false
	}

	descs, _ := l.GetLimitDescriptionsForFilter(filterData)
	for _, desc := range descs {
		for _, item := range data.LoadControlLimitData {
			if item.LimitId != nil &&
				desc.LimitId != nil &&
				*item.LimitId == *desc.LimitId &&
				item.Value != nil {
				return true
			}
		}
	}

	return false
}

var _ api.LoadControlCommonInterface = (*LoadControlCommon)(nil)

// Get the description for a given limitId
//
// Will return nil if no matching description is found
func (l *LoadControlCommon) GetLimitDescriptionForId(limitId model.LoadControlLimitIdType) (*model.LoadControlLimitDescriptionDataType, error) {
	data, err := l.GetLimitDescriptionsForFilter(model.LoadControlLimitDescriptionDataType{LimitId: &limitId})

	if err != nil || len(data) != 1 {
		return nil, api.ErrDataNotAvailable
	}

	return &data[0], nil
}

// Get the description for a given filter
//
// Returns an error if no matching description is found
func (l *LoadControlCommon) GetLimitDescriptionsForFilter(
	filter model.LoadControlLimitDescriptionDataType,
) ([]model.LoadControlLimitDescriptionDataType, error) {
	function := model.FunctionTypeLoadControlLimitDescriptionListData

	data, err := featureDataCopyOfType[model.LoadControlLimitDescriptionListDataType](l.featureLocal, l.featureRemote, function)
	if err != nil || data == nil || data.LoadControlLimitDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.LoadControlLimitDescriptionDataType](data.LoadControlLimitDescriptionData, filter)
	return result, nil
}

// Get the description for a given limitId
//
// Will return nil if no data is available
func (l *LoadControlCommon) GetLimitDataForId(limitId model.LoadControlLimitIdType) (*model.LoadControlLimitDataType, error) {
	result, err := l.GetLimitDataForFilter(model.LoadControlLimitDescriptionDataType{LimitId: &limitId})
	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}
	return &result[0], nil
}

// Get limit data for a given filter
//
// Will return nil if no data is available
func (l *LoadControlCommon) GetLimitDataForFilter(filter model.LoadControlLimitDescriptionDataType) ([]model.LoadControlLimitDataType, error) {
	function := model.FunctionTypeLoadControlLimitListData

	descriptions, err := l.GetLimitDescriptionsForFilter(filter)
	if err != nil || len(descriptions) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	data, err := featureDataCopyOfType[model.LoadControlLimitListDataType](l.featureLocal, l.featureRemote, function)
	if err != nil || data == nil || data.LoadControlLimitData == nil {
		return nil, api.ErrDataNotAvailable
	}

	var result []model.LoadControlLimitDataType

	for _, desc := range descriptions {
		filter2 := model.LoadControlLimitDataType{
			LimitId: desc.LimitId,
		}

		elements := searchFilterInList[model.LoadControlLimitDataType](data.LoadControlLimitData, filter2)
		result = append(result, elements...)
	}
	return result, nil
}
