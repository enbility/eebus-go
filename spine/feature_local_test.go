package spine

import (
	"sync"
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

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Request_Reply() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	suite.senderMock.On("Request", model.CmdClassifierTypeRead, suite.sut.Address(), suite.remoteFeature.Address(), false, mock.AnythingOfType("[]model.CmdType")).Return(&suite.msgCounter, nil)

	manufacturerData := &model.DeviceClassificationManufacturerDataType{
		BrandName:    util.Ptr(model.DeviceClassificationStringType("brand name")),
		VendorName:   util.Ptr(model.DeviceClassificationStringType("vendor name")),
		DeviceName:   util.Ptr(model.DeviceClassificationStringType("device name")),
		DeviceCode:   util.Ptr(model.DeviceClassificationStringType("device code")),
		SerialNumber: util.Ptr(model.DeviceClassificationStringType("serial number")),
	}

	go func() {
		defer wg.Done()
		// Act
		result := suite.sut.RequestData(suite.function, suite.remoteFeature)
		assert.Nil(suite.T(), result.errorResult)
		assert.NotNil(suite.T(), result.data)
		assert.IsType(suite.T(), &model.DeviceClassificationManufacturerDataType{}, result.data, "Data has wrong type")
		receivedData := result.data.(*model.DeviceClassificationManufacturerDataType)

		assert.Equal(suite.T(), manufacturerData.BrandName, receivedData.BrandName)
		assert.Equal(suite.T(), manufacturerData.VendorName, receivedData.VendorName)
		assert.Equal(suite.T(), manufacturerData.DeviceName, receivedData.DeviceName)
		assert.Equal(suite.T(), manufacturerData.DeviceCode, receivedData.DeviceCode)
		assert.Equal(suite.T(), manufacturerData.SerialNumber, receivedData.SerialNumber)
	}()

	// to make sure that the above go routine will start before HandleMessage is called
	time.Sleep(10 * time.Millisecond)

	replyMsg := Message{
		Cmd: model.CmdType{
			DeviceClassificationManufacturerData: manufacturerData,
		},
		CmdClassifier: model.CmdClassifierTypeReply,
		RequestHeader: &model.HeaderType{
			MsgCounter:          util.Ptr(model.MsgCounterType(1)),
			MsgCounterReference: &suite.msgCounter,
		},
		featureRemote: suite.remoteFeature,
	}
	msgErr := suite.sut.HandleMessage(&replyMsg)
	if assert.Nil(suite.T(), msgErr) {
		remoteData := suite.remoteFeature.Data(suite.function)
		assert.IsType(suite.T(), &model.DeviceClassificationManufacturerDataType{}, remoteData, "Data has wrong type")
	}

	wg.Wait()
}

func (suite *DeviceClassificationTestSuite) TestDeviceClassification_Request_Error() {
	wg := sync.WaitGroup{}
	wg.Add(1)

	suite.senderMock.On("Request", model.CmdClassifierTypeRead, suite.sut.Address(), suite.remoteFeature.Address(), false, mock.AnythingOfType("[]model.CmdType")).Return(&suite.msgCounter, nil)

	const errorNumber = model.ErrorNumberTypeGeneralError
	const errorDescription = "error occured"

	go func() {
		defer wg.Done()
		// Act
		result := suite.sut.RequestData(suite.function, suite.remoteFeature)
		assert.Nil(suite.T(), result.data)
		assert.NotNil(suite.T(), result.errorResult)
		assert.Equal(suite.T(), errorNumber, result.errorResult.ErrorNumber)
		assert.Equal(suite.T(), errorDescription, string(result.errorResult.Description))
	}()

	// to make sure that the above go routine will start before HandleMessage is called
	time.Sleep(10 * time.Millisecond)

	replyMsg := Message{
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
		featureRemote: suite.remoteFeature,
	}

	msgErr := suite.sut.HandleMessage(&replyMsg)
	if assert.Nil(suite.T(), msgErr) {
		remoteData := suite.remoteFeature.Data(suite.function)
		assert.Nil(suite.T(), remoteData)
	}

	wg.Wait()
}

func CreateLocalDeviceAndFeature(entityId uint, featureType model.FeatureTypeType, role model.RoleType) *FeatureLocalImpl {
	localDevice := NewDeviceLocalImpl("Vendor", "DeviceName", "DeviceCode", "SerialNumber", "Address", model.DeviceTypeTypeEnergyManagementSystem)
	localEntity := NewEntityLocalImpl(localDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(entityId)})
	localDevice.AddEntity(localEntity)
	localFeature := NewFeatureLocalImpl(localEntity.NextFeatureId(), localEntity, featureType, role)
	localEntity.AddFeature(localFeature)
	return localFeature
}
