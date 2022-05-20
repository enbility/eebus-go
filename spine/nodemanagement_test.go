package spine

import (
	"reflect"
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine/mocks"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodemanagement_SubscriptionRequestCall(t *testing.T) {

	// const serverName = "Server"
	// const clientName = "Client"

	const subscriptionEntityId uint = 1
	//const subscriptionFeatureId uint = 2
	const featureType = model.FeatureTypeTypeDeviceClassification

	senderMock := mocks.NewSender(t)

	//localDevice := NewDeviceLocalImpl(model.AddressDeviceType("server"))

	serverFeature := CreateLocalDeviceAndFeature(subscriptionEntityId, featureType, model.RoleTypeServer)
	clientFeature := CreateRemoteDeviceAndFeature(subscriptionEntityId, featureType, model.RoleTypeClient, senderMock)

	// serverAddress := featureAddress(serverName, subscriptionEntityId, subscriptionFeatureId)
	// clientAddress := featureAddress(clientName, subscriptionEntityId, subscriptionFeatureId)
	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 1, len(cmd.NodeManagementSubscriptionData.SubscriptionEntry))
		assert.True(t, reflect.DeepEqual(cmd.NodeManagementSubscriptionData.SubscriptionEntry[0].ClientAddress, clientFeature.Address()))
		assert.True(t, reflect.DeepEqual(cmd.NodeManagementSubscriptionData.SubscriptionEntry[0].ServerAddress, serverFeature.Address()))

	}).Return(nil).Once()

	// serverFeatureMock := newFeatureLocalMock(serverAddress, model.RoleTypeServer, featureType, senderMock)
	// clientFeatureMock := newFeatureRemoteMock(clientAddress, model.RoleTypeClient, featureType)

	requestMsg := Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionRequestCall: NewNodeManagementSubscriptionRequestCallType(
				clientFeature.Address(), serverFeature.Address(), featureType),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		featureRemote: clientFeature,
	}

	sut := NewNodeManagementImpl(0, serverFeature.Entity())

	// Act
	err := sut.HandleMessage(&requestMsg)
	if assert.Nil(t, err) {

		dataMsg := Message{
			Cmd: model.CmdType{
				NodeManagementSubscriptionData: &model.NodeManagementSubscriptionDataType{},
			},
			CmdClassifier: model.CmdClassifierTypeCall,
			featureRemote: clientFeature,
		}
		err = sut.HandleMessage(&dataMsg)
		assert.Nil(t, err)
	}
}

// func newFeatureLocalMock(address *model.FeatureAddressType, role model.RoleType, ftype model.FeatureTypeType, sender spine.Sender) *mocks.FeatureLocal {
// 	deviceMock := new(mocks.DeviceLocal)
// 	entityMock := new(mocks.EntityLocal)
// 	featureMock := new(mocks.FeatureLocal)

// 	deviceMock.On("FeatureByAddress", address).Return(featureMock)
// 	deviceMock.On("Address").Return(address.Device)
// 	deviceMock.On("Sender").Return(sender)

// 	entityMock.On("Device").Return(deviceMock)
// 	entityMock.On("Address").Return(address.Entity)

// 	featureMock.On("Role").Return(role)
// 	featureMock.On("Type").Return(model.FeatureTypeEnumType(ftype))
// 	featureMock.On("Device").Return(deviceMock)
// 	featureMock.On("Entity").Return(entityMock)

// 	return featureMock
// }

// func newFeatureRemoteMock(address *model.FeatureAddressType, role model.RoleType, ftype model.FeatureTypeType) *mocks.FeatureRemote {
// 	deviceMock := new(mocks.DeviceRemote)
// 	entityMock := new(mocks.EntityRemote)
// 	featureMock := new(mocks.FeatureRemote)

// 	deviceMock.On("FeatureByAddress", address).Return(featureMock)
// 	deviceMock.On("Address").Return(address.Device)
// 	//deviceMock.On("Sender").Return(sender)

// 	entityMock.On("Device").Return(deviceMock)
// 	entityMock.On("Address").Return(address.Entity)

// 	featureMock.On("Role").Return(role)
// 	featureMock.On("Type").Return(model.FeatureTypeEnumType(ftype))
// 	featureMock.On("Device").Return(deviceMock)
// 	featureMock.On("Entity").Return(entityMock)

// 	return featureMock
// }
