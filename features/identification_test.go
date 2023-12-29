package features_test

import (
	"testing"

	"github.com/enbility/eebus-go/features"
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestIdentificationSuite(t *testing.T) {
	suite.Run(t, new(IdentificationSuite))
}

type IdentificationSuite struct {
	suite.Suite

	localEntity  *spine.EntityLocalImpl
	remoteEntity *spine.EntityRemoteImpl

	identification *features.Identification
	sentMessage    []byte
}

var _ spine.SpineDataConnection = (*IdentificationSuite)(nil)

func (s *IdentificationSuite) WriteSpineMessage(message []byte) {
	s.sentMessage = message
}

func (s *IdentificationSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeIdentification,
				functions: []model.FunctionType{
					model.FunctionTypeIdentificationListData,
				},
			},
		},
	)

	var err error
	s.identification, err = features.NewIdentification(model.RoleTypeServer, model.RoleTypeClient, s.localEntity, s.remoteEntity)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.identification)
}

func (s *IdentificationSuite) Test_RequestValues() {
	counter, err := s.identification.RequestValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), counter)
}

func (s *IdentificationSuite) Test_GetValues() {
	data, err := s.identification.GetValues()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	s.addData()

	data, err = s.identification.GetValues()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *IdentificationSuite) addData() {
	rF := s.remoteEntity.Feature(util.Ptr(model.AddressFeatureType(1)))
	fData := &model.IdentificationListDataType{
		IdentificationData: []model.IdentificationDataType{
			{
				IdentificationId:    util.Ptr(model.IdentificationIdType(0)),
				IdentificationType:  util.Ptr(model.IdentificationTypeTypeEui64),
				IdentificationValue: util.Ptr(model.IdentificationValueType("test")),
			},
		},
	}
	rF.UpdateData(model.FunctionTypeIdentificationListData, fData, nil, nil)
}
