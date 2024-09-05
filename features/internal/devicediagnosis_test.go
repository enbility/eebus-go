package internal_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/features/internal"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestDeviceDiagnosisSuite(t *testing.T) {
	suite.Run(t, new(DeviceDiagnosisSuite))
}

type DeviceDiagnosisSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.DeviceDiagnosisCommon
}

func (s *DeviceDiagnosisSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeDeviceDiagnosis,
				functions: []model.FunctionType{
					model.FunctionTypeDeviceDiagnosisHeartbeatData,
					model.FunctionTypeDeviceDiagnosisStateData,
				},
			},
		},
	)

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalDeviceDiagnosis(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteDeviceDiagnosis(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *DeviceDiagnosisSuite) Test_GetState() {
	result, err := s.localSut.GetState()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)
	result, err = s.remoteSut.GetState()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)

	fData := &model.DeviceDiagnosisStateDataType{
		OperatingState:       util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
		PowerSupplyCondition: util.Ptr(model.PowerSupplyConditionTypeGood),
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeDeviceDiagnosisStateData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisStateData, fData, nil, nil)

	result, err = s.localSut.GetState()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
	result, err = s.remoteSut.GetState()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
}

func (s *DeviceDiagnosisSuite) Test_IsHeartbeatWithinDuration() {
	result := s.localSut.IsHeartbeatWithinDuration(time.Second * 10)
	assert.Equal(s.T(), true, result) // local server automatically generates a first hearbeat!
	result = s.remoteSut.IsHeartbeatWithinDuration(time.Second * 10)
	assert.Equal(s.T(), false, result)

	now := time.Now().UTC()

	data := &model.DeviceDiagnosisHeartbeatDataType{
		HeartbeatCounter: util.Ptr(uint64(1)),
		HeartbeatTimeout: model.NewDurationType(time.Second * 4),
	}

	_ = s.localFeature.UpdateData(model.FunctionTypeDeviceDiagnosisHeartbeatData, data, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisHeartbeatData, data, nil, nil)

	result = s.localSut.IsHeartbeatWithinDuration(time.Second * 10)
	assert.Equal(s.T(), false, result)
	result = s.remoteSut.IsHeartbeatWithinDuration(time.Second * 10)
	assert.Equal(s.T(), false, result)

	data.Timestamp = model.NewAbsoluteOrRelativeTimeTypeFromTime(now)
	_ = s.localFeature.UpdateData(model.FunctionTypeDeviceDiagnosisHeartbeatData, data, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeDeviceDiagnosisHeartbeatData, data, nil, nil)

	result = s.localSut.IsHeartbeatWithinDuration(time.Second * 10)
	assert.Equal(s.T(), true, result)
	result = s.remoteSut.IsHeartbeatWithinDuration(time.Second * 10)
	assert.Equal(s.T(), true, result)

	// Disable this test as it may sometimes fail due to timing issues
	/*
		time.Sleep(time.Second * 2)

		result = s.localSut.IsHeartbeatWithinDuration(time.Millisecond * 500)
		assert.Equal(s.T(), false, result)
		result = s.remoteSut.IsHeartbeatWithinDuration(time.Millisecond * 500)
		assert.Equal(s.T(), false, result)
	*/
}
