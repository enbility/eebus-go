package server_test

import (
	"fmt"
	"testing"

	"github.com/enbility/eebus-go/api"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/mock"
)

const remoteSki string = "testremoteski"

func setupFeatures(
	eebusService api.ServiceInterface, t *testing.T) (
	spineapi.DeviceRemoteInterface,
	[]spineapi.EntityRemoteInterface) {
	localDevice := eebusService.LocalDevice()
	localEntity := localDevice.EntityForType(model.EntityTypeTypeCEM)

	f := spine.NewFeatureLocal(1, localEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(2, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(3, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(4, localEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(5, localEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(6, localEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitListData, true, true)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(7, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionParameterDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionPermittedValueSetListData, true, false)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionCharacteristicListData, true, true)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(8, localEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueListData, true, true)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(9, localEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceClassificationManufacturerData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceClassificationUserData, true, true)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(9, localEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisStateData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceDiagnosisHeartbeatData, true, true)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(10, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeMeasurementDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeMeasurementListData, true, false)
	localEntity.AddFeature(f)

	writeHandler := shipmocks.NewShipConnectionDataWriterInterface(t)
	writeHandler.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()
	sender := spine.NewSender(writeHandler)
	remoteDevice := spine.NewDeviceRemote(localDevice, remoteSki, sender)

	var remoteFeatures = []struct {
		featureType   model.FeatureTypeType
		role          model.RoleType
		supportedFcts []model.FunctionType
	}{
		{model.FeatureTypeTypeLoadControl,
			model.RoleTypeServer,
			[]model.FunctionType{
				model.FunctionTypeLoadControlLimitDescriptionListData,
				model.FunctionTypeLoadControlLimitConstraintsListData,
				model.FunctionTypeLoadControlLimitListData,
			},
		},
		{model.FeatureTypeTypeElectricalConnection,
			model.RoleTypeServer,
			[]model.FunctionType{
				model.FunctionTypeElectricalConnectionParameterDescriptionListData,
				model.FunctionTypeElectricalConnectionPermittedValueSetListData,
			},
		},
		{model.FeatureTypeTypeMeasurement,
			model.RoleTypeServer,
			[]model.FunctionType{
				model.FunctionTypeMeasurementDescriptionListData,
				model.FunctionTypeMeasurementListData,
			},
		},
		{model.FeatureTypeTypeDeviceClassification,
			model.RoleTypeServer,
			[]model.FunctionType{
				model.FunctionTypeDeviceClassificationManufacturerData,
				model.FunctionTypeDeviceClassificationUserData,
			},
		},
		{model.FeatureTypeTypeDeviceConfiguration,
			model.RoleTypeServer,
			[]model.FunctionType{
				model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData,
				model.FunctionTypeDeviceConfigurationKeyValueListData,
			},
		},
	}

	remoteDeviceName := "remote"

	var featureInformations []model.NodeManagementDetailedDiscoveryFeatureInformationType
	for index, feature := range remoteFeatures {
		supportedFcts := []model.FunctionPropertyType{}
		for _, fct := range feature.supportedFcts {
			supportedFct := model.FunctionPropertyType{
				Function: util.Ptr(fct),
				PossibleOperations: &model.PossibleOperationsType{
					Read: &model.PossibleOperationsReadType{},
				},
			}
			supportedFcts = append(supportedFcts, supportedFct)
		}

		featureInformation := model.NodeManagementDetailedDiscoveryFeatureInformationType{
			Description: &model.NetworkManagementFeatureDescriptionDataType{
				FeatureAddress: &model.FeatureAddressType{
					Device:  util.Ptr(model.AddressDeviceType(remoteDeviceName)),
					Entity:  []model.AddressEntityType{1, 1},
					Feature: util.Ptr(model.AddressFeatureType(index)),
				},
				FeatureType:       util.Ptr(feature.featureType),
				Role:              util.Ptr(feature.role),
				SupportedFunction: supportedFcts,
			},
		}
		featureInformations = append(featureInformations, featureInformation)
	}

	detailedData := &model.NodeManagementDetailedDiscoveryDataType{
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
			{
				Description: &model.NetworkManagementEntityDescriptionDataType{
					EntityAddress: &model.EntityAddressType{
						Device: util.Ptr(model.AddressDeviceType(remoteDeviceName)),
						Entity: []model.AddressEntityType{1, 1},
					},
					EntityType: util.Ptr(model.EntityTypeTypeEV),
				},
			},
		},
		FeatureInformation: featureInformations,
	}

	entities, err := remoteDevice.AddEntityAndFeatures(true, detailedData)
	if err != nil {
		fmt.Println(err)
	}

	localDevice.AddRemoteDeviceForSki(remoteSki, remoteDevice)

	return remoteDevice, entities
}
