package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Actor: Configuration Appliance
// UseCase: Configuration of DHW Temperature
type CaCDTInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the current DHW temperature setpoints
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	Setpoints(entity spineapi.EntityRemoteInterface) ([]Setpoint, error)

	// return the constraints for the DHW temperature setpoints
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	SetpointConstraints(entity spineapi.EntityRemoteInterface) ([]SetpointConstraints, error)

	// write the DHW temperature setpoint for a DHW operation mode
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//   - mode: the DHW operation mode the setpoint is used for
	//   - degC: the temperature setpoint in degree Celsius
	//   - resultCB: called when the remote device returns a result; may be nil
	//
	// possible errors:
	//   - ErrNotSupported if the setpoint or write operation is not changeable
	//   - ErrDataNotAvailable if the required data is not (yet) available
	//   - ErrDataInvalid if the value is outside the constraints or not on a valid step
	//   - and others
	WriteSetpoint(
		entity spineapi.EntityRemoteInterface,
		mode HvacOperationModeType,
		degC float64,
		resultCB func(result model.ResultDataType, msgCounter model.MsgCounterType),
	) (*model.MsgCounterType, error)
}
