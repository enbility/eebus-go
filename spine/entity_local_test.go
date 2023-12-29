package spine_test

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestEntityLocalSuite(t *testing.T) {
	suite.Run(t, new(EntityLocalTestSuite))
}

type EntityLocalTestSuite struct {
	suite.Suite
}

func (suite *EntityLocalTestSuite) Test_Entity() {
	device := spine.NewDeviceLocalImpl("brand", "model", "serial", "code", "address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)
	entity := spine.NewEntityLocalImpl(device, model.EntityTypeTypeCEM, spine.NewAddressEntityType([]uint{1}))
	device.AddEntity(entity)

	f := spine.NewFeatureLocalImpl(1, entity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
	entity.AddFeature(f)
	assert.Equal(suite.T(), 1, len(entity.Features()))

	entity.AddFeature(f)
	assert.Equal(suite.T(), 1, len(entity.Features()))

	f1 := entity.Feature(nil)
	assert.Nil(suite.T(), f1)

	f1 = entity.Feature(f.Address().Feature)
	assert.NotNil(suite.T(), f1)

	fakeAddress := model.AddressFeatureType(5)
	f1 = entity.Feature(&fakeAddress)
	assert.Nil(suite.T(), f1)

	f2 := entity.GetOrAddFeature(model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
	assert.NotNil(suite.T(), f2)

	assert.Equal(suite.T(), 2, len(entity.Features()))

	f3 := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(suite.T(), f3)

	assert.Equal(suite.T(), 3, len(entity.Features()))

	f4 := entity.GetOrAddFeature(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(suite.T(), f4)

	assert.Equal(suite.T(), 3, len(entity.Features()))
}
