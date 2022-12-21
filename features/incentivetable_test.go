package features

import (
	"testing"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestIncentiveTableSuite(t *testing.T) {
	suite.Run(t, new(IncentiveTableSuite))
}

type IncentiveTableSuite struct {
	suite.Suite

	localDevice  *spine.DeviceLocalImpl
	remoteEntity *spine.EntityRemoteImpl

	incentiveTable *IncentiveTable
	sentMessage    []byte
}

var _ spine.SpineDataConnection = (*IncentiveTableSuite)(nil)

func (s *IncentiveTableSuite) WriteSpineMessage(message []byte) {
	s.sentMessage = message
}

func (s *IncentiveTableSuite) BeforeTest(suiteName, testName string) {
	s.localDevice, s.remoteEntity = setupFeatures(
		s.T(),
		s,
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

	var err error
	s.incentiveTable, err = NewIncentiveTable(model.RoleTypeServer, model.RoleTypeClient, s.localDevice, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.incentiveTable)
}

func (s *IncentiveTableSuite) Test_RequestDescriptions() {
	err := s.incentiveTable.RequestDescriptions()
	assert.Nil(s.T(), err)
}

func (s *IncentiveTableSuite) Test_RequestConstraints() {
	err := s.incentiveTable.RequestConstraints()
	assert.Nil(s.T(), err)
}

func (s *IncentiveTableSuite) Test_RequestValues() {
	counter, err := s.incentiveTable.RequestValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *IncentiveTableSuite) Test_GetValues() {
	data, err := s.incentiveTable.GetValues()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addData()

	data, err = s.incentiveTable.GetValues()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

func (s *IncentiveTableSuite) Test_GetDescriptions() {
	data, err := s.incentiveTable.GetDescriptions()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addDescription()

	data, err = s.incentiveTable.GetDescriptions()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

func (s *IncentiveTableSuite) Test_GetDescriptionsForScope() {
	scope := model.ScopeTypeTypeSimpleIncentiveTable
	data, err := s.incentiveTable.GetDescriptionsForScope(scope)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addDescription()

	data, err = s.incentiveTable.GetDescriptionsForScope(scope)
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

func (s *IncentiveTableSuite) Test_GetConstraints() {
	data, err := s.incentiveTable.GetConstraints()
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), 0, len(data))

	s.addConstraints()

	data, err = s.incentiveTable.GetConstraints()
	assert.Nil(s.T(), err)
	assert.NotEqual(s.T(), nil, data)
}

// helpers

func (s *IncentiveTableSuite) addData() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))

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
	rF.UpdateData(model.FunctionTypeIncentiveTableData, fData, nil, nil)
}

func (s *IncentiveTableSuite) addDescription() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
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
	rF.UpdateData(model.FunctionTypeIncentiveTableDescriptionData, fData, nil, nil)
}

func (s *IncentiveTableSuite) addConstraints() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
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
	rF.UpdateData(model.FunctionTypeIncentiveTableConstraintsData, fData, nil, nil)
}
