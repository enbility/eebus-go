package features

import (
	"testing"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

func TestDeviceDiagnosis_GetState(t *testing.T) {
	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := spine.NewEntityLocalImpl(localDevice, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	localDevice.AddEntity(localEntity)

	f := spine.NewFeatureLocalImpl(1, localEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	localEntity.AddFeature(f)

	remoteDeviceName := "remoteDevice"
	remoteDevice := spine.NewDeviceRemoteImpl(localDevice, "test", nil)
	data := &model.NodeManagementDetailedDiscoveryDataType{
		DeviceInformation: &model.NodeManagementDetailedDiscoveryDeviceInformationType{
			Description: &model.NetworkManagementDeviceDescriptionDataType{
				DeviceAddress: &model.DeviceAddressType{
					Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
				},
			},
		},
		EntityInformation: []model.NodeManagementDetailedDiscoveryEntityInformationType{
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{1},
					},
					EntityType: util.Ptr(model.EntityTypeTypeEVSE),
				},
			},
		},
		FeatureInformation: []model.NodeManagementDetailedDiscoveryFeatureInformationType{
			{
				Description: &model.NetworkManagementFeatureDescriptionDataType{
					FeatureAddress: &model.FeatureAddressType{
						Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity:  []model.AddressEntityType{1},
						Feature: util.Ptr(model.AddressFeatureType(1)),
					},
					FeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
					Role:        util.Ptr(model.RoleTypeClient),
				},
			},
		},
	}
	remoteEntities, err := remoteDevice.AddEntityAndFeatures(true, data)
	assert.Nil(t, err)
	assert.NotNil(t, remoteEntities)
	assert.NotEqual(t, 0, len(remoteEntities))

	remoteEntity := remoteEntities[0]

	d, err := NewDeviceDiagnosis(model.RoleTypeServer, model.RoleTypeClient, localDevice, remoteEntity)
	assert.Nil(t, err)
	assert.NotNil(t, d)

	result, err := d.GetState()
	assert.NotNil(t, err)
	assert.Nil(t, result)

	rF := remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.DeviceDiagnosisStateDataType{
		OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
	}
	rF.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, fData, nil, nil)

	result, err = d.GetState()
	assert.Nil(t, err)
	assert.NotNil(t, result)

}
