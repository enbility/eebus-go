package spine

import (
	"testing"
	"time"

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
	remoteFeature *FeatureRemoteImpl
	sut           *FeatureLocalImpl
}

func (suite *DeviceClassificationTestSuite) SetupSuite() {
	suite.senderMock = mocks.NewSender(suite.T())
	suite.function = model.FunctionTypeDeviceClassificationManufacturerData
	suite.featureType = model.FeatureTypeTypeDeviceClassification
	suite.msgCounter = model.MsgCounterType(1)

	suite.remoteFeature = CreateRemoteDeviceAndFeature(1, suite.featureType, model.RoleTypeServer, suite.senderMock)
	suite.sut = CreateLocalDeviceAndFeature(1, suite.featureType, model.RoleTypeClient)
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Request() {
	suite.senderMock.On("Request", model.CmdClassifierTypeRead, suite.sut.Address(), suite.remoteFeature.Address(), false, mock.AnythingOfType("[]model.CmdType")).Return(&suite.msgCounter, nil)

	// Act
	usedMsgCounter, err := suite.sut.RequestData(suite.function, suite.remoteFeature, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), &suite.msgCounter, usedMsgCounter)
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Request_Reply() {
	suite.senderMock.On("Request", model.CmdClassifierTypeRead, suite.sut.Address(), suite.remoteFeature.Address(), false, mock.AnythingOfType("[]model.CmdType")).Return(&suite.msgCounter, nil)

	manufacturerData := &model.DeviceClassificationManufacturerDataType{
		BrandName:    util.Ptr(model.DeviceClassificationStringType("brand name")),
		VendorName:   util.Ptr(model.DeviceClassificationStringType("vendor name")),
		DeviceName:   util.Ptr(model.DeviceClassificationStringType("device name")),
		DeviceCode:   util.Ptr(model.DeviceClassificationStringType("device code")),
		SerialNumber: util.Ptr(model.DeviceClassificationStringType("serial number")),
	}

	requestChannel := make(chan *model.DeviceClassificationManufacturerDataType)
	_, err := suite.sut.RequestData(suite.function, suite.remoteFeature, requestChannel)
	assert.NoError(suite.T(), err)

	replyMsg := Message{
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: manufacturerData,
		},
		CmdClassifier: model.CmdClassifierTypeReply,
		RequestHeader: &model.HeaderType{
			MsgCounter: &suite.msgCounter,
		},
		featureRemote: suite.remoteFeature,
	}

	// Act
	go func() {
		msgErr := suite.sut.HandleMessage(&replyMsg)
		if assert.Nil(suite.T(), msgErr) {
			remoteData := suite.remoteFeature.Data(suite.function)
			assert.IsType(suite.T(), &model.DeviceClassificationManufacturerDataType{}, remoteData, "Data has wrong type")
		}
	}()

	channelData := util.ReceiveWithTimeout(requestChannel, time.Duration(time.Second*2))
	assert.NotNil(suite.T(), channelData)
	if channelData != nil {
		assert.Equal(suite.T(), manufacturerData.BrandName, channelData.BrandName)
		assert.Equal(suite.T(), manufacturerData.VendorName, channelData.VendorName)
		assert.Equal(suite.T(), manufacturerData.DeviceName, channelData.DeviceName)
		assert.Equal(suite.T(), manufacturerData.DeviceCode, channelData.DeviceCode)
		assert.Equal(suite.T(), manufacturerData.SerialNumber, channelData.SerialNumber)
	}
}

func CreateLocalDeviceAndFeature(entityId uint, featureType model.FeatureTypeType, role model.RoleType) *FeatureLocalImpl {
	localDevice := NewDeviceLocalImpl("Vendor", "DeviceName", "DeviceCode", "SerialNumber", "Address", model.DeviceTypeTypeEnergyManagementSystem)
	localEntity := NewEntityLocalImpl(localDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(entityId)})
	localDevice.AddEntity(localEntity)
	localFeature := NewFeatureLocalImpl(localEntity.NextFeatureId(), localEntity, featureType, role)
	localEntity.AddFeature(localFeature)
	return localFeature
}
