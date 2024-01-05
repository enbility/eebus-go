package spine

import (
	"time"

	"github.com/enbility/eebus-go/spine/model"
)

func createRemoteDeviceAndFeature(entityId uint, featureType model.FeatureTypeType, role model.RoleType, sender Sender) *FeatureRemoteImpl {
	localDevice := NewDeviceLocalImpl("Vendor", "DeviceName", "SerialNumber", "DeviceCode", "Address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart, time.Second*4)

	remoteDevice := NewDeviceRemoteImpl(localDevice, "ski", sender)
	// remoteDevice.address = util.Ptr(model.AddressDeviceType("Address"))
	remoteEntity := NewEntityRemoteImpl(remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(entityId)})
	remoteDevice.AddEntity(remoteEntity)
	remoteFeature := NewFeatureRemoteImpl(remoteEntity.NextFeatureId(), remoteEntity, featureType, role)
	remoteEntity.AddFeature(remoteFeature)
	return remoteFeature
}
