package client

import (
	"testing"

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

	deviceDiagnosis *DeviceDiagnosis
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
					model.FunctionTypeDeviceDiagnosisHeartbeatData,
				},
			},
		},
	)

	var err error
	s.deviceDiagnosis, err = NewDeviceDiagnosis(s.localEntity, nil)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), s.deviceDiagnosis)

	s.deviceDiagnosis, err = NewDeviceDiagnosis(s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.deviceDiagnosis)
}

func (s *DeviceDiagnosisSuite) Test_RequestHeartbeat() {
	counter, err := s.deviceDiagnosis.RequestHeartbeat()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *DeviceDiagnosisSuite) Test_RequestState() {
	counter, err := s.deviceDiagnosis.RequestState()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}
