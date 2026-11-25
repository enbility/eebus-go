package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Actor: Monitoring Appliance
// UseCase: Monitoring of Room Temperature
type MaMRTInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the HVAC room temperature converted to the requested unit
	//
	// parameters:
	//   - entity: the entity of the device (e.g. HVAC room)
	//   - unit: the requested temperature unit (degC, degF or K)
	//
	// possible errors:
	//   - ErrDataNotAvailable if no such value is (yet) available
	//   - ErrDataInvalid if the currently available data is invalid and should be ignored
	//   - and others
	Temperature(entity spineapi.EntityRemoteInterface, unit model.UnitOfMeasurementType) (float64, error)
}
