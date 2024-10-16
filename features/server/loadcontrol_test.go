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

func TestLoadControlSuite(t *testing.T) {
	suite.Run(t, new(LoadControlSuite))
}

type LoadControlSuite struct {
	suite.Suite

	sut *server.LoadControl

	service api.ServiceInterface

	localEntity spineapi.EntityLocalInterface

	remoteDevice     spineapi.DeviceRemoteInterface
	remoteEntity     spineapi.EntityRemoteInterface
	mockRemoteEntity *spinemocks.EntityRemoteInterface
}

func (s *LoadControlSuite) BeforeTest(suiteName, testName string) {
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
	s.sut, err = server.NewLoadControl(nil)
	assert.NotNil(s.T(), err)

	s.sut, err = server.NewLoadControl(s.localEntity)
	assert.Nil(s.T(), err)
}

func (s *LoadControlSuite) Test_CheckEventPayloadDataForFilter() {
	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
		ScopeType:     util.Ptr(model.ScopeTypeTypeSelfConsumption),
	}

	exists := s.sut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	exists = s.sut.CheckEventPayloadDataForFilter(s.mockRemoteEntity, filter)
	assert.False(s.T(), exists)

	descData := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(0)),
				LimitCategory: filter.LimitCategory,
				MeasurementId: util.Ptr(model.MeasurementIdType(0)),
				LimitType:     filter.LimitType,
				ScopeType:     filter.ScopeType,
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
				LimitCategory: filter.LimitCategory,
				MeasurementId: util.Ptr(model.MeasurementIdType(1)),
				LimitType:     filter.LimitType,
				ScopeType:     filter.ScopeType,
			},
			{
				LimitId:       util.Ptr(model.LoadControlLimitIdType(2)),
				LimitCategory: filter.LimitCategory,
				MeasurementId: util.Ptr(model.MeasurementIdType(2)),
				LimitType:     filter.LimitType,
				ScopeType:     filter.ScopeType,
			},
		},
	}

	entity := s.service.LocalDevice().EntityForType(model.EntityTypeTypeCEM)
	feature := entity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	feature.SetData(model.FunctionTypeLoadControlLimitDescriptionListData, descData)

	exists = s.sut.CheckEventPayloadDataForFilter(nil, filter)
	assert.False(s.T(), exists)

	limitData := &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{},
	}

	exists = s.sut.CheckEventPayloadDataForFilter(limitData, filter)
	assert.False(s.T(), exists)

	limitData = &model.LoadControlLimitListDataType{
		LoadControlLimitData: []model.LoadControlLimitDataType{
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(0)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(1)),
				Value:   model.NewScaledNumberType(16),
			},
			{
				LimitId: util.Ptr(model.LoadControlLimitIdType(2)),
			},
		},
	}

	exists = s.sut.CheckEventPayloadDataForFilter(limitData, filter)
	assert.True(s.T(), exists)
}

func (s *LoadControlSuite) Test_Description() {
	data, err := s.sut.GetLimitDescriptionForId(model.LoadControlLimitIdType(100))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	desc := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}
	limitId1 := s.sut.AddLimitDescription(desc)
	assert.NotNil(s.T(), limitId1)

	data, err = s.sut.GetLimitDescriptionForId(*limitId1)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)

	desc = model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeProduce),
	}

	limitId2 := s.sut.AddLimitDescription(desc)
	assert.NotNil(s.T(), limitId2)

	limitData, err := s.sut.GetLimitDescriptionForId(*limitId2)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), limitData)
}

func (s *LoadControlSuite) Test_GetDescriptionsForFilter() {
	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}

	data, err := s.sut.GetLimitDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	feature := s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

	desc := &model.LoadControlLimitDescriptionListDataType{
		LoadControlLimitDescriptionData: []model.LoadControlLimitDescriptionDataType{
			{
				LimitId:        util.Ptr(model.LoadControlLimitIdType(0)),
				LimitType:      filter.LimitType,
				LimitCategory:  filter.LimitCategory,
				LimitDirection: filter.LimitDirection,
				ScopeType:      filter.ScopeType,
			},
		},
	}
	feature.SetData(model.FunctionTypeLoadControlLimitDescriptionListData, desc)

	data, err = s.sut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1, len(data))
	assert.NotNil(s.T(), data[0].LimitId)

	filter = model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeProduce),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}

	data, err = s.sut.GetLimitDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))
}

func (s *LoadControlSuite) Test_GetLimitData() {
	ids := []api.LoadControlLimitDataForID{
		{
			Id: model.LoadControlLimitIdType(100),
		},
	}

	err := s.sut.UpdateLimitDataForIds(ids)
	assert.NotNil(s.T(), err)

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
		LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
		LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
		ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
	}

	data := []api.LoadControlLimitDataForFilter{
		{
			Filter: filter,
		},
	}
	err = s.sut.UpdateLimitDataForFilters(data, nil, nil)
	assert.NotNil(s.T(), err)

	data = []api.LoadControlLimitDataForFilter{
		{
			Data: model.LoadControlLimitDataType{
				LimitId: util.Ptr(model.LoadControlLimitIdType(100)),
			},
			Filter: filter,
		},
	}
	err = s.sut.UpdateLimitDataForFilters(data, nil, nil)
	assert.NotNil(s.T(), err)

	limitId := s.sut.AddLimitDescription(filter)
	assert.NotNil(s.T(), limitId)

	descData, err := s.sut.GetLimitDescriptionForId(*limitId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), descData)

	result, err := s.sut.GetLimitDataForId(*limitId)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)

	data = []api.LoadControlLimitDataForFilter{
		{
			Data: model.LoadControlLimitDataType{
				LimitId:    limitId,
				Value:      model.NewScaledNumberType(16),
				TimePeriod: model.NewTimePeriodTypeWithRelativeEndTime(time.Minute * 2),
			},
			Filter: filter,
		},
	}
	err = s.sut.UpdateLimitDataForFilters(data, nil, nil)
	assert.Nil(s.T(), err)

	result, err = s.sut.GetLimitDataForId(*limitId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)

	result, err = s.sut.GetLimitDataForId(model.LoadControlLimitIdType(100))
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)

	dataIds := []api.LoadControlLimitDataForFilter{}
	deleteSelectors := &model.LoadControlLimitListDataSelectorsType{
		LimitId: limitId,
	}
	deleteElements := &model.LoadControlLimitDataElementsType{
		TimePeriod: &model.TimePeriodElementsType{},
	}
	err = s.sut.UpdateLimitDataForFilters(dataIds, deleteSelectors, deleteElements)
	assert.Nil(s.T(), err)

	result, err = s.sut.GetLimitDataForId(*limitId)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Nil(s.T(), result.TimePeriod)
}
