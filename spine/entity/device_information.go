package entity

import (
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

//  Entities:
//   e[0] type=DeviceInformation
//  Features:
//   e[0] f-0 special.NodeManagement
//    {RO} nodeManagementDetailedDiscoveryData
//    {--} nodeManagementSubscriptionRequestCall
//    {--} nodeManagementBindingRequestCall
//    {--} nodeManagementSubscriptionDeleteCall
//    {--} nodeManagementBindingDeleteCall
//    {RO} nodeManagementSubscriptionData
//    {RO} nodeManagementBindingData
//    {RO} nodeManagementUseCaseData
//   e[0] f-1 server.DeviceClassification
//    {RO} deviceClassificationManufacturerData
func NewDeviceInformation(device *spine.DeviceLocalImpl) *spine.EntityLocalImpl {
	entityType := model.EntityTypeType(model.EntityTypeTypeDeviceInformation)
	entity := spine.NewEntityLocalImpl(device, entityType, []model.AddressEntityType{model.AddressEntityType(spine.DeviceInformationEntityId)})

	{
		// TODO
		//f := feature.NewNodeManagementLocal(entity.NextFeatureId(), entity)
		//entity.Add(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
		entity.AddFeature(f)
	}

	return entity
}
