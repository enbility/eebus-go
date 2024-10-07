package evcc

import (
	"fmt"
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	spinemocks "github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestCemEVCCSuite(t *testing.T) {
	suite.Run(t, new(CemEVCCSuite))
}

type CemEVCCSuite struct {
	suite.Suite

	sut *EVCC

	service api.ServiceInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	mockSender       *spinemocks.SenderInterface
	mockRemoteEntity *spinemocks.EntityRemoteInterface
	evEntity         spineapi.EntityRemoteInterface

	eventCalled bool
}

func (s *CemEVCCSuite) Event(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event api.EventType) {
	s.eventCalled = true
}

func (s *CemEVCCSuite) BeforeTest(suiteName, testName string) {
	s.eventCalled = false
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
	mockRemoteFeature.EXPECT().Address().Return(&model.FeatureAddressType{}).Maybe()
	mockRemoteFeature.EXPECT().Operations().Return(nil).Maybe()

	localEntity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	s.sut = NewEVCC(s.service, localEntity, s.Event)
	s.sut.AddFeatures()
	s.sut.AddUseCase()

	var entities []spineapi.EntityRemoteInterface
	s.remoteDevice, s.mockSender, entities = setupDevices(s.service, s.T())
	s.evEntity = entities[1]
}

const remoteSki string = "testremoteski"

func setupDevices(
	eebusService api.ServiceInterface, t *testing.T) (
	spineapi.DeviceRemoteInterface,
	*spinemocks.SenderInterface,
	[]spineapi.EntityRemoteInterface) {
	localDevice := eebusService.LocalDevice()

	mockSender := spinemocks.NewSenderInterface(t)
	defaultMsgCounter := model.MsgCounterType(100)
	mockSender.
		EXPECT().
		Request(mock.Anything, mock.Anything, mock.Anything, mock.Anything, mock.Anything).
		Return(&defaultMsgCounter, nil).
		Maybe()
	mockSender.
		EXPECT().
		Subscribe(mock.Anything, mock.Anything, mock.Anything).
		Return(&defaultMsgCounter, nil).
		Maybe()
	remoteDevice := spine.NewDeviceRemote(localDevice, remoteSki, mockSender)

	remoteDeviceName := "remote"

	var remoteFeatures = []struct {
		featureType   model.FeatureTypeType
		supportedFcts []model.FunctionType
	}{
		{model.FeatureTypeTypeDeviceConfiguration,
			[]model.FunctionType{
				model.FunctionTypeDeviceConfigurationKeyValueDescriptionListData,
				model.FunctionTypeDeviceConfigurationKeyValueListData,
			},
		},
		{model.FeatureTypeTypeIdentification,
			[]model.FunctionType{
				model.FunctionTypeIdentificationListData,
			},
		},
		{model.FeatureTypeTypeDeviceClassification,
			[]model.FunctionType{
				model.FunctionTypeDeviceClassificationManufacturerData,
			},
		},
		{model.FeatureTypeTypeElectricalConnection,
			[]model.FunctionType{
				model.FunctionTypeElectricalConnectionDescriptionListData,
				model.FunctionTypeElectricalConnectionParameterDescriptionListData,
				model.FunctionTypeElectricalConnectionPermittedValueSetListData,
			},
		},
		{model.FeatureTypeTypeDeviceDiagnosis,
			[]model.FunctionType{
				model.FunctionTypeDeviceDiagnosisStateData,
			},
		},
	}

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
				Role:              util.Ptr(model.RoleTypeServer),
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

	return remoteDevice, mockSender, entities
}
