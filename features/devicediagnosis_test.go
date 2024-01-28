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

func TestDeviceDiagnosisSuite(t *testing.T) {
	suite.Run(t, new(DeviceDiagnosisSuite))
}

type DeviceDiagnosisSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	deviceDiagnosis *features.DeviceDiagnosis
	sentMessage     []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*DeviceDiagnosisSuite)(nil)

func (s *DeviceDiagnosisSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *DeviceDiagnosisSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeDeviceDiagnosis,
				functions: []model.FunctionType{
					model.FunctionTypeDeviceDiagnosisStateData,
				},
			},
		},
	)

	var err error
	s.deviceDiagnosis, err = features.NewDeviceDiagnosis(model.RoleTypeClient, model.RoleTypeServer, s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.deviceDiagnosis)
}

func (s *DeviceDiagnosisSuite) Test_RequestState() {
	counter, err := s.deviceDiagnosis.RequestState()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *DeviceDiagnosisSuite) Test_GetState() {
	result, err := s.deviceDiagnosis.GetState()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)

	rF := s.remoteEntity.FeatureOfAddress(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.DeviceDiagnosisStateDataType{
		OperatingState:       util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
		PowerSupplyCondition: util.Ptr(model.PowerSupplyConditionTypeGood),
	}
	rF.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, fData, nil, nil)

	result, err = s.deviceDiagnosis.GetState()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
}

func (s *DeviceDiagnosisSuite) Test_SetState() {
	data := &model.DeviceDiagnosisStateDataType{
		OperatingState:       util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
		PowerSupplyCondition: util.Ptr(model.PowerSupplyConditionTypeGood),
	}
	s.deviceDiagnosis.SetLocalState(data)
	assert.NotNil(s.T(), s.sentMessage)
}
