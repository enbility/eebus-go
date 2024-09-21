package client

import (
	"github.com/enbility/eebus-go/api"
	internal "github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

type DeviceConfiguration struct {
	*Feature

	*internal.DeviceConfigurationCommon
}

// Get a new DeviceConfiguration features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewDeviceConfiguration(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*DeviceConfiguration, error) {
	feature, err := NewFeature(model.FeatureTypeTypeDeviceConfiguration, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	dc := &DeviceConfiguration{
		Feature:                   feature,
		DeviceConfigurationCommon: internal.NewRemoteDeviceConfiguration(feature.featureRemote),
	}

	return dc, nil
}

var _ api.DeviceConfigurationClientInterface = (*DeviceConfiguration)(nil)

// request DeviceConfigurationKeyValueDescriptionDataType from a remote entity
func (d *DeviceConfiguration) RequestKeyValueDescriptions(
	selector *model.DeviceConfigurationKeyValueDescriptionListDataSelectorsType,
	elements *model.DeviceConfigurationKeyValueDescriptionDataElementsType,
) (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, selector, elements)
}

// request DeviceConfigurationKeyValueListData from a remote entity
func (d *DeviceConfiguration) RequestKeyValues(
	selector *model.DeviceConfigurationKeyValueListDataSelectorsType,
	elements *model.DeviceConfigurationKeyValueDataElementsType,
) (*model.MsgCounterType, error) {
	return d.requestData(model.FunctionTypeDeviceConfigurationKeyValueListData, selector, elements)
}

// write key values
// returns an error if this failed
func (d *DeviceConfiguration) WriteKeyValues(data []model.DeviceConfigurationKeyValueDataType) (*model.MsgCounterType, error) {
	if len(data) == 0 {
		return nil, api.ErrMissingData
	}

	filters := []model.FilterType{*model.NewFilterTypePartial()}

	// does the remote server feature not support partials?
	operation := d.featureRemote.Operations()[model.FunctionTypeDeviceConfigurationKeyValueListData]
	if operation == nil || !operation.WritePartial() {
		filters = nil
		// we need to send all data
		updateData := &model.DeviceConfigurationKeyValueListDataType{
			DeviceConfigurationKeyValueData: data,
		}

		if mergedData, err := d.featureRemote.UpdateData(false, model.FunctionTypeDeviceConfigurationKeyValueListData, updateData, nil, nil); err == nil {
			data = mergedData.([]model.DeviceConfigurationKeyValueDataType)
		}
	}

	cmd := model.CmdType{
		DeviceConfigurationKeyValueListData: &model.DeviceConfigurationKeyValueListDataType{
			DeviceConfigurationKeyValueData: data,
		},
	}

	if filters != nil {
		cmd.Filter = filters
		cmd.Function = util.Ptr(model.FunctionTypeDeviceConfigurationKeyValueListData)
	}

	return d.remoteDevice.Sender().Write(d.featureLocal.Address(), d.featureRemote.Address(), cmd)
}
