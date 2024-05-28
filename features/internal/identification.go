package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type IdentificationCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalIdentification(featureLocal spineapi.FeatureLocalInterface) *IdentificationCommon {
	return &IdentificationCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteIdentification(featureRemote spineapi.FeatureRemoteInterface) *IdentificationCommon {
	return &IdentificationCommon{
		featureRemote: featureRemote,
	}
}

var _ api.IdentificationCommonInterface = (*IdentificationCommon)(nil)

// check if spine.EventPayload Data contains identification data
//
// data type will be checked for model.IdentificationListDataType
func (i *IdentificationCommon) CheckEventPayloadDataForFilter(payloadData any) bool {
	if payloadData == nil {
		return false
	}

	data, ok := payloadData.(*model.IdentificationListDataType)
	if !ok {
		return false
	}

	for _, item := range data.IdentificationData {
		if item.IdentificationId == nil || item.IdentificationValue == nil {
			continue
		}

		return true
	}

	return false
}

// return current values for Identification
func (i *IdentificationCommon) GetDataForFilter(filter model.IdentificationDataType) ([]model.IdentificationDataType, error) {
	function := model.FunctionTypeIdentificationListData

	data, err := featureDataCopyOfType[model.IdentificationListDataType](i.featureLocal, i.featureRemote, function)
	if err != nil || data == nil || data.IdentificationData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.IdentificationDataType](data.IdentificationData, filter)
	return result, nil
}
