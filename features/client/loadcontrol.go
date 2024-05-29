package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type LoadControl struct {
	*Feature

	*internal.LoadControlCommon
}

// Get a new LoadControl features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewLoadControl(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*LoadControl, error) {
	feature, err := NewFeature(model.FeatureTypeTypeLoadControl, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	lc := &LoadControl{
		Feature:           feature,
		LoadControlCommon: internal.NewRemoteLoadControl(feature.featureRemote),
	}

	return lc, nil
}

var _ api.LoadControlClientInterface = (*LoadControl)(nil)

// request FunctionTypeLoadControlLimitDescriptionListData from a remote device
func (l *LoadControl) RequestLimitDescriptions() (*model.MsgCounterType, error) {
	return l.requestData(model.FunctionTypeLoadControlLimitDescriptionListData, nil, nil)
}

// request FunctionTypeLoadControlLimitConstraintsListData from a remote device
func (l *LoadControl) RequestLimitConstraints() (*model.MsgCounterType, error) {
	return l.requestData(model.FunctionTypeLoadControlLimitConstraintsListData, nil, nil)
}

// request FunctionTypeLoadControlLimitListData from a remote device
func (l *LoadControl) RequestLimitData() (*model.MsgCounterType, error) {
	return l.requestData(model.FunctionTypeLoadControlLimitListData, nil, nil)
}

// write load control limits
// returns an error if this failed
func (l *LoadControl) WriteLimitData(
	data []model.LoadControlLimitDataType,
	deleteSelectors *model.LoadControlLimitListDataSelectorsType,
	deleteElements *model.LoadControlLimitDataElementsType,
) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	var filters []model.FilterType
	if deleteElements != nil && deleteSelectors != nil {
		delFilter := model.FilterType{
			CmdControl: &model.CmdControlType{
				Delete: &model.ElementTagType{},
			},
			LoadControlLimitListDataSelectors: deleteSelectors,
			LoadControlLimitDataElements:      deleteElements,
		}
		filters = append(filters, delFilter)
	}
	filters = append(filters, *model.NewFilterTypePartial())

	cmd := model.CmdType{
		Function: util.Ptr(model.FunctionTypeLoadControlLimitListData),
		Filter:   filters,
		LoadControlLimitListData: &model.LoadControlLimitListDataType{
			LoadControlLimitData: data,
		},
	}

	return l.remoteDevice.Sender().Write(l.featureLocal.Address(), l.featureRemote.Address(), cmd)
}
