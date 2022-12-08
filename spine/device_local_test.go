package spine_test

import (
	"testing"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestDeviceLocalSuite(t *testing.T) {
	suite.Run(t, new(DeviceLocalTestSuite))
}

type DeviceLocalTestSuite struct {
	suite.Suite
}

var _ spine.SpineDataConnection = (*DeviceLocalTestSuite)(nil)

func (d *DeviceLocalTestSuite) WriteSpineMessage([]byte) {}

func (d *DeviceLocalTestSuite) Test_RemoteDevice() {
	sut := spine.NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)
	localEntity := spine.NewEntityLocalImpl(sut, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	sut.AddEntity(localEntity)

	f := spine.NewFeatureLocalImpl(1, localEntity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	localEntity.AddFeature(f)
	f = spine.NewFeatureLocalImpl(2, localEntity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	localEntity.AddFeature(f)

	ski := "test"
	remote := sut.RemoteDeviceForSki(ski)
	assert.Nil(d.T(), remote)

	devices := sut.RemoteDevices()
	assert.Equal(d.T(), 0, len(devices))

	sut.AddRemoteDevice(ski, d)
	remote = sut.RemoteDeviceForSki(ski)
	assert.NotNil(d.T(), remote)

	devices = sut.RemoteDevices()
	assert.Equal(d.T(), 1, len(devices))

	entities := sut.Entities()
	assert.Equal(d.T(), 2, len(entities))

	entity1 := sut.Entity([]model.AddressEntityType{1})
	assert.NotNil(d.T(), entity1)

	entity2 := sut.Entity([]model.AddressEntityType{2})
	assert.Nil(d.T(), entity2)

	featureAddress := &model.FeatureAddressType{
		Entity:  []model.AddressEntityType{1},
		Feature: util.Ptr(model.AddressFeatureType(1)),
	}
	feature1 := sut.FeatureByAddress(featureAddress)
	assert.NotNil(d.T(), feature1)

	feature2 := sut.FeatureByTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	assert.NotNil(d.T(), feature2)

	sut.RemoveEntity(entity1)
	entities = sut.Entities()
	assert.Equal(d.T(), 1, len(entities))

	sut.RemoveRemoteDevice(ski)
	remote = sut.RemoteDeviceForSki(ski)
	assert.Nil(d.T(), remote)
}
