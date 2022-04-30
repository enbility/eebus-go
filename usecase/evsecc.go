package usecase

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type evseCC struct {
	*UseCaseImpl
}

func RegisterEvseCC(service *service.EEBUSService) {
	entity := service.LocalEntity()

	// add the use case
	useCase := &evseCC{
		UseCaseImpl: NewUseCase(
			entity,
			model.UseCaseNameTypeEVSECommissioningAndConfiguration,
			[]model.UseCaseScenarioSupportType{1, 2}),
	}

	if useCase.Actor == model.UseCaseActorTypeCEM {
		spine.Events.Subscribe(useCase)
	}

	// add the features
	{
		f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient)
		f.SetDescriptionString("Device Classification")

		entity.AddFeature(f)
	}
	switch useCase.Actor {
	case model.UseCaseActorTypeEVSE:
		{
			f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
			f.SetDescriptionString("Device Diagnosis")
			entity.AddFeature(f)
		}
	case model.UseCaseActorTypeCEM:
		{
			f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
			f.SetDescriptionString("Device Diagnosis")
			entity.AddFeature(f)
		}
	}
}

// EventHandler Interface
func (r *evseCC) HandleEvent(payload spine.EventPayload) {
	if payload.EventType == spine.EventTypeDeviceChange && payload.ChangeType == spine.ElementChangeAdd {
		if payload.Device.DeviceType() == model.DeviceTypeTypeChargingStation {
			r.requestManufacturer(payload.Device)
			// TOV-TODO: Add request diagnosis state data
		}
	}
}

func (r *evseCC) requestManufacturer(remoteDevice *spine.DeviceRemoteImpl) {
	featureLocal := r.Entity.Device().FeatureByTypeAndRole(model.FeatureTypeTypeDeviceClassification, model.RoleTypeClient)
	featureRemote := remoteDevice.FeatureByTypeAndRole(model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)

	if featureLocal == nil || featureRemote == nil {
		fmt.Println("ERROR: local or remote feature not found")
		return
	}

	requestChannel := make(chan *model.DeviceClassificationManufacturerDataType)
	featureLocal.RequestData(model.FunctionTypeDeviceClassificationManufacturerData, featureRemote, requestChannel)
}
