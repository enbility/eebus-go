package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
)

// Actor: Visualization Appliance
// UseCase: Visualization of Heating Area Name
type VaVHANInterface interface {
	api.UseCaseInterface

	// return the heating area name of the remote entity, scenario derived
	// from entity type: circuit (1), zone (2), room (3). See ErrDataNotAvailable.
	Name(entity spineapi.EntityRemoteInterface) (string, error)
}
