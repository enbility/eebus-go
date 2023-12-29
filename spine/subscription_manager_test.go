package spine_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestSubscriptionManagerSuite(t *testing.T) {
	suite.Run(t, new(SubscriptionManagerSuite))
}

type SubscriptionManagerSuite struct {
	suite.Suite

	localDevice  *spine.DeviceLocalImpl
	remoteDevice *spine.DeviceRemoteImpl
	sut          spine.SubscriptionManager
}

func (suite *SubscriptionManagerSuite) WriteSpineMessage([]byte) {}

func (suite *SubscriptionManagerSuite) SetupSuite() {
	suite.localDevice = spine.NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)

	ski := "test"
	sender := spine.NewSender(suite)
	suite.remoteDevice = spine.NewDeviceRemoteImpl(suite.localDevice, ski, sender)

	suite.localDevice.AddRemoteDevice(ski, suite)

	suite.sut = spine.NewSubscriptionManager(suite.localDevice)
}

func (suite *SubscriptionManagerSuite) Test_Subscriptions() {
	entity := spine.NewEntityLocalImpl(suite.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	suite.localDevice.AddEntity(entity)

	localFeature := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)

	remoteEntity := spine.NewEntityRemoteImpl(suite.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})
	suite.remoteDevice.AddEntity(remoteEntity)

	remoteFeature := spine.NewFeatureRemoteImpl(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	remoteFeature.Address().Device = util.Ptr(model.AddressDeviceType("remoteDevice"))
	remoteEntity.AddFeature(remoteFeature)

	subscrRequest := model.SubscriptionManagementRequestCallType{
		ClientAddress:     remoteFeature.Address(),
		ServerAddress:     localFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
	}

	subMgr := suite.localDevice.SubscriptionManager()
	err := subMgr.AddSubscription(suite.remoteDevice, subscrRequest)
	assert.Nil(suite.T(), err)

	err = subMgr.AddSubscription(suite.remoteDevice, subscrRequest)
	assert.NotNil(suite.T(), err)

	subs := subMgr.Subscriptions(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	subscrDelete := model.SubscriptionManagementDeleteCallType{
		ClientAddress: remoteFeature.Address(),
		ServerAddress: localFeature.Address(),
	}

	err = subMgr.RemoveSubscription(subscrDelete, suite.remoteDevice)
	assert.Nil(suite.T(), err)

	subs = subMgr.Subscriptions(suite.remoteDevice)
	assert.Equal(suite.T(), 0, len(subs))

	err = subMgr.RemoveSubscription(subscrDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)
}
