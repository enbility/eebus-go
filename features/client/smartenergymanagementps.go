package client

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type SmartEnergyManagementPs struct {
	*Feature

	*internal.SmartEnergyManagementPsCommon
}

// Get a new Identification features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewSmartEnergyManagementPs(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*SmartEnergyManagementPs, error) {
	feature, err := NewFeature(model.FeatureTypeTypeSmartEnergyManagementPs, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	i := &SmartEnergyManagementPs{
		Feature:                       feature,
		SmartEnergyManagementPsCommon: internal.NewRemoteSmartEnergyManagementPs(feature.featureRemote),
	}

	return i, nil
}

// request FunctionTypeSmartEnergyManagementPsData from a remote entity
func (i *SmartEnergyManagementPs) RequestData() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeSmartEnergyManagementPsData, nil, nil)
}

// write SmartEnergyManagementPsData
// returns an error if this failed
func (l *SmartEnergyManagementPs) WriteData(data *model.SmartEnergyManagementPsDataType) (*model.MsgCounterType, error) {
	if data == nil {
		return nil, api.ErrMissingData
	}

	cmd := model.CmdType{
		Function:                    util.Ptr(model.FunctionTypeSmartEnergyManagementPsData),
		Filter:                      []model.FilterType{*model.NewFilterTypePartial()},
		SmartEnergyManagementPsData: data,
	}

	return l.remoteDevice.Sender().Write(l.featureLocal.Address(), l.featureRemote.Address(), cmd)
}
