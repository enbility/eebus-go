package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Monitoring Appliance
// UseCase: Monitoring of Room Heating System Function
type MaMRHSFInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the heating operation modes supported by the HVAC room
	//
	// parameters:
	//   - entity: the entity of the HVAC room
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	OperationModes(entity spineapi.EntityRemoteInterface) ([]HvacOperationModeType, error)

	// return the current heating operation mode of the HVAC room
	//
	// parameters:
	//   - entity: the entity of the HVAC room
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	CurrentOperationMode(entity spineapi.EntityRemoteInterface) (HvacOperationModeType, error)
}
