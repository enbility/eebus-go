package spine

import (
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

func CreateRemoteDeviceAndFeature(entityId uint, featureType model.FeatureTypeType, role model.RoleType, sender Sender) *FeatureRemoteImpl {
	localDevice := NewDeviceLocalImpl("Vendor", "DeviceName", "SerialNumber", "DeviceCode", "Address", model.DeviceTypeTypeEnergyManagementSystem, model.NetworkManagementFeatureSetTypeSmart)

	remoteDevice := NewDeviceRemoteImpl(localDevice, "ski", nil, nil)
	remoteDevice.address = util.Ptr(model.AddressDeviceType("Address"))
	remoteDevice.sender = sender
	remoteEntity := NewEntityRemoteImpl(remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(entityId)})
	remoteDevice.addEntity(remoteEntity)
	remoteFeature := NewFeatureRemoteImpl(remoteEntity.NextFeatureId(), remoteEntity, featureType, role)
	remoteEntity.AddFeature(remoteFeature)
	return remoteFeature
}
