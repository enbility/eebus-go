package spine_test

import (
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/mocks"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

func TestDeviceClassificationSuite(t *testing.T) {
	suite.Run(t, new(DeviceClassificationTestSuite))
}

type DeviceClassificationTestSuite struct {
	suite.Suite
	senderMock    *mocks.Sender
	function      model.FunctionType
	featureType   model.FeatureTypeType
	msgCounter    model.MsgCounterType
	remoteFeature *spine.FeatureRemoteImpl
	sut           *spine.FeatureLocalImpl
}

func (suite *DeviceClassificationTestSuite) SetupSuite() {
	suite.senderMock = mocks.NewSender(suite.T())
	suite.function = model.FunctionTypeDeviceClassificationManufacturerData
	suite.featureType = model.FeatureTypeTypeDeviceClassification
	suite.msgCounter = model.MsgCounterType(1)

	suite.remoteFeature = spine.CreateRemoteDeviceAndFeature(1, suite.featureType, model.RoleTypeServer, suite.senderMock)
	suite.sut = CreateLocalDeviceAndFeature(1, suite.featureType, model.RoleTypeClient)
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Request_Reply() {
	suite.senderMock.On("Request", model.CmdClassifierTypeRead, suite.sut.Address(), suite.remoteFeature.Address(), false, mock.AnythingOfType("[]model.CmdType")).Return(&suite.msgCounter, nil)

	// send data request
	msgCounter, err := suite.sut.RequestData(suite.function, nil, suite.remoteFeature)
	assert.Nil(suite.T(), err)

	manufacturerData := &model.DeviceClassificationManufacturerDataType{
		BrandName:    util.Ptr(model.DeviceClassificationStringType("brand name")),
		VendorName:   util.Ptr(model.DeviceClassificationStringType("vendor name")),
		DeviceName:   util.Ptr(model.DeviceClassificationStringType("device name")),
		DeviceCode:   util.Ptr(model.DeviceClassificationStringType("device code")),
		SerialNumber: util.Ptr(model.DeviceClassificationStringType("serial number")),
	}

	replyMsg := spine.Message{
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: manufacturerData,
		},
		CmdClassifier: model.CmdClassifierTypeReply,
		RequestHeader: &model.HeaderType{
			MsgCounter:          util.Ptr(model.MsgCounterType(1)),
			MsgCounterReference: &suite.msgCounter,
		},
		FeatureRemote: suite.remoteFeature,
	}
	// set response
	msgErr := suite.sut.HandleMessage(&replyMsg)
	if assert.Nil(suite.T(), msgErr) {
		remoteData := suite.remoteFeature.Data(suite.function)
		assert.IsType(suite.T(), &model.DeviceClassificationManufacturerDataType{}, remoteData, "Data has wrong type")
	}

	// Act
	result, err := suite.sut.FetchRequestData(*msgCounter, suite.remoteFeature)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.IsType(suite.T(), &model.DeviceClassificationManufacturerDataType{}, result, "Data has wrong type")
	receivedData := result.(*model.DeviceClassificationManufacturerDataType)

	assert.Equal(suite.T(), manufacturerData.BrandName, receivedData.BrandName)
	assert.Equal(suite.T(), manufacturerData.VendorName, receivedData.VendorName)
	assert.Equal(suite.T(), manufacturerData.DeviceName, receivedData.DeviceName)
	assert.Equal(suite.T(), manufacturerData.DeviceCode, receivedData.DeviceCode)
	assert.Equal(suite.T(), manufacturerData.SerialNumber, receivedData.SerialNumber)
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Request_Error() {
	suite.senderMock.On("Request", model.CmdClassifierTypeRead, suite.sut.Address(), suite.remoteFeature.Address(), false, mock.AnythingOfType("[]model.CmdType")).Return(&suite.msgCounter, nil)

	const errorNumber = model.ErrorNumberTypeGeneralError
	const errorDescription = "error occured"

	// send data request
	msgCounter, err := suite.sut.RequestData(suite.function, nil, suite.remoteFeature)
	assert.Nil(suite.T(), err)

	replyMsg := spine.Message{
		Cmd: model.CmdType{
			ResultData: &model.ResultDataType{
				ErrorNumber: util.Ptr(model.ErrorNumberType(errorNumber)),
				Description: util.Ptr(model.DescriptionType(errorDescription)),
			},
		},
		CmdClassifier: model.CmdClassifierTypeResult,
		RequestHeader: &model.HeaderType{
			MsgCounter:          util.Ptr(model.MsgCounterType(1)),
			MsgCounterReference: &suite.msgCounter,
		},
		FeatureRemote: suite.remoteFeature,
	}

	// set response
	msgErr := suite.sut.HandleMessage(&replyMsg)
	if assert.Nil(suite.T(), msgErr) {
		remoteData := suite.remoteFeature.Data(suite.function)
		assert.Nil(suite.T(), remoteData)
	}

	// Act
	result, err := suite.sut.FetchRequestData(*msgCounter, suite.remoteFeature)
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errorNumber, err.ErrorNumber)
	assert.Equal(suite.T(), errorDescription, string(*err.Description))
}

func CreateLocalDeviceAndFeature(entityId uint, featureType model.FeatureTypeType, role model.RoleType) *spine.FeatureLocalImpl {
	localDevice := spine.NewDeviceLocalImpl("Vendor", "DeviceName", "DeviceCode", "SerialNumber", "Address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := spine.NewEntityLocalImpl(localDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(entityId)})
	localDevice.AddEntity(localEntity)
	localFeature := spine.NewFeatureLocalImpl(localEntity.NextFeatureId(), localEntity, featureType, role)
	localEntity.AddFeature(localFeature)
	return localFeature
}
