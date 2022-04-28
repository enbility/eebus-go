package entity

import (
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

//  Entities:
//   e[1] type=EVSE
//  Features:
//   e[1] f-0 server.DeviceDiagnosis
//    {RO} deviceDiagnosisStateData
//    {RO} deviceDiagnosisHeartbeatData
//   e[1] f-1 client.DeviceClassification
//   e[1] f-2 client.DeviceDiagnosis
func NewEVSE(device *spine.DeviceLocalImpl, address []model.AddressEntityType) *spine.EntityLocalImpl {
	entityType := model.EntityTypeTypeEVSE
	entity := spine.NewEntityLocalImpl(device, entityType, address)

	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
		f.SetDescriptionString("Device Diagnosis EVSE")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient)
		f.SetDescriptionString("EMS Device Classification")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
		f.SetDescriptionString("Device Diagnosis Heartbeat")
		entity.AddFeature(f)
	}

	return entity
}
