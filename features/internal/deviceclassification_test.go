package internal_test

import (
	"testing"

	"github.com/enbility/eebus-go/features/internal"
	shipapi "github.com/enbility/ship-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDeviceClassificationSuite(t *testing.T) {
	suite.Run(t, new(DeviceClassificationSuite))
}

type DeviceClassificationSuite struct {
	suite.Suite

	localEntity  spineapi.EntityLocalInterface
	remoteEntity spineapi.EntityRemoteInterface

	localFeature  spineapi.FeatureLocalInterface
	remoteFeature spineapi.FeatureRemoteInterface

	localSut,
	remoteSut *internal.DeviceClassificationCommon
	sentMessage []byte
}

var _ shipapi.ShipConnectionDataWriterInterface = (*DeviceClassificationSuite)(nil)

func (s *DeviceClassificationSuite) WriteShipMessageWithPayload(message []byte) {
	s.sentMessage = message
}

func (s *DeviceClassificationSuite) BeforeTest(suiteName, testName string) {
	s.localEntity, s.remoteEntity = setupFeatures(
		s.T(),
		s,
		[]featureFunctions{
			{
				featureType: model.FeatureTypeTypeDeviceClassification,
				functions: []model.FunctionType{
					model.FunctionTypeDeviceClassificationManufacturerData,
				},
			},
		},
	)

	s.localFeature = s.localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	assert.NotNil(s.T(), s.localFeature)
	s.localSut = internal.NewLocalDeviceClassification(s.localFeature)
	assert.NotNil(s.T(), s.localSut)

	s.remoteFeature = s.remoteEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
	assert.NotNil(s.T(), s.remoteFeature)
	s.remoteSut = internal.NewRemoteDeviceClassification(s.remoteFeature)
	assert.NotNil(s.T(), s.remoteSut)
}

func (s *DeviceClassificationSuite) Test_GetManufacturerDetails() {
	result, err := s.remoteSut.GetManufacturerDetails()
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), result)

	fData := &model.DeviceClassificationManufacturerDataType{
		DeviceName:                     util.Ptr(model.DeviceClassificationStringType("brand")),
		DeviceCode:                     util.Ptr(model.DeviceClassificationStringType("brand")),
		SerialNumber:                   util.Ptr(model.DeviceClassificationStringType("brand")),
		SoftwareRevision:               util.Ptr(model.DeviceClassificationStringType("brand")),
		HardwareRevision:               util.Ptr(model.DeviceClassificationStringType("brand")),
		VendorName:                     util.Ptr(model.DeviceClassificationStringType("brand")),
		VendorCode:                     util.Ptr(model.DeviceClassificationStringType("brand")),
		BrandName:                      util.Ptr(model.DeviceClassificationStringType("brand")),
		PowerSource:                    util.Ptr(model.PowerSourceType("brand")),
		ManufacturerNodeIdentification: util.Ptr(model.DeviceClassificationStringType("brand")),
		ManufacturerLabel:              util.Ptr(model.LabelType("label")),
		ManufacturerDescription:        util.Ptr(model.DescriptionType("description")),
	}
	_, _ = s.remoteFeature.UpdateData(true, model.FunctionTypeDeviceClassificationManufacturerData, fData, nil, nil)
	_ = s.localFeature.UpdateData(model.FunctionTypeDeviceClassificationManufacturerData, fData, nil, nil)

	result, err = s.remoteSut.GetManufacturerDetails()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)

	result, err = s.localSut.GetManufacturerDetails()
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), result)
}
