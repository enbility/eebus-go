package server_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/server"
	"github.com/enbility/eebus-go/mocks"
	"github.com/enbility/eebus-go/service"
	shipapi "github.com/enbility/ship-go/api"
	"github.com/enbility/ship-go/cert"
	spineapi "github.com/enbility/spine-go/api"
	spinemocks "github.com/enbility/spine-go/mocks"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestElectricalConnectionSuite(t *testing.T) {
	suite.Run(t, new(ElectricalConnectionSuite))
}

type ElectricalConnectionSuite struct {
	suite.Suite

	sut *server.ElectricalConnection

	service api.ServiceInterface

	localEntity spineapi.EntityLocalInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	remoteEntity     spineapi.EntityRemoteInterface
	mockRemoteEntity *spinemocks.EntityRemoteInterface
}

func (s *ElectricalConnectionSuite) BeforeTest(suiteName, testName string) {
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
	s.localEntity = s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)

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

	s.remoteDevice, entities = setupFeatures(s.service, s.T())
	s.remoteEntity = entities[1]

	var err error
	s.sut, err = server.NewElectricalConnection(nil)
	assert.NotNil(s.T(), err)

	s.sut, err = server.NewElectricalConnection(s.localEntity)
	assert.Nil(s.T(), err)
}

func (s *ElectricalConnectionSuite) Test_Description() {
	filter := model.ElectricalConnectionDescriptionDataType{}

	data, err := s.sut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	desc := model.ElectricalConnectionDescriptionDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		PowerSupplyType:        util.Ptr(model.ElectricalConnectionVoltageTypeTypeAc),
		AcConnectedPhases:      util.Ptr(uint(3)),
		ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	err = s.sut.AddDescription(desc)
	assert.Nil(s.T(), err)

	filter.ElectricalConnectionId = desc.ElectricalConnectionId
	data, err = s.sut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	desc = model.ElectricalConnectionDescriptionDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(1)),
		PowerSupplyType:        util.Ptr(model.ElectricalConnectionVoltageTypeTypeDc),
		ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	err = s.sut.AddDescription(desc)
	assert.Nil(s.T(), err)

	desc = model.ElectricalConnectionDescriptionDataType{
		PowerSupplyType: util.Ptr(model.ElectricalConnectionVoltageTypeTypeDc),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	err = s.sut.AddDescription(desc)
	assert.NotNil(s.T(), err)

	filter.ElectricalConnectionId = desc.ElectricalConnectionId
	data, err = s.sut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_ParameterDescription() {
	filter := model.ElectricalConnectionParameterDescriptionDataType{}

	data, err := s.sut.GetParameterDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	desc := model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),

		MeasurementId: util.Ptr(model.MeasurementIdType(0)),
		ScopeType:     util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	pId := s.sut.AddParameterDescription(desc)
	assert.NotNil(s.T(), pId)

	filter.ElectricalConnectionId = desc.ElectricalConnectionId
	data, err = s.sut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	desc = model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
		ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}
	pId = s.sut.AddParameterDescription(desc)
	assert.NotNil(s.T(), pId)

	desc = model.ElectricalConnectionParameterDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}

	pId = s.sut.AddParameterDescription(desc)
	assert.Nil(s.T(), pId)

	filter.ElectricalConnectionId = desc.ElectricalConnectionId
	data, err = s.sut.GetParameterDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ElectricalConnectionSuite) Test_GetCharacteristicsForFilter() {
	filter := model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
		CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeApparentPowerConsumptionNominalMax),
	}

	result, err := s.sut.GetCharacteristicsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)

	charData := model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
		CharacteristicContext:  filter.CharacteristicContext,
		CharacteristicType:     filter.CharacteristicType,
		Value:                  model.NewScaledNumberType(10),
	}
	charId, err := s.sut.AddCharacteristic(charData)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), charId)

	result, err = s.sut.GetCharacteristicsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)

	charData.CharacteristicType = util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeContractualConsumptionNominalMax)
	charId, err = s.sut.AddCharacteristic(charData)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), charId)

	result, err = s.sut.GetCharacteristicsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)

	charData.CharacteristicId = util.Ptr(model.ElectricalConnectionCharacteristicIdType(100))
	charId, err = s.sut.AddCharacteristic(charData)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), charId)

	charData.ElectricalConnectionId = nil
	charId, err = s.sut.AddCharacteristic(charData)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), charId)
}

func (s *ElectricalConnectionSuite) Test_UpdateCharacteristic() {
	charData := model.ElectricalConnectionCharacteristicDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		ParameterId:            util.Ptr(model.ElectricalConnectionParameterIdType(0)),
		CharacteristicContext:  util.Ptr(model.ElectricalConnectionCharacteristicContextTypeEntity),
		CharacteristicType:     util.Ptr(model.ElectricalConnectionCharacteristicTypeTypeApparentPowerConsumptionNominalMax),
		Value:                  model.NewScaledNumberType(10),
	}

	err := s.sut.UpdateCharacteristic(charData, nil)
	assert.NotNil(s.T(), err)

	charId, err := s.sut.AddCharacteristic(charData)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), charId)

	filter := model.ElectricalConnectionCharacteristicDataType{
		CharacteristicId: charId,
	}
	data, err := s.sut.GetCharacteristicsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), 1, len(data))
	assert.NotNil(s.T(), data[0].Value)
	assert.Equal(s.T(), 10.0, data[0].Value.GetValue())

	charData.CharacteristicId = util.Ptr(model.ElectricalConnectionCharacteristicIdType(100))
	charData.Value = model.NewScaledNumberType(20)
	err = s.sut.UpdateCharacteristic(charData, nil)
	assert.NotNil(s.T(), err)

	charData.CharacteristicId = charId
	err = s.sut.UpdateCharacteristic(charData, nil)
	assert.Nil(s.T(), err)

	deleteElements := &model.ElectricalConnectionCharacteristicDataElementsType{
		Value: &model.ScaledNumberElementsType{},
	}
	charData.Value = nil
	err = s.sut.UpdateCharacteristic(charData, deleteElements)
	assert.Nil(s.T(), err)

	data, err = s.sut.GetCharacteristicsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
	assert.Equal(s.T(), 1, len(data))
	assert.Nil(s.T(), data[0].Value)
}

func (s *ElectricalConnectionSuite) Test_PermittedData() {
	ids := []api.ElectricalConnectionPermittedValueSetForID{
		{
			ElectricalConnectionId: model.ElectricalConnectionIdType(0),
			ParameterId:            model.ElectricalConnectionParameterIdType(0),
			Data: model.ElectricalConnectionPermittedValueSetDataType{
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(10),
						},
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(0),
								Max: model.NewScaledNumberType(100),
							},
						},
					},
				},
			},
		},
	}

	err := s.sut.UpdatePermittedValueSetForIds(ids)
	assert.NotNil(s.T(), err)

	filter := model.ElectricalConnectionParameterDescriptionDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		MeasurementId:          util.Ptr(model.MeasurementIdType(0)),
		ScopeType:              util.Ptr(model.ScopeTypeTypeACPowerTotal),
	}

	data := []api.ElectricalConnectionPermittedValueSetForFilter{
		{
			Filter: filter,
		},
	}
	err = s.sut.UpdatePermittedValueSetForFilters(data, nil, nil)
	assert.NotNil(s.T(), err)

	data = []api.ElectricalConnectionPermittedValueSetForFilter{
		{
			Data: model.ElectricalConnectionPermittedValueSetDataType{
				ParameterId: util.Ptr(model.ElectricalConnectionParameterIdType(0)),
			},
			Filter: filter,
		},
	}
	err = s.sut.UpdatePermittedValueSetForFilters(data, nil, nil)
	assert.NotNil(s.T(), err)

	pdId := s.sut.AddParameterDescription(filter)
	assert.NotNil(s.T(), pdId)

	filter.ParameterId = pdId

	vsFilter := model.ElectricalConnectionPermittedValueSetDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		ParameterId:            pdId,
	}
	result, err := s.sut.GetPermittedValueSetForFilter(vsFilter)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)

	data = []api.ElectricalConnectionPermittedValueSetForFilter{
		{
			Data: model.ElectricalConnectionPermittedValueSetDataType{
				PermittedValueSet: []model.ScaledNumberSetType{
					{
						Value: []model.ScaledNumberType{
							*model.NewScaledNumberType(10),
						},
						Range: []model.ScaledNumberRangeType{
							{
								Min: model.NewScaledNumberType(0),
								Max: model.NewScaledNumberType(100),
							},
						},
					},
				},
			},
			Filter: filter,
		},
	}
	err = s.sut.UpdatePermittedValueSetForFilters(data, nil, nil)
	assert.Nil(s.T(), err)

	result, err = s.sut.GetPermittedValueSetForFilter(vsFilter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)

	vsFilter2 := model.ElectricalConnectionPermittedValueSetDataType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(101)),
		ParameterId:            pdId,
	}
	result, err = s.sut.GetPermittedValueSetForFilter(vsFilter2)
	assert.Nil(s.T(), err)
	assert.Nil(s.T(), result)

	dataFilter := []api.ElectricalConnectionPermittedValueSetForFilter{}
	deleteSelectors := &model.ElectricalConnectionPermittedValueSetListDataSelectorsType{
		ElectricalConnectionId: util.Ptr(model.ElectricalConnectionIdType(0)),
		ParameterId:            pdId,
	}
	deleteElements := &model.ElectricalConnectionPermittedValueSetDataElementsType{
		PermittedValueSet: &model.ElementTagType{},
	}
	err = s.sut.UpdatePermittedValueSetForFilters(dataFilter, deleteSelectors, deleteElements)
	assert.Nil(s.T(), err)

	result, err = s.sut.GetPermittedValueSetForFilter(vsFilter)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Nil(s.T(), result[0].PermittedValueSet)
}
