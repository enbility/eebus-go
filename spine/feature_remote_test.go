package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

func CreateRemoteDeviceAndFeature(entityId uint, featureType model.FeatureTypeType, role model.RoleType, sender Sender) *FeatureRemoteImpl {
	localDevice := NewDeviceLocalImpl("Vendor", "DeviceName", "DeviceCode", "SerialNumber", model.DeviceTypeTypeEnergyManagementSystem)

	remoteDevice := NewDeviceRemoteImpl(localDevice, "ski", nil, nil)
	remoteEntity := NewEntityRemoteImpl(remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(entityId)})
	remoteDevice.addEntity(remoteEntity)
	remoteFeature := NewFeatureRemoteImpl(remoteEntity.NextFeatureId(), remoteEntity, featureType, role)
	remoteEntity.AddFeature(remoteFeature)
	return remoteFeature
}
