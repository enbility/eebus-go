package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Actor: Customer Energy Management
// UseCase: EVSE Commissioning and Configuration
type CemEVSECCInterface interface {
	api.UseCaseInterface

	// the manufacturer data of an EVSE
	//
	// parameters:
	//   - entity: the entity of the EV
	//
	// returns deviceName, serialNumber, error
	ManufacturerData(entity spineapi.EntityRemoteInterface) (ManufacturerData, error)

	// the operating state data of an EVSE
	//
	// parameters:
	//   - entity: the entity of the EV
	//
	// returns operatingState, lastErrorCode, error
	OperatingState(entity spineapi.EntityRemoteInterface) (model.DeviceDiagnosisOperatingStateType, string, error)
}
