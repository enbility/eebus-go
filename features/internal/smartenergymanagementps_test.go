package internal_test

import (
	"testing"

	"github.com/enbility/eebus-go/features/internal"
	shipmocks "github.com/enbility/ship-go/mocks"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestSmartEnergyManagementPsSuite(t *testing.T) {
	suite.Run(t, new(SmartEnergyManagementPsSuite))
}

type SmartEnergyManagementPsSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.SmartEnergyManagementPsCommon
}

func (s *SmartEnergyManagementPsSuite) BeforeTest(suiteName, testName string) {
	mockWriter := shipmocks.NewShipConnectionDataWriterInterface(s.T())
	mockWriter.EXPECT().WriteShipMessageWithPayload(mock.Anything).Return().Maybe()

	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		mockWriter,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeSmartEnergyManagementPs,
				functions: []model.FunctionType{
					model.FunctionTypeSmartEnergyManagementPsData,
				},
			},
		},
	)

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalSmartEnergyManagementPs(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteSmartEnergyManagementPs(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *SmartEnergyManagementPsSuite) Test_GetData() {
	value, err := s.localSut.GetData()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)
	value, err = s.remoteSut.GetData()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), value)

	s.addData()

	value, err = s.localSut.GetData()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
	value, err = s.remoteSut.GetData()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), value)
}

// helper

func (s *SmartEnergyManagementPsSuite) addData() {
	fData := &model.SmartEnergyManagementPsDataType{
		NodeScheduleInformation: &model.PowerSequenceNodeScheduleInformationDataType{},
		Alternatives:            []model.SmartEnergyManagementPsAlternativesType{},
	}
	_ = s.localFeature.UpdateData(model.FunctionTypeSmartEnergyManagementPsData, fData, nil, nil)
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, fData, nil, nil)
}
