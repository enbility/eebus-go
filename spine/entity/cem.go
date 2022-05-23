package entity

import (
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

//  Entities:
//   e[1] type=CEM, CEM Energy Guard
//  Features:
//   e[1] f-1 client.DeviceClassification - Device Classification
//   e[1] f-2 client.DeviceDiagnosis - Device Diagnosis
//   e[1] f-3 client.Measurement - Measurement for client
//   e[1] f-4 client.DeviceConfiguration - Device Configuration
//   e[1] f-5 server.DeviceDiagnosis - DeviceDiag
//    {RO} deviceDiagnosisStateData
//    {RO} deviceDiagnosisHeartbeatData
//   e[1] f-7 client.LoadControl - LoadControl client for CEM
//   e[1] f-8 client.Identification - EV identification
//   e[1] f-9 client.ElectricalConnection - Electrical Connection
func NewCEM(device *spine.DeviceLocalImpl, address []model.AddressEntityType) *spine.EntityLocalImpl {
	entityType := model.EntityTypeTypeCEM
	entity := spine.NewEntityLocalImpl(device, entityType, address)

	// UseCases
	localUseCaseManager := device.UseCaseManager()

	// UseCase EVSECommissioningAndConfiguration
	localUseCaseManager.Add(
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeEVSECommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2},
	)
	// UseCase EVCommissioningAndConfiguration
	localUseCaseManager.Add(
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeEVCommissioningAndConfiguration,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8},
	)
	// UseCase MeasurementOfElectricityDuringEVCharging
	localUseCaseManager.Add(
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeMeasurementOfElectricityDuringEVCharging,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3},
	)
	// UseCase OverloadProtectionByEVChargingCurrentCurtailment
	localUseCaseManager.Add(
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeOverloadProtectionByEVChargingCurrentCurtailment,
		model.SpecificationVersionType("1.0.1b"),
		[]model.UseCaseScenarioSupportType{1, 2, 3},
	)
	// UseCase OptimizationOfSelfConsumptionDuringEVCharging
	localUseCaseManager.Add(
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeOptimizationOfSelfConsumptionDuringEVCharging,
		model.SpecificationVersionType("1.0.1b"),
		[]model.UseCaseScenarioSupportType{1, 2, 3},
	)
	// UseCase EVStateOfCharge
	localUseCaseManager.Add(
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeEVStateOfCharge,
		model.SpecificationVersionType("1.0.0"),
		[]model.UseCaseScenarioSupportType{1},
	)
	// UseCase CoordinatedEVCharging
	localUseCaseManager.Add(
		model.UseCaseActorTypeCEM,
		model.UseCaseNameTypeCoordinatedEVCharging,
		model.SpecificationVersionType("1.0.1"),
		[]model.UseCaseScenarioSupportType{1, 2, 3, 4, 5, 6, 7, 8},
	)

	// Features
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient)
		f.SetDescriptionString("Device Classification Client")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
		f.SetDescriptionString("Device Diagnosis Client")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeMeasurement, model.RoleTypeClient)
		f.SetDescriptionString("Measurements Client")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeClient)
		f.SetDescriptionString("Device Configuration Client")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
		f.SetDescriptionString("Device Diagnoses Server")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeLoadControl, model.RoleTypeClient)
		f.SetDescriptionString("Load Control Client")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeIdentification, model.RoleTypeClient)
		f.SetDescriptionString("Identification Client")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeElectricalConnection, model.RoleTypeClient)
		f.SetDescriptionString("Electrical Connection Client")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeTimeSeries, model.RoleTypeClient)
		f.SetDescriptionString("Time Series Client")
		entity.AddFeature(f)
	}
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeIncentiveTable, model.RoleTypeClient)
		f.SetDescriptionString("Incentive Table Client")
		entity.AddFeature(f)
	}

	return entity
}
