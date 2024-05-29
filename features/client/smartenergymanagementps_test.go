package client_test

import (
	"testing"

	features "github.com/enbility/eebus-go/features/client"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSmartEnergyManagementPsSuite(t *testing.T) {
	suite.Run(t, new(SmartEnergyManagementPsSuite))
}

type SmartEnergyManagementPsSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	smartenergymgmtps *features.SmartEnergyManagementPs
	sentMessage       []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*SmartEnergyManagementPsSuite)(nil)

func (s *SmartEnergyManagementPsSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *SmartEnergyManagementPsSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeSmartEnergyManagementPs,
				functions: []model.FunctionType{
					model.FunctionTypeSmartEnergyManagementPsData,
				},
			},
		},
	)

	var err error
	s.smartenergymgmtps, err = features.NewSmartEnergyManagementPs(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.smartenergymgmtps)

	s.smartenergymgmtps, err = features.NewSmartEnergyManagementPs(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.smartenergymgmtps)
}

func (s *SmartEnergyManagementPsSuite) Test_RequestData() {
	counter, err := s.smartenergymgmtps.RequestData()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *SmartEnergyManagementPsSuite) Test_WriteData() {
	counter, err := s.smartenergymgmtps.WriteData(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := &model.SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &model.PowerSequenceNodeScheduleInformationDataType{},
		Alternatives:            []model.SmartEnergyManagementPsAlternativesType{},
	}
	counter, err = s.smartenergymgmtps.WriteData(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
