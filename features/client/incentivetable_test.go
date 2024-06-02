package client

import (
	"testing"

	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestIncentiveTableSuite(t *testing.T) {
	suite.Run(t, new(IncentiveTableSuite))
}

type IncentiveTableSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	incentiveTable *IncentiveTable
	sentMessage    []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*IncentiveTableSuite)(nil)

func (s *IncentiveTableSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *IncentiveTableSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
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
	s.incentiveTable, err = NewIncentiveTable(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.incentiveTable)

	s.incentiveTable, err = NewIncentiveTable(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.incentiveTable)
}

func (s *IncentiveTableSuite) Test_RequestDescriptions() {
	counter, err := s.incentiveTable.RequestDescriptions()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *IncentiveTableSuite) Test_RequestConstraints() {
	counter, err := s.incentiveTable.RequestConstraints()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *IncentiveTableSuite) Test_RequestValues() {
	counter, err := s.incentiveTable.RequestValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *IncentiveTableSuite) Test_WriteValues() {
	counter, err := s.incentiveTable.WriteValues(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.IncentiveTableType{}
	counter, err = s.incentiveTable.WriteValues(data)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data = []model.IncentiveTableType{
		{
			Tariff: &model.TariffDataType{
				TariffId: util.Ptr(model.TariffIdType(0)),
			},
			IncentiveSlot: []model.IncentiveTableIncentiveSlotType{
				{
					TimeInterval: &model.TimeTableDataType{
						StartTime: &model.AbsoluteOrRecurringTimeType{
							Relative: model.NewDurationType(0),
						},
					},
					Tier: []model.IncentiveTableTierType{
						{
							Tier: &model.TierDataType{
								TierId: util.Ptr(model.TierIdType(0)),
							},
							Boundary: []model.TierBoundaryDataType{
								{
									BoundaryId:         util.Ptr(model.TierBoundaryIdType(0)),
									LowerBoundaryValue: model.NewScaledNumberType(0),
								},
							},
							Incentive: []model.IncentiveDataType{
								{
									IncentiveId: util.Ptr(model.IncentiveIdType(1)),
									Value:       model.NewScaledNumberType(100),
								},
							},
						},
					},
				},
			},
		},
	}

	counter, err = s.incentiveTable.WriteValues(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *IncentiveTableSuite) Test_WriteDescriptions() {
	counter, err := s.incentiveTable.WriteDescriptions(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := []model.IncentiveTableDescriptionType{}
	counter, err = s.incentiveTable.WriteDescriptions(data)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data = []model.IncentiveTableDescriptionType{
		{
			TariffDescription: &model.TariffDescriptionDataType{
				TariffId: util.Ptr(model.TariffIdType(0)),
			},
			Tier: []model.IncentiveTableDescriptionTierType{
				{
					TierDescription: &model.TierDescriptionDataType{
						TierId:   util.Ptr(model.TierIdType(0)),
						TierType: util.Ptr(model.TierTypeTypeFixedCost),
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
	}

	counter, err = s.incentiveTable.WriteDescriptions(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
