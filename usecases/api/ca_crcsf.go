package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Actor: Configuration Appliance
// UseCase: Configuration of Room Cooling System Function
type CaCRCSFInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the cooling operation modes supported by the HVAC room
	//
	// parameters:
	//   - entity: the entity of the HVAC room
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	OperationModes(entity spineapi.EntityRemoteInterface) ([]HvacOperationModeType, error)

	// return the current cooling operation mode of the HVAC room
	//
	// parameters:
	//   - entity: the entity of the HVAC room
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such data is (yet) available
	//   - and others
	CurrentOperationMode(entity spineapi.EntityRemoteInterface) (HvacOperationModeType, error)

	// set the cooling operation mode of the HVAC room
	//
	// parameters:
	//   - entity: the entity of the HVAC room
	//   - mode: the cooling operation mode to set
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
}
