package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Configuration Appliance
// UseCase: Configuration of DHW System Function
type CaCDSFInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the DHW operation modes supported by the DHW circuit
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	OperationModes(entity spineapi.EntityRemoteInterface) ([]HvacOperationModeType, error)

	// return the current DHW operation mode of the DHW circuit
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	CurrentOperationMode(entity spineapi.EntityRemoteInterface) (HvacOperationModeType, error)

	// set the DHW operation mode of the DHW circuit
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//   - mode: the DHW operation mode to set
	//
	// possible errors:
	//   - ErrNotSupported if the operation mode is not changeable or not supported
	//   - ErrDataNotAvailable if the required data is not (yet) available
	//   - and others
	WriteOperationMode(entity spineapi.EntityRemoteInterface, mode HvacOperationModeType) error

	// Scenario 2

	// start the one-time DHW loading overrun of the DHW circuit
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//
	// possible errors:
	//   - ErrNotSupported if the overrun status is not changeable
	//   - ErrDataNotAvailable if the required data is not (yet) available
	//   - and others
	StartOneTimeDhw(entity spineapi.EntityRemoteInterface) error

	// Scenario 3

	// stop the one-time DHW loading overrun of the DHW circuit
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//
	// possible errors:
	//   - ErrNotSupported if the overrun status is not changeable
	//   - ErrDataNotAvailable if the required data is not (yet) available
	//   - and others
	StopOneTimeDhw(entity spineapi.EntityRemoteInterface) error
}
