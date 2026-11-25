package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
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
	//   - resultCB: callback for the device result; a non-zero ResultData.ErrorNumber signals a rejected write
	//
	// possible errors:
	//   - ErrNotSupported if the operation mode is not changeable or not supported
	//   - ErrDataNotAvailable if the required data is not (yet) available
	//   - and others
	WriteOperationMode(
		entity spineapi.EntityRemoteInterface,
		mode HvacOperationModeType,
		resultCB func(result model.ResultDataType, msgCounter model.MsgCounterType),
	) (*model.MsgCounterType, error)

	// Scenario 2

	// start the one-time DHW loading overrun of the DHW circuit
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//   - resultCB: callback for the device result; a non-zero ResultData.ErrorNumber signals a rejected write
	//
	// possible errors:
	//   - ErrNotSupported if the overrun status is not changeable
	//   - ErrDataNotAvailable if the required data is not (yet) available
	//   - and others
	StartOneTimeDhw(
		entity spineapi.EntityRemoteInterface,
		resultCB func(result model.ResultDataType, msgCounter model.MsgCounterType),
	) (*model.MsgCounterType, error)

	// Scenario 3

	// stop the one-time DHW loading overrun of the DHW circuit
	//
	// parameters:
	//   - entity: the entity of the DHW circuit
	//   - resultCB: callback for the device result; a non-zero ResultData.ErrorNumber signals a rejected write
	//
	// possible errors:
	//   - ErrNotSupported if the overrun status is not changeable
	//   - ErrDataNotAvailable if the required data is not (yet) available
	//   - and others
	StopOneTimeDhw(
		entity spineapi.EntityRemoteInterface,
		resultCB func(result model.ResultDataType, msgCounter model.MsgCounterType),
	) (*model.MsgCounterType, error)
}
