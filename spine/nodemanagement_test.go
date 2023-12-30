package spine_test

import (
	"reflect"
	"testing"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/mocks"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestNodemanagement_BindingCalls(t *testing.T) {
	const bindingEntityId uint = 1
	const featureType = model.FeatureTypeTypeLoadControl

	senderMock := mocks.NewSender(t)

	serverFeature := createLocalDeviceAndFeature(bindingEntityId, featureType, model.RoleTypeServer)
	clientFeature := createRemoteDeviceAndFeature(bindingEntityId, featureType, model.RoleTypeClient, senderMock)

	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 1, len(cmd.NodeManagementBindingData.BindingEntry))
		assert.True(t, reflect.DeepEqual(cmd.NodeManagementBindingData.BindingEntry[0].ClientAddress, clientFeature.Address()))
		assert.True(t, reflect.DeepEqual(cmd.NodeManagementBindingData.BindingEntry[0].ServerAddress, serverFeature.Address()))
	}).Return(nil).Once()

	requestMsg := spine.Message{
		Cmd: model.CmdType{
			NodeManagementBindingRequestCall: spine.NewNodeManagementBindingRequestCallType(
				clientFeature.Address(), serverFeature.Address(), featureType),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		FeatureRemote: clientFeature,
	}

	sut := spine.NewNodeManagementImpl(0, serverFeature.Entity())

	// Act
	err := sut.HandleMessage(&requestMsg)
	if assert.Nil(t, err) {

		dataMsg := spine.Message{
			Cmd: model.CmdType{
				NodeManagementBindingData: &model.NodeManagementBindingDataType{},
			},
			CmdClassifier: model.CmdClassifierTypeCall,
			FeatureRemote: clientFeature,
		}
		err = sut.HandleMessage(&dataMsg)
		assert.Nil(t, err)
	}

	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 0, len(cmd.NodeManagementBindingData.BindingEntry))
	}).Return(nil).Once()

	deleteMsg := spine.Message{
		Cmd: model.CmdType{
			NodeManagementBindingDeleteCall: spine.NewNodeManagementBindingDeleteCallType(
				clientFeature.Address(), serverFeature.Address()),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		FeatureRemote: clientFeature,
	}

	// Act
	err = sut.HandleMessage(&deleteMsg)
	if assert.Nil(t, err) {

		dataMsg := spine.Message{
			Cmd: model.CmdType{
				NodeManagementBindingData: &model.NodeManagementBindingDataType{},
			},
			CmdClassifier: model.CmdClassifierTypeCall,
			FeatureRemote: clientFeature,
		}
		err = sut.HandleMessage(&dataMsg)
		assert.Nil(t, err)
	}
}

func TestNodemanagement_SubscriptionCalls(t *testing.T) {
	const subscriptionEntityId uint = 1
	const featureType = model.FeatureTypeTypeDeviceClassification

	senderMock := mocks.NewSender(t)

	serverFeature := createLocalDeviceAndFeature(subscriptionEntityId, featureType, model.RoleTypeServer)
	clientFeature := createRemoteDeviceAndFeature(subscriptionEntityId, featureType, model.RoleTypeClient, senderMock)

	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 1, len(cmd.NodeManagementSubscriptionData.SubscriptionEntry))
		assert.True(t, reflect.DeepEqual(cmd.NodeManagementSubscriptionData.SubscriptionEntry[0].ClientAddress, clientFeature.Address()))
		assert.True(t, reflect.DeepEqual(cmd.NodeManagementSubscriptionData.SubscriptionEntry[0].ServerAddress, serverFeature.Address()))
	}).Return(nil).Once()

	requestMsg := spine.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionRequestCall: spine.NewNodeManagementSubscriptionRequestCallType(
				clientFeature.Address(), serverFeature.Address(), featureType),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		FeatureRemote: clientFeature,
	}

	sut := spine.NewNodeManagementImpl(0, serverFeature.Entity())

	// Act
	err := sut.HandleMessage(&requestMsg)
	if assert.Nil(t, err) {

		dataMsg := spine.Message{
			Cmd: model.CmdType{
				NodeManagementSubscriptionData: &model.NodeManagementSubscriptionDataType{},
			},
			CmdClassifier: model.CmdClassifierTypeCall,
			FeatureRemote: clientFeature,
		}
		err = sut.HandleMessage(&dataMsg)
		assert.Nil(t, err)
	}

	senderMock.On("Reply", mock.Anything, mock.Anything, mock.Anything).Run(func(args mock.Arguments) {
		cmd := args.Get(2).(model.CmdType)
		assert.Equal(t, 0, len(cmd.NodeManagementSubscriptionData.SubscriptionEntry))
	}).Return(nil).Once()

	deleteMsg := spine.Message{
		Cmd: model.CmdType{
			NodeManagementSubscriptionDeleteCall: spine.NewNodeManagementSubscriptionDeleteCallType(
				clientFeature.Address(), serverFeature.Address(), featureType),
		},
		CmdClassifier: model.CmdClassifierTypeCall,
		FeatureRemote: clientFeature,
	}

	// Act
	err = sut.HandleMessage(&deleteMsg)
	if assert.Nil(t, err) {

		dataMsg := spine.Message{
			Cmd: model.CmdType{
				NodeManagementSubscriptionData: &model.NodeManagementSubscriptionDataType{},
			},
			CmdClassifier: model.CmdClassifierTypeCall,
			FeatureRemote: clientFeature,
		}
		err = sut.HandleMessage(&dataMsg)
		assert.Nil(t, err)
	}
}
