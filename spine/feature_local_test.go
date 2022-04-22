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
	suite.senderMock = new(mocks.Sender)
	suite.function = model.FunctionTypeDeviceClassificationManufacturerData
	suite.featureType = model.FeatureTypeTypeDeviceClassification
	suite.msgCounter = model.MsgCounterType(1)

	suite.remoteFeature = createRemoteFeature(suite.featureType, model.RoleTypeServer, suite.senderMock)
	suite.sut = createLocalFeature(suite.featureType, model.RoleTypeClient)
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Request() {
	suite.senderMock.On("Request", model.CmdClassifierTypeRead, suite.sut.Address(), false, mock.AnythingOfType("[]model.CmdType")).Return(&suite.msgCounter, nil)

	// Act
	usedMsgCounter, err := suite.sut.RequestData(suite.function, suite.remoteFeature, nil)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), &suite.msgCounter, usedMsgCounter)
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Request_Reply() {
	suite.senderMock.On("Request", model.CmdClassifierTypeRead, suite.sut.Address(), false, mock.AnythingOfType("[]model.CmdType")).Return(&suite.msgCounter, nil)

	manufacturerData := &model.DeviceClassificationManufacturerDataType{
		BrandName:  util.Ptr(model.DeviceClassificationStringType("brand name")),
		DeviceName: util.Ptr(model.DeviceClassificationStringType("device name")),
		DeviceCode: util.Ptr(model.DeviceClassificationStringType("device code")),
	}

	requestChannel := make(chan *model.DeviceClassificationManufacturerDataType)
	suite.sut.RequestData(suite.function, suite.remoteFeature, requestChannel)

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
		err := suite.sut.HandleMessage(&replyMsg)
		if assert.NoError(suite.T(), err) {
			remoteData := suite.remoteFeature.Data(suite.function)
			assert.IsType(suite.T(), &model.DeviceClassificationManufacturerDataType{}, remoteData, "Data has wrong type")
			remoteManufacturerData := remoteData.(*model.DeviceClassificationManufacturerDataType)
			assert.Equal(suite.T(), manufacturerData.BrandName, remoteManufacturerData.BrandName)
			assert.Equal(suite.T(), manufacturerData.DeviceName, remoteManufacturerData.DeviceName)
			assert.Equal(suite.T(), manufacturerData.DeviceCode, remoteManufacturerData.DeviceCode)
		}
	}()

	channelData := util.ReceiveWithTimeout(requestChannel, time.Duration(time.Second*2))
	assert.NotNil(suite.T(), channelData)
	assert.Equal(suite.T(), manufacturerData.BrandName, channelData.BrandName)
	assert.Equal(suite.T(), manufacturerData.DeviceName, channelData.DeviceName)
	assert.Equal(suite.T(), manufacturerData.DeviceCode, channelData.DeviceCode)
}

func createLocalFeature(featureType model.FeatureTypeType, role model.RoleType) *FeatureLocalImpl {
	localDevice := NewDeviceLocalImpl(model.AddressDeviceType("localDevice"))
	localEntity := NewEntityLocalImpl(localDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(0)})
	return NewFeatureLocalImpl(0, localEntity, featureType, role)
}

func createRemoteFeature(featureType model.FeatureTypeType, role model.RoleType, sender Sender) *FeatureRemoteImpl {
	remoteDevice := NewDeviceRemoteImpl(model.AddressDeviceType("remoteDevice"))
	remoteEntity := NewEntityRemoteImpl(remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(0)})
	return NewFeatureRemoteImpl(0, remoteEntity, featureType, role, sender)
}
