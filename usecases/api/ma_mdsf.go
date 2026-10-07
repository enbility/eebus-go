package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Actor: Monitoring Appliance
// UseCase: Monitoring of DHW System Function
type MaMDSFInterface interface {
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

	// Scenario 2

	// return whether an overrun affecting the DHW system function is active
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	IsOverrunActive(entity spineapi.EntityRemoteInterface) (bool, error)

	// return the status of the one-time DHW overrun (active, running, finished or inactive)
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	OverrunStatus(entity spineapi.EntityRemoteInterface) (model.HvacOverrunStatusType, error)
}
