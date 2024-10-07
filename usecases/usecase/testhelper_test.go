package usecase

import (
	"fmt"
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	spinemocks "github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestUseCaseSuite(t *testing.T) {
	suite.Run(t, new(UseCaseSuite))
}

type UseCaseSuite struct {
	suite.Suite

	service api.ServiceInterface

	uc *UseCaseBase

	localEntity spineapi.EntityLocalInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	mockRemoteEntity *spinemocks.EntityRemoteInterface
	evseEntity       spineapi.EntityRemoteInterface
	monitoredEntity  spineapi.EntityRemoteInterface
}

func (s *UseCaseSuite) Event(
	ski string,
	Device spineapi.DeviceRemoteInterface,
	entity spineapi.EntityRemoteInterface,
	event api.EventType) {
}

const (
	useCaseUpdateEvent = "test-update-event"
	useCaseActor       = model.UseCaseActorTypeCEM
	useCaseName        = model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment
)

func (s *UseCaseSuite) BeforeTest(suiteName, testName string) {
	cert, _ := cert.CreateCertificate("test", "test", "DE", "test")
	configuration, _ := api.NewConfiguration(
		"test", "test", "test", "test",
		[]shipapi.DeviceCategoryType{shipapi.DeviceCategoryTypeEnergyManagementSystem},
		model.DeviceTypeTypeEnergyManagementSystem,
		[]model.EntityTypeType{model.EntityTypeTypeCEM},
		9999, cert, time.Second*4)

	serviceHandler := mocks.NewServiceReaderInterface(s.T())
	serviceHandler.EXPECT().ServicePairingDetailUpdate(mock.Anything, mock.Anything).Return().Maybe()

	s.service = service.NewService(configuration, serviceHandler)
	_ = s.service.Setup()

	mockRemoteDevice := spinemocks.NewDeviceRemoteInterface(s.T())
	s.mockRemoteEntity = spinemocks.NewEntityRemoteInterface(s.T())
	mockRemoteFeature := spinemocks.NewFeatureRemoteInterface(s.T())
	mockRemoteDevice.EXPECT().FeatureByEntityTypeAndRole(mock.Anything, mock.Anything, mock.Anything).Return(mockRemoteFeature).Maybe()
	mockRemoteDevice.EXPECT().Ski().Return(remoteSki).Maybe()
	s.mockRemoteEntity.EXPECT().Device().Return(mockRemoteDevice).Maybe()
	s.mockRemoteEntity.EXPECT().EntityType().Return(mock.Anything).Maybe()
	entityAddress := &model.EntityAddressType{}
	s.mockRemoteEntity.EXPECT().Address().Return(entityAddress).Maybe()
	mockRemoteFeature.EXPECT().DataCopy(mock.Anything).Return(mock.Anything).Maybe()

	var entities []spineapi.EntityRemoteInterface

	s.localEntity, s.remoteDevice, entities = setupDevices(s.service, s.T())
	s.evseEntity = entities[0]
	s.monitoredEntity = entities[1]
	s.remoteDevice.AddEntity(s.monitoredEntity)

	validActorTypes := []model.UseCaseActorType{model.UseCaseActorTypeEV}
	validEntityTypes := []model.EntityTypeType{model.EntityTypeTypeEV}
	scenarios := []api.UseCaseScenario{
		{
			Scenario:       1,
			ServerFeatures: []model.FeatureTypeType{model.FeatureTypeTypeLoadControl},
		},
		{
			Scenario: 2,
		},
		{
			Scenario: 3,
		},
	}
	s.uc = NewUseCaseBase(
		s.localEntity,
		useCaseActor,
		useCaseName,
		"1.0.0",
		"release",
		scenarios,
		s.Event,
		useCaseUpdateEvent,
		validActorTypes,
		validEntityTypes,
	)
}

const remoteSki string = "testremoteski"

func setupDevices(
	eebusService api.ServiceInterface, t *testing.T) (
	spineapi.EntityLocalInterface,
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
	f = spine.NewFeatureLocal(1, localEntity, model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeLoadControlLimitListData, true, true)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(2, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionParameterDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionPermittedValueSetListData, true, false)
	f.AddFunctionType(model.FunctionTypeElectricalConnectionCharacteristicListData, true, true)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(3, localEntity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceConfigurationKeyValueListData, true, true)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocal(4, localEntity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	f.AddFunctionType(model.FunctionTypeDeviceClassificationManufacturerData, true, false)
	f.AddFunctionType(model.FunctionTypeDeviceClassificationUserData, true, true)
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
	remoteDevice.UpdateDevice(detailedData.DeviceInformation.Description)

	for _, entity := range entities {
		entity.UpdateDeviceAddress(*remoteDevice.Address())
	}

	localDevice.AddRemoteDeviceForSki(remoteSki, remoteDevice)

	return localEntity, remoteDevice, entities
}
