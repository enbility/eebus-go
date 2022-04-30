package usecase

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type evseCCDelegate interface {
	// handle error state updates from an EVSE device
	HandleEVSEErrorState(failure bool, errorCode string)
}

type evseCC struct {
	*UseCaseImpl
	service *service.EEBUSService

	// only required by CEM
	Delegate evseCCDelegate
}

func RegisterEvseCC(service *service.EEBUSService) evseCC {
	entity := service.LocalEntity()

	// add the use case
	useCase := &evseCC{
		UseCaseImpl: NewUseCase(
			entity,
			model.UseCaseNameTypeEVSECommissioningAndConfiguration,
			[]model.UseCaseScenarioSupportType{1, 2}),
		service: service,
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

			// Set the initial state
			deviceDiagnosisStateDate := &model.DeviceDiagnosisStateDataType{
				OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
			}
			f.SetData(model.FunctionTypeDeviceDiagnosisStateData, deviceDiagnosisStateDate)

			entity.AddFeature(f)
		}
	case model.UseCaseActorTypeCEM:
		{
			f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeClient)
			f.SetDescriptionString("Device Diagnosis")
			entity.AddFeature(f)
		}
	}

	return *useCase
}

func (r *evseCC) SetEVSEErrorState(failure bool, errorCode string) {
	deviceDiagnosisStateDate := &model.DeviceDiagnosisStateDataType{}
	if failure {
		deviceDiagnosisStateDate.OperatingState = util.Ptr(model.DeviceDiagnosisOperatingStateTypeFailure)
		deviceDiagnosisStateDate.LastErrorCode = util.Ptr(model.LastErrorCodeType(errorCode))
	} else {
		deviceDiagnosisStateDate.OperatingState = util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation)
	}
	r.notifyDeviceDiagnosisState(deviceDiagnosisStateDate)
}

// EventHandler Interface
func (r *evseCC) HandleEvent(payload spine.EventPayload) {
	switch payload.EventType {
	case spine.EventTypeDeviceChange:
		if payload.ChangeType == spine.ElementChangeAdd {
			if payload.Device.DeviceType() == model.DeviceTypeTypeChargingStation {
				r.requestManufacturer(payload.Device)
				r.requestDeviceDiagnosisState(payload.Device)
			}
		}
	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceDiagnosisStateDataType:
				if r.Delegate == nil {
					return
				}

				deviceDiagnosisStateData := payload.Data.(model.DeviceDiagnosisStateDataType)
				failure := *deviceDiagnosisStateData.OperatingState == model.DeviceDiagnosisOperatingStateTypeFailure
				r.Delegate.HandleEVSEErrorState(failure, string(*deviceDiagnosisStateData.LastErrorCode))
			}
		}
	}
}

func (r *evseCC) requestManufacturer(remoteDevice *spine.DeviceRemoteImpl) {
	featureLocal, featureRemote, err := r.getLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceClassification, remoteDevice)

	if err != nil {
		fmt.Println(err)
		return
	}

	requestChannel := make(chan *model.DeviceClassificationManufacturerDataType)
	featureLocal.RequestData(model.FunctionTypeDeviceClassificationManufacturerData, featureRemote, requestChannel)
}

func (r *evseCC) requestDeviceDiagnosisState(remoteDevice *spine.DeviceRemoteImpl) {
	featureLocal, featureRemote, err := r.getLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, remoteDevice)

	if err != nil {
		fmt.Println(err)
		return
	}

	requestChannel := make(chan *model.DeviceDiagnosisStateDataType)
	featureLocal.RequestData(model.FunctionTypeDeviceDiagnosisStateData, featureRemote, requestChannel)

	// subscribe to device diagnosis state updates
	remoteDevice.Sender().Subscribe(featureLocal.Address(), featureRemote.Address(), model.FeatureTypeTypeDeviceDiagnosis)
}

func (r *evseCC) notifyDeviceDiagnosisState(operatingState *model.DeviceDiagnosisStateDataType) {
	remoteDevice := r.service.RemoteDeviceOfType(model.DeviceTypeTypeEnergyManagementSystem)

	if remoteDevice == nil {
		fmt.Println("Remote device not found")
		return
	}

	featureLocal, featureRemote, err := r.getLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceDiagnosis, remoteDevice)

	if err != nil {
		fmt.Println(err)
		return
	}

	featureLocal.SetData(model.FunctionTypeDeviceDiagnosisStateData, operatingState)

	featureLocal.NotifyData(model.FunctionTypeDeviceDiagnosisStateData, featureRemote)
}

func (r *evseCC) getLocalClientAndRemoteServerFeatures(featureType model.FeatureTypeType, remoteDevice *spine.DeviceRemoteImpl) (spine.FeatureLocal, *spine.FeatureRemoteImpl, error) {
	featureLocal := r.Entity.Device().FeatureByTypeAndRole(featureType, model.RoleTypeClient)
	featureRemote := remoteDevice.FeatureByTypeAndRole(featureType, model.RoleTypeServer)

	if featureLocal == nil || featureRemote == nil {
		return nil, nil, errors.New("local or remote feature not found")
	}

	return featureLocal, featureRemote, nil
}
