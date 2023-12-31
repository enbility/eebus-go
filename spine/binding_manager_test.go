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

func TestBindingManagerSuite(t *testing.T) {
	suite.Run(t, new(BindingManagerSuite))
}

type BindingManagerSuite struct {
	suite.Suite

	localDevice  *spine.DeviceLocalImpl
	remoteDevice *spine.DeviceRemoteImpl
	sut          spine.BindingManager
}

func (suite *BindingManagerSuite) WriteSpineMessage([]byte) {}

func (suite *BindingManagerSuite) SetupSuite() {
	suite.localDevice = spine.NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)

	ski := "test"
	sender := spine.NewSender(suite)
	suite.remoteDevice = spine.NewDeviceRemoteImpl(suite.localDevice, ski, sender)

	suite.localDevice.AddRemoteDevice(ski, suite)

	suite.sut = spine.NewBindingManager(suite.localDevice)
}

func (suite *BindingManagerSuite) Test_Bindings() {
	entity := spine.NewEntityLocalImpl(suite.localDevice, model.EntityTypeTypeCEM, []model.AddressEntityType{1})
	suite.localDevice.AddEntity(entity)

	localFeature := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)

	remoteEntity := spine.NewEntityRemoteImpl(suite.remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{1})
	suite.remoteDevice.AddEntity(remoteEntity)

	remoteFeature := spine.NewFeatureRemoteImpl(remoteEntity.NextFeatureId(), remoteEntity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
	remoteFeature.Address().Device = util.Ptr(model.AddressDeviceType("remoteDevice"))
	remoteEntity.AddFeature(remoteFeature)

	bindingRequest := model.BindingManagementRequestCallType{
		ClientAddress:     remoteFeature.Address(),
		ServerAddress:     localFeature.Address(),
		ServerFeatureType: util.Ptr(model.FeatureTypeTypeDeviceDiagnosis),
	}

	bindingMgr := suite.localDevice.BindingManager()
	err := bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.Nil(suite.T(), err)

	subs := bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.NotNil(suite.T(), err)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	bindingDelete := model.BindingManagementDeleteCallType{
		ClientAddress: remoteFeature.Address(),
		ServerAddress: localFeature.Address(),
	}

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.Nil(suite.T(), err)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 0, len(subs))

	err = bindingMgr.RemoveBinding(bindingDelete, suite.remoteDevice)
	assert.NotNil(suite.T(), err)

	err = bindingMgr.AddBinding(suite.remoteDevice, bindingRequest)
	assert.Nil(suite.T(), err)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 1, len(subs))

	bindingMgr.RemoveBindingsForDevice(suite.remoteDevice)

	subs = bindingMgr.Bindings(suite.remoteDevice)
	assert.Equal(suite.T(), 0, len(subs))
}
