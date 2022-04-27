package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

func CreateRemoteDeviceAndFeature(entityId uint, featureType model.FeatureTypeType, role model.RoleType, sender Sender) *FeatureRemoteImpl {
	remoteDevice := NewDeviceRemoteImpl(model.AddressDeviceType("remoteDevice"), sender)
	remoteEntity := NewEntityRemoteImpl(remoteDevice, model.EntityTypeTypeEVSE, []model.AddressEntityType{model.AddressEntityType(entityId)})
	remoteDevice.addEntity(remoteEntity)
	remoteFeature := NewFeatureRemoteImpl(remoteEntity.NextFeatureId(), remoteEntity, featureType, role)
	remoteEntity.AddFeature(remoteFeature)
	return remoteFeature
}
