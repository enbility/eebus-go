package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type SetpointCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

// NewLocalSetpoint creates a new SetpointCommon helper for local entities
func NewLocalSetpoint(featureLocal spineapi.FeatureLocalInterface) *SetpointCommon {
	return &SetpointCommon{
		featureLocal: featureLocal,
	}
}

// NewRemoteSetpoint creates a new SetpointCommon helper for remote entities
func NewRemoteSetpoint(featureRemote spineapi.FeatureRemoteInterface) *SetpointCommon {
	return &SetpointCommon{
		featureRemote: featureRemote,
	}
}

// GetSetpointDescriptionsForFilter returns the setpoint descriptions for a given filter
func (s *SetpointCommon) GetSetpointDescriptionsForFilter(
	filter model.SetpointDescriptionDataType,
) ([]model.SetpointDescriptionDataType, error) {
	function := model.FunctionTypeSetpointDescriptionListData

	data, err := featureDataCopyOfType[model.SetpointDescriptionListDataType](s.featureLocal, s.featureRemote, function)
	if err != nil || data == nil || data.SetpointDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.SetpointDescriptionDataType](data.SetpointDescriptionData, filter)

	return result, nil
}

// GetSetpointDescriptionForId returns the setpoint description for a given setpoint ID
func (s *SetpointCommon) GetSetpointDescriptionForId(
	id model.SetpointIdType,
) (*model.SetpointDescriptionDataType, error) {
	filter := model.SetpointDescriptionDataType{
		SetpointId: &id,
	}

	result, err := s.GetSetpointDescriptionsForFilter(filter)
	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

// GetSetpointDataForFilter returns the setpoint data for a given filter
func (s *SetpointCommon) GetSetpointDataForFilter(
	filter model.SetpointDataType,
) ([]model.SetpointDataType, error) {
	function := model.FunctionTypeSetpointListData

	data, err := featureDataCopyOfType[model.SetpointListDataType](s.featureLocal, s.featureRemote, function)
	if err != nil || data == nil || data.SetpointData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.SetpointDataType](data.SetpointData, filter)

	return result, nil
}

// GetSetpointForId returns the setpoint data for a given setpoint ID
func (s *SetpointCommon) GetSetpointForId(
	id model.SetpointIdType,
) (*model.SetpointDataType, error) {
	filter := model.SetpointDataType{
		SetpointId: &id,
	}

	result, err := s.GetSetpointDataForFilter(filter)
	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

// GetSetpointConstraintsForFilter returns the setpoint constraints for a given filter
func (s *SetpointCommon) GetSetpointConstraintsForFilter(
	filter model.SetpointConstraintsDataType,
) ([]model.SetpointConstraintsDataType, error) {
	function := model.FunctionTypeSetpointConstraintsListData

	data, err := featureDataCopyOfType[model.SetpointConstraintsListDataType](s.featureLocal, s.featureRemote, function)
	if err != nil || data == nil || data.SetpointConstraintsData == nil {
		return nil, api.ErrDataNotAvailable
	}

	result := searchFilterInList[model.SetpointConstraintsDataType](data.SetpointConstraintsData, filter)

	return result, nil
}

// GetSetpointConstraintsForId returns the setpoint constraints for a given setpoint ID
func (s *SetpointCommon) GetSetpointConstraintsForId(
	id model.SetpointIdType,
) (*model.SetpointConstraintsDataType, error) {
	filter := model.SetpointConstraintsDataType{
		SetpointId: &id,
	}

	result, err := s.GetSetpointConstraintsForFilter(filter)
	if err != nil || len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return util.Ptr(result[0]), nil
}

// CheckEventPayloadDataForFilter checks if the given event payload data contains
// setpoint data for the given setpoint ID
func (s *SetpointCommon) CheckEventPayloadDataForFilter(payloadData any, filter model.SetpointDataType) bool {
	if payloadData == nil {
		return false
	}

	data, ok := payloadData.(*model.SetpointListDataType)
	if !ok || data.SetpointData == nil {
		return false
	}

	result := searchFilterInList[model.SetpointDataType](data.SetpointData, filter)

	return len(result) > 0
}
