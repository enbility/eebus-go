package internal

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type SetPointCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

// NewLocalSetPoint creates a new SetPointCommon helper for local entities
func NewLocalSetPoint(featureLocal spineapi.FeatureLocalInterface) *SetPointCommon {
	return &SetPointCommon{
		featureLocal: featureLocal,
	}
}

// NewRemoteSetPoint creates a new SetPointCommon helper for remote entities
func NewRemoteSetPoint(featureRemote spineapi.FeatureRemoteInterface) *SetPointCommon {
	return &SetPointCommon{
		featureRemote: featureRemote,
	}
}

// GetSetpointDescriptions returns the setpoint descriptions
func (s *SetPointCommon) GetSetpointDescriptions() ([]model.SetpointDescriptionDataType, error) {
	function := model.FunctionTypeSetpointDescriptionListData

	data, err := featureDataCopyOfType[model.SetpointDescriptionListDataType](s.featureLocal, s.featureRemote, function)
	if err != nil || data == nil || data.SetpointDescriptionData == nil {
		return nil, api.ErrDataNotAvailable
	}

	return data.SetpointDescriptionData, nil
}

// GetSetpointForId returns the setpoint data for a given setpoint ID
func (s *SetPointCommon) GetSetpointForId(
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

// GetSetpoints returns the setpoints
func (s *SetPointCommon) GetSetpoints() []model.SetpointDataType {
	function := model.FunctionTypeSetpointListData

	data, err := featureDataCopyOfType[model.SetpointListDataType](s.featureLocal, s.featureRemote, function)
	if err != nil || data == nil || data.SetpointData == nil {
		return []model.SetpointDataType{}
	}

	return data.SetpointData
}

// GetSetpointDataForFilter returns the setpoint data for a given filter
func (s *SetPointCommon) GetSetpointDataForFilter(
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

// GetSetpointConstraints returns the setpoints constraints.
func (s *SetPointCommon) GetSetpointConstraints() []model.SetpointConstraintsDataType {
	function := model.FunctionTypeSetpointConstraintsListData

	data, err := featureDataCopyOfType[model.SetpointConstraintsListDataType](s.featureLocal, s.featureRemote, function)
	if err != nil || data == nil || data.SetpointConstraintsData == nil {
		return []model.SetpointConstraintsDataType{}
	}

	return data.SetpointConstraintsData
}

// GetSetpointConstraintsForId returns the setpoint constraints for a given setpoint ID
func (s *SetPointCommon) GetSetpointConstraintsForId(
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

// GetSetpointConstraintsForFilter returns the setpoint constraints for a given filter
func (s *SetPointCommon) GetSetpointConstraintsForFilter(
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
