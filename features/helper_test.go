package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
)

type featureFunctions struct {
	featureType model.FeatureTypeType
	functions   []model.FunctionType
}

func setupFeatures(t assert.TestingT, dataCon spine.SpineDataConnection, featureFunctions []featureFunctions) (*spine.DeviceLocalImpl, *spine.EntityRemoteImpl) {
	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := spine.NewEntityLocalImpl(localDevice, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	localDevice.AddEntity(localEntity)

	for i, item := range featureFunctions {
		f := spine.NewFeatureLocalImpl(uint(i+1), localEntity, item.featureType, model.RoleTypeServer)
		localEntity.AddFeature(f)
	}

	remoteDeviceName := "remoteDevice"
	remoteDevice := spine.NewDeviceRemoteImpl(localDevice, "test", dataCon)
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
	}

	var features []model.NodeManagementDetailedDiscoveryFeatureInformationType
	for i, item := range featureFunctions {
		featureI := model.NodeManagementDetailedDiscoveryFeatureInformationType{
			Description: &model.NetworkManagementFeatureDescriptionDataType{
				FeatureAddress: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
					Entity:  []model.AddressEntityType{1},
					Feature: util.Ptr(model.AddressFeatureType(i + 1)),
				},
				FeatureType: util.Ptr(item.featureType),
				Role:        util.Ptr(model.RoleTypeClient),
			},
		}
		var supportedFcts []model.FunctionPropertyType
		for _, function := range item.functions {
			supportedFct := model.FunctionPropertyType{
				Function: util.Ptr(function),
				PossibleOperations: &model.PossibleOperationsType{
					Read:  &model.PossibleOperationsReadType{},
					Write: &model.PossibleOperationsWriteType{},
				},
			}

			supportedFcts = append(supportedFcts, supportedFct)
		}
		featureI.Description.SupportedFunction = supportedFcts
		features = append(features, featureI)
	}
	data.FeatureInformation = features

	remoteEntities, err := remoteDevice.AddEntityAndFeatures(true, data)
	assert.Nil(t, err)
	assert.NotNil(t, remoteEntities)
	assert.NotEqual(t, 0, len(remoteEntities))

	return localDevice, remoteEntities[0]
}
