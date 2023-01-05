package features

import (
	"testing"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDeviceDiagnosisSuite(t *testing.T) {
	suite.Run(t, new(DeviceDiagnosisSuite))
}

type DeviceDiagnosisSuite struct {
	suite.Suite

	localDevice  *spine.DeviceLocalImpl
	remoteEntity *spine.EntityRemoteImpl

	deviceDiagnosis *DeviceDiagnosis
	sentMessage     []byte
}

var _ spine.SpineDataConnection = (*DeviceDiagnosisSuite)(nil)

func (s *DeviceDiagnosisSuite) WriteSpineMessage(message []byte) {
	s.sentMessage = message
}

func (s *DeviceDiagnosisSuite) BeforeTest(suiteName, testName string) {
	s.localDevice, s.remoteEntity = setupFeatures(
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
	s.deviceDiagnosis, err = NewDeviceDiagnosis(model.RoleTypeServer, model.RoleTypeClient, s.localDevice, s.remoteEntity)
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

	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.DeviceDiagnosisStateDataType{
		OperatingState:       util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
		PowerSupplyCondition: util.Ptr(model.PowerSupplyConditionTypeGood),
	}
	rF.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, fData, nil, nil)

	result, err = s.deviceDiagnosis.GetState()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
}

func (s *DeviceDiagnosisSuite) Test_SendState() {
	data := &model.DeviceDiagnosisStateDataType{
		OperatingState:       util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
		PowerSupplyCondition: util.Ptr(model.PowerSupplyConditionTypeGood),
	}
	s.deviceDiagnosis.SendState(data)
	assert.NotNil(s.T(), s.sentMessage)
}
