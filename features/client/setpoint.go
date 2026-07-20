package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Setpoint struct {
	*Feature

	*internal.SetpointCommon
}

// Get a new Setpoint features helper
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
		SetpointCommon: internal.NewRemoteSetpoint(feature.featureRemote),
	}

	return sp, nil
}

// request FunctionTypeSetpointDescriptionListData from a remote device
func (s *Setpoint) RequestSetpointDescriptions(
	selector *model.SetpointDescriptionListDataSelectorsType,
	elements *model.SetpointDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return s.requestData(model.FunctionTypeSetpointDescriptionListData, selector, elements)
}

// request FunctionTypeSetpointConstraintsListData from a remote device
func (s *Setpoint) RequestSetpointConstraints(
	selector *model.SetpointConstraintsListDataSelectorsType,
	elements *model.SetpointConstraintsDataElementsType,
) (*model.MsgCounterType, error) {
	return s.requestData(model.FunctionTypeSetpointConstraintsListData, selector, elements)
}

// request FunctionTypeSetpointListData from a remote device
func (s *Setpoint) RequestSetpoints(
	selector *model.SetpointListDataSelectorsType,
	elements *model.SetpointDataElementsType,
) (*model.MsgCounterType, error) {
	return s.requestData(model.FunctionTypeSetpointListData, selector, elements)
}

// write the given setpoint data to the remote device,
// returns the message counter of the sent message
func (s *Setpoint) WriteSetpointListData(
	data []model.SetpointDataType,
) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	// the remote server has to advertise the write operation for this function
	operation := s.featureRemote.Operations()[model.FunctionTypeSetpointListData]
	if operation == nil || !operation.Write() {
		return nil, api.ErrNotSupported
	}

	cmd := model.CmdType{
		SetpointListData: &model.SetpointListDataType{
			SetpointData: data,
		},
	}

	return s.remoteDevice.Sender().Write(s.featureLocal.Address(), s.featureRemote.Address(), cmd)
}
