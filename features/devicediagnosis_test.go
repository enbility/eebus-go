package features_test

import (
	"testing"
	"time"

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
					model.FunctionTypeDeviceDiagnosisHeartbeatData,
				},
			},
		},
	)

	var err error
	s.deviceDiagnosis, err = features.NewDeviceDiagnosis(s.localEntity, s.remoteEntity)
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

func (s *DeviceDiagnosisSuite) Test_IsHeartbeatWithinDuration() {
	rF := s.remoteEntity.FeatureOfAddress(util.Ptr(model.AddressFeatureType(1)))

	result := s.deviceDiagnosis.IsHeartbeatWithinDuration(time.Second * 10)
	assert.Equal(s.T(), false, result)

	now := time.Now().UTC()

	data := &model.DeviceDiagnosisHeartbeatDataType{
		HeartbeatCounter: util.Ptr(uint64(1)),
		HeartbeatTimeout: model.NewDurationType(time.Second * 4),
	}

	rF.UpdateData(model.FunctionTypeDeviceDiagnosisHeartbeatData, data, nil, nil)

	result = s.deviceDiagnosis.IsHeartbeatWithinDuration(time.Second * 10)
	assert.Equal(s.T(), false, result)

	data.Timestamp = model.NewAbsoluteOrRelativeTimeTypeFromTime(now)
	rF.UpdateData(model.FunctionTypeDeviceDiagnosisHeartbeatData, data, nil, nil)

	result = s.deviceDiagnosis.IsHeartbeatWithinDuration(time.Second * 10)
	assert.Equal(s.T(), true, result)

	time.Sleep(time.Second * 2)

	result = s.deviceDiagnosis.IsHeartbeatWithinDuration(time.Second * 1)
	assert.Equal(s.T(), false, result)
}
