package internal_test

import (
	"testing"

	"github.com/enbility/eebus-go/features/internal"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestIncentiveTableSuite(t *testing.T) {
	suite.Run(t, new(IncentiveTableSuite))
}

type IncentiveTableSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.IncentiveTableCommon
}

func (s *IncentiveTableSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeIncentiveTable,
				functions: []model.FunctionType{
					model.FunctionTypeIncentiveTableDescriptionData,
					model.FunctionTypeIncentiveTableConstraintsData,
					model.FunctionTypeIncentiveTableData,
				},
			},
		},
	)

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalIncentiveTable(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeIncentiveTable, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteIncentiveTable(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *IncentiveTableSuite) Test_GetValues() {
	data, err := s.localSut.GetData()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))
	data, err = s.remoteSut.GetData()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addData()

	data, err = s.localSut.GetData()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
	data, err = s.remoteSut.GetData()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

func (s *IncentiveTableSuite) Test_GetDescriptions() {
	filter := model.TariffDescriptionDataType{}
	data, err := s.localSut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addDescription()

	data, err = s.localSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

func (s *IncentiveTableSuite) Test_GetDescriptionsForScope() {
	filter := model.TariffDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeSimpleIncentiveTable),
	}
	data, err := s.localSut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addDescription()

	data, err = s.localSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
	data, err = s.remoteSut.GetDescriptionsForFilter(filter)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

func (s *IncentiveTableSuite) Test_GetConstraints() {
	data, err := s.localSut.GetConstraints()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))
	data, err = s.remoteSut.GetConstraints()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addConstraints()

	data, err = s.localSut.GetConstraints()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
	data, err = s.remoteSut.GetConstraints()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

// helpers

func (s *IncentiveTableSuite) addData() {
	fData := &model.IncentiveTableDataType{
		IncentiveTable: []model.IncentiveTableType{
			{
				Tariff: &model.TariffDataType{
					TariffId: util.Ptr(model.TariffIdType(0)),
				},
				IncentiveSlot: []model.IncentiveTableIncentiveSlotType{
					{
						TimeInterval: &model.TimeTableDataType{},
					},
				},
			},
		},
	}
	s.localFeature.UpdateData(model.FunctionTypeIncentiveTableData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeIncentiveTableData, fData, nil, nil)
}

func (s *IncentiveTableSuite) addDescription() {
	fData := &model.IncentiveTableDescriptionDataType{
		IncentiveTableDescription: []model.IncentiveTableDescriptionType{
			{
				TariffDescription: &model.TariffDescriptionDataType{
					TariffId:        util.Ptr(model.TariffIdType(0)),
					TariffWriteable: util.Ptr(true),
					UpdateRequired:  util.Ptr(true),
					ScopeType:       util.Ptr(model.ScopeTypeTypeSimpleIncentiveTable),
				},
				Tier: []model.IncentiveTableDescriptionTierType{
					{
						TierDescription: &model.TierDescriptionDataType{
							TierId:   util.Ptr(model.TierIdType(0)),
							TierType: util.Ptr(model.TierTypeTypeDynamicCost),
						},
						BoundaryDescription: []model.TierBoundaryDescriptionDataType{
							{
								BoundaryId:   util.Ptr(model.TierBoundaryIdType(0)),
								BoundaryType: util.Ptr(model.TierBoundaryTypeTypePowerBoundary),
								BoundaryUnit: util.Ptr(model.UnitOfMeasurementTypeW),
							},
						},
						IncentiveDescription: []model.IncentiveDescriptionDataType{
							{
								IncentiveId:   util.Ptr(model.IncentiveIdType(0)),
								IncentiveType: util.Ptr(model.IncentiveTypeTypeAbsoluteCost),
								Currency:      util.Ptr(model.CurrencyTypeEur),
							},
						},
					},
				},
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeIncentiveTableDescriptionData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeIncentiveTableDescriptionData, fData, nil, nil)
}

func (s *IncentiveTableSuite) addConstraints() {
	fData := &model.IncentiveTableConstraintsDataType{
		IncentiveTableConstraints: []model.IncentiveTableConstraintsType{
			{
				Tariff: &model.TariffDataType{
					TariffId: util.Ptr(model.TariffIdType(0)),
				},
				TariffConstraints: &model.TariffOverallConstraintsDataType{
					MaxTiersPerTariff:    util.Ptr(model.TierCountType(3)),
					MaxBoundariesPerTier: util.Ptr(model.TierBoundaryCountType(1)),
					MaxIncentivesPerTier: util.Ptr(model.IncentiveCountType(3)),
				},
				IncentiveSlotConstraints: &model.TimeTableConstraintsDataType{
					SlotCountMax: util.Ptr(model.TimeSlotCountType(24)),
				},
			},
		},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeIncentiveTableConstraintsData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeIncentiveTableConstraintsData, fData, nil, nil)
}
