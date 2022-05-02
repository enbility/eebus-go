package usecase

import (
	"errors"
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

// Interface for the evseCC use case for CEM device
type UCEvseCCCemDelegate interface {
	// handle device state updates from the remote EVSE device
	HandleEVSEDeviceState(ski string, failure bool, errorCode string)
}

// EVSE Commissioning and Configuration Use Case implementation
type UCEvseCC struct {
	*UseCaseImpl
	service *service.EEBUSService

	// only required by CEM
	CemDelegate UCEvseCCCemDelegate
}

// Register the use case
func RegisterUCEvseCC(service *service.EEBUSService) UCEvseCC {
	entity := service.LocalEntity()

	// add the use case
	useCase := &UCEvseCC{
		UseCaseImpl: NewUseCase(
			entity,
			model.UseCaseNameTypeEVSECommissioningAndConfiguration,
			[]model.UseCaseScenarioSupportType{1, 2}),
		service: service,
	}

	// only the CEM needs to subscribe to get incoming EVSE events
	if useCase.Actor == model.UseCaseActorTypeCEM {
		spine.Events.Subscribe(useCase)
	}

	// add the features
	switch useCase.Actor {
	case model.UseCaseActorTypeEVSE:
		{
			f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceClassification, model.RoleTypeServer)
			f.SetDescriptionString("Device Classification Server")

			entity.AddFeature(f)
		}
		{
			f := spine.NewFeatureLocalImpl(entity.NextFeatureId(), entity, model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
			f.SetDescriptionString("Device Diagnosis Server")

			// Set the initial state
			deviceDiagnosisStateDate := &model.DeviceDiagnosisStateDataType{
				OperatingState: util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation),
			}
			f.SetData(model.FunctionTypeDeviceDiagnosisStateData, deviceDiagnosisStateDate)

			entity.AddFeature(f)
		}
	case model.UseCaseActorTypeCEM:
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
	}

	return *useCase
}

// public method to allow updating the EVSE device state
// this will be sent to the CEM remote device
func (r *UCEvseCC) UpdateEVSEErrorState(failure bool, errorCode string) {
	deviceDiagnosisStateDate := &model.DeviceDiagnosisStateDataType{}
	if failure {
		deviceDiagnosisStateDate.OperatingState = util.Ptr(model.DeviceDiagnosisOperatingStateTypeFailure)
		deviceDiagnosisStateDate.LastErrorCode = util.Ptr(model.LastErrorCodeType(errorCode))
	} else {
		deviceDiagnosisStateDate.OperatingState = util.Ptr(model.DeviceDiagnosisOperatingStateTypeNormalOperation)
	}
	r.notifyDeviceDiagnosisState(deviceDiagnosisStateDate)
}

// Internal EventHandler Interface
func (r *UCEvseCC) HandleEvent(payload spine.EventPayload) {
	switch payload.EventType {
	case spine.EventTypeDeviceChange:
		if payload.ChangeType == spine.ElementChangeAdd {
			r.requestManufacturer(payload.Device)
			r.requestDeviceDiagnosisState(payload.Device)
		}
	case spine.EventTypeDataChange:
		if payload.ChangeType == spine.ElementChangeUpdate {
			switch payload.Data.(type) {
			case *model.DeviceDiagnosisStateDataType:
				if r.CemDelegate == nil {
					return
				}

				deviceDiagnosisStateData := payload.Data.(model.DeviceDiagnosisStateDataType)
				failure := *deviceDiagnosisStateData.OperatingState == model.DeviceDiagnosisOperatingStateTypeFailure
				r.CemDelegate.HandleEVSEDeviceState(payload.Ski, failure, string(*deviceDiagnosisStateData.LastErrorCode))
			}
		}
	}
}

// request DeviceClassificationManufacturerData from a remote device
func (r *UCEvseCC) requestManufacturer(remoteDevice *spine.DeviceRemoteImpl) {
	featureLocal, featureRemote, err := r.getLocalClientAndRemoteServerFeatures(model.FeatureTypeTypeDeviceClassification, remoteDevice)
	if err != nil {
		fmt.Println(err)
		return
	}

	requestChannel := make(chan *model.DeviceClassificationManufacturerDataType)
	featureLocal.RequestData(model.FunctionTypeDeviceClassificationManufacturerData, featureRemote, requestChannel)
}

// request DeviceDiagnosisStateData from a remote device
func (r *UCEvseCC) requestDeviceDiagnosisState(remoteDevice *spine.DeviceRemoteImpl) {
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

// notify remote devices about the new device diagnosis state
func (r *UCEvseCC) notifyDeviceDiagnosisState(operatingState *model.DeviceDiagnosisStateDataType) {
	remoteDevice := r.service.RemoteDeviceOfType(model.DeviceTypeTypeEnergyManagementSystem)
	if remoteDevice == nil {
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

// internal helper method for getting local and remote feature for a given featureType and a given remoteDevice
func (r *UCEvseCC) getLocalClientAndRemoteServerFeatures(featureType model.FeatureTypeType, remoteDevice *spine.DeviceRemoteImpl) (spine.FeatureLocal, *spine.FeatureRemoteImpl, error) {
	featureLocal := r.Entity.Device().FeatureByTypeAndRole(featureType, model.RoleTypeClient)
	featureRemote := remoteDevice.FeatureByTypeAndRole(featureType, model.RoleTypeServer)

	if featureLocal == nil || featureRemote == nil {
		return nil, nil, errors.New("local or remote feature not found")
	}

	return featureLocal, featureRemote, nil
}
