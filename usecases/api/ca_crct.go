package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Configuration Appliance
// UseCase: Configuration of Room Cooling Temperature
type CaCRCTInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the current room cooling temperature setpoints
	//
	// parameters:
	//   - entity: the entity of the HVAC room
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	Setpoints(entity spineapi.EntityRemoteInterface) ([]Setpoint, error)

	// return the constraints for the room cooling temperature setpoints
	//
	// parameters:
	//   - entity: the entity of the HVAC room
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	SetpointConstraints(entity spineapi.EntityRemoteInterface) ([]SetpointConstraints, error)

	// write the room cooling temperature setpoint for a cooling operation mode
	//
	// parameters:
	//   - entity: the entity of the HVAC room
	//   - mode: the cooling operation mode the setpoint is used for (on, off or eco)
	//   - degC: the temperature setpoint in degree Celsius
	//
	// possible errors:
	//   - ErrNotSupported if the setpoint is not changeable or the mode is auto
	//   - ErrDataNotAvailable if the required data is not (yet) available
	//   - and others
	WriteSetpoint(entity spineapi.EntityRemoteInterface, mode HvacOperationModeType, degC float64) error
}
