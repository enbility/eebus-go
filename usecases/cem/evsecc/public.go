package evsecc

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// the manufacturer data of an EVSE
// returns deviceName, serialNumber, error
func (e *EVSECC) ManufacturerData(
	entity spineapi.EntityRemoteInterface,
) (
	ucapi.ManufacturerData,
	error,
) {
	if !e.IsCompatibleEntityType(entity) {
		return ucapi.ManufacturerData{}, api.ErrNoCompatibleEntity
	}

	return internal.ManufacturerData(e.LocalEntity, entity)
}

// the operating state data of an EVSE
// returns operatingState, lastErrorCode, error
func (e *EVSECC) OperatingState(
	entity spineapi.EntityRemoteInterface,
) (
	model.DeviceDiagnosisOperatingStateType, string, error,
) {
	operatingState := model.DeviceDiagnosisOperatingStateTypeNormalOperation
	lastErrorCode := ""

	if !e.IsCompatibleEntityType(entity) {
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
