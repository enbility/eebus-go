package spine

import (
	"testing"
	"time"

	"github.com/DerAndereAndi/eebus-go/spine/mocks"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestDeviceClassification_Request(t *testing.T) {
	featureType := model.FeatureTypeEnumTypeDeviceClassification
	function := model.FunctionEnumTypeDeviceClassificationManufacturerData

	senderMock := new(mocks.Sender)
	msgCounter := model.MsgCounterType(1)
	clientAddress := model.FeatureAddressType{}
	serverAddress := model.FeatureAddressType{}
	senderMock.On("Request", model.CmdClassifierTypeRead, &clientAddress, &serverAddress, false, mock.AnythingOfType("[]model.CmdType")).Return(&msgCounter, nil)

	sut := NewFeatureLocalImpl(&clientAddress, featureType, senderMock)

	// Act
	usedMsgCounter, err := sut.RequestData(function, &serverAddress, nil)
	assert.NoError(t, err)
	assert.Equal(t, &msgCounter, usedMsgCounter)
}

func TestDeviceClassification_Request_Reply(t *testing.T) {
	featureType := model.FeatureTypeEnumTypeDeviceClassification
	function := model.FunctionEnumTypeDeviceClassificationManufacturerData

	senderMock := new(mocks.Sender)
	msgCounter := model.MsgCounterType(1)
	clientAddress := model.FeatureAddressType{}
	serverAddress := model.FeatureAddressType{}
	senderMock.On("Request", model.CmdClassifierTypeRead, &clientAddress, &serverAddress, false, mock.AnythingOfType("[]model.CmdType")).Return(&msgCounter, nil)

	manufacturerData := &model.DeviceClassificationManufacturerDataType{
		BrandName:  util.Ptr(model.DeviceClassificationStringType("brand name")),
		DeviceName: util.Ptr(model.DeviceClassificationStringType("device name")),
		DeviceCode: util.Ptr(model.DeviceClassificationStringType("device code")),
	}

	remoteFeature := NewFeatureRemoteImpl(featureType)

	sut := NewFeatureLocalImpl(&clientAddress, featureType, senderMock)

	requestChannel := make(chan *model.DeviceClassificationManufacturerDataType)
	sut.RequestData(function, &serverAddress, requestChannel)

	replyMsg := Message{
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: manufacturerData,
		},
		CmdClassifier: model.CmdClassifierTypeReply,
		RequestHeader: &model.HeaderType{
			MsgCounter: &msgCounter,
		},
		featureRemote: remoteFeature,
	}

	// Act
	go func() {
		err := sut.HandleMessage(&replyMsg)
		if assert.NoError(t, err) {
			remoteData := remoteFeature.Data(function)
			assert.IsType(t, &model.DeviceClassificationManufacturerDataType{}, remoteData, "Data has wrong type")
			remoteManufacturerData := remoteData.(*model.DeviceClassificationManufacturerDataType)
			assert.Equal(t, manufacturerData.BrandName, remoteManufacturerData.BrandName)
			assert.Equal(t, manufacturerData.DeviceName, remoteManufacturerData.DeviceName)
			assert.Equal(t, manufacturerData.DeviceCode, remoteManufacturerData.DeviceCode)
		}
	}()

	channelData := util.ReceiveWithTimeout(requestChannel, time.Duration(time.Second*2))
	assert.NotNil(t, channelData)
	assert.Equal(t, manufacturerData.BrandName, channelData.BrandName)
	assert.Equal(t, manufacturerData.DeviceName, channelData.DeviceName)
	assert.Equal(t, manufacturerData.DeviceCode, channelData.DeviceCode)
}
