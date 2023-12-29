package spine_test

import (
	"time"

	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

func createRemoteDeviceAndFeature(entityId uint, featureType model.FeatureTypeType, role model.RoleType, sender spine.Sender) *spine.FeatureRemoteImpl {
	localDevice := spine.NewDeviceLocalImpl("Vendor", "DeviceName", "SerialNumber", "DeviceCode", "Address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)

	remoteDevice := spine.NewDeviceRemoteImpl(localDevice, "ski", sender)
	// remoteDevice.address = util.Ptr(model.AddressDeviceType("Address"))
	remoteEntity := spine.NewEntityRemoteImpl(remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(entityId)})
	remoteDevice.AddEntity(remoteEntity)
	remoteFeature := spine.NewFeatureRemoteImpl(remoteEntity.NextFeatureId(), remoteEntity, featureType, role)
	remoteEntity.AddFeature(remoteFeature)
	return remoteFeature
}
