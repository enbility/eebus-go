package evsecc

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// the manufacturer data of an EVSE
// returns deviceName, serialNumber, error
func (e *CemEVSECC) ManufacturerData(
	entity spineapi.EntityRemoteInterface,
) (
	api.ManufacturerData,
	error,
) {
	if !e.IsCompatibleEntity(entity) {
		return api.ManufacturerData{}, api.ErrNoCompatibleEntity
	}

	return internal.ManufacturerData(e.LocalEntity, entity)
}

// the operating state data of an EVSE
// returns operatingState, lastErrorCode, error
func (e *CemEVSECC) OperatingState(
	entity spineapi.EntityRemoteInterface,
) (
	model.DeviceDiagnosisOperatingStateType, string, error,
) {
	operatingState := model.DeviceDiagnosisOperatingStateTypeNormalOperation
	lastErrorCode := ""

	if !e.IsCompatibleEntity(entity) {
		return operatingState, lastErrorCode, api.ErrNoCompatibleEntity
	}

	evseDeviceDiagnosis, err := client.NewDeviceDiagnosis(e.LocalEntity, entity)
	if err != nil {
		return operatingState, lastErrorCode, err
	}

	data, err := evseDeviceDiagnosis.GetState()
	if err != nil {
		return operatingState, lastErrorCode, err
	}

	if data.OperatingState != nil {
		operatingState = *data.OperatingState
	}
	if data.LastErrorCode != nil {
		lastErrorCode = string(*data.LastErrorCode)
	}

	return operatingState, lastErrorCode, nil
}
