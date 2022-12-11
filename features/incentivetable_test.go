package features

import (
	"testing"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
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

func (s *IncentiveTableSuite) Test_RequestDescription() {
	err := s.incentiveTable.RequestDescription()
	assert.Nil(s.T(), err)
}

func (s *IncentiveTableSuite) Test_RequestConstraints() {
	err := s.incentiveTable.RequestConstraints()
	assert.Nil(s.T(), err)
}

func (s *IncentiveTableSuite) Test_RequestValues() {
	err := s.incentiveTable.RequestValues()
	assert.Nil(s.T(), err)
}
