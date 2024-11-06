package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Setpoint struct {
	*Feature

	*internal.SetPointCommon
}

// Get a new SetPoint features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewSetpoint(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
) (*Setpoint, error) {
	feature, err := NewFeature(model.FeatureTypeTypeSetpoint, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	sp := &Setpoint{
		Feature:        feature,
		SetPointCommon: internal.NewRemoteSetPoint(feature.featureRemote),
	}

	return sp, nil
}

// request FunctionTypeSetpointDescriptionListData from a remote device
func (s *Setpoint) RequestSetPointDescriptions(
	selector *model.SetpointDescriptionListDataSelectorsType,
	elements *model.SetpointDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return s.requestData(model.FunctionTypeSetpointDescriptionListData, selector, elements)
}

// request FunctionTypeSetpointConstraintsListData from a remote device
func (s *Setpoint) RequestSetPointConstraints(
	selector *model.SetpointConstraintsListDataSelectorsType,
	elements *model.SetpointConstraintsDataElementsType,
) (*model.MsgCounterType, error) {
	return s.requestData(model.FunctionTypeSetpointConstraintsListData, selector, elements)
}

// request FunctionTypeSetpointListData from a remote device
func (s *Setpoint) RequestSetPoints(
	selector *model.SetpointListDataSelectorsType,
	elements *model.SetpointDataElementsType,
) (*model.MsgCounterType, error) {
	return s.requestData(model.FunctionTypeSetpointListData, selector, elements)
}

// WriteSetPointListData writes the given setpoint data
//
// Parameters:
// - data: the setpoint data to write
//
// Returns:
// - the message counter of the sent message
// - an error if the data could not be written
func (s *Setpoint) WriteSetPointListData(
	data []model.SetpointDataType,
) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	cmd := model.CmdType{
		SetpointListData: &model.SetpointListDataType{
			SetpointData: data,
		},
	}

	return s.remoteDevice.Sender().Write(s.featureLocal.Address(), s.featureRemote.Address(), cmd)
}
