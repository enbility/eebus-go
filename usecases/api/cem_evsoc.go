package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Customer Energy Management
// UseCase: EV State Of Charge
type CemEVSOCInterface interface {
	api.UseCaseInterface

	// Scenario 1

	// return the EVscurrent state of charge of the EV or an error it is unknown
	//
	// parameters:
	//   - entity: the entity of the EV
	StateOfCharge(entity spineapi.EntityRemoteInterface) (float64, error)

	// Scenario 2 to 4 are not supported, as there is no EV supporting this as of today
}
