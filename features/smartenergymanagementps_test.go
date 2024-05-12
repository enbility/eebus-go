package features_test

import (
	"testing"

	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/util"
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
	s.smartenergymgmtps, err = features.NewSmartEnergyManagementPs(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.smartenergymgmtps)
}

func (s *SmartEnergyManagementPsSuite) Test_RequestValues() {
	counter, err := s.smartenergymgmtps.RequestValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *SmartEnergyManagementPsSuite) Test_WriteValues() {
	counter, err := s.smartenergymgmtps.WriteValues(nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), counter)

	data := &model.SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &model.PowerSequenceNodeScheduleInformationDataType{},
		Alternatives:            []model.SmartEnergyManagementPsAlternativesType{},
	}
	counter, err = s.smartenergymgmtps.WriteValues(data)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *SmartEnergyManagementPsSuite) Test_GetValues() {
	value, err := s.smartenergymgmtps.GetValues()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)

	s.addData()

	value, err = s.smartenergymgmtps.GetValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
}

func (s *SmartEnergyManagementPsSuite) addData() {
	rF := s.remoteEntity.FeatureOfAddress(util.Ptr(model.AddressFeatureType(1)))

	fData := &model.SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &model.PowerSequenceNodeScheduleInformationDataType{},
		Alternatives:            []model.SmartEnergyManagementPsAlternativesType{},
	}
	rF.UpdateData(model.FunctionTypeSmartEnergyManagementPsData, fData, nil, nil)
}
