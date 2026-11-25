package vhan

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
)

// return the heating area name of the remote entity, scenario derived from
// entity type: heating circuit (1), heating zone (2) or HVAC room (3).
func (e *VHAN) Name(entity spineapi.EntityRemoteInterface) (string, error) {
	if !e.IsCompatibleEntityType(entity) {
		return "", api.ErrNoCompatibleEntity
	}

	deviceClassification, err := client.NewDeviceClassification(e.LocalEntity, entity)
	if err != nil {
		return "", err
	}

	data, err := deviceClassification.GetUserData()
	if err != nil || data == nil || data.UserLabel == nil {
		return "", api.ErrDataNotAvailable
	}

	return string(*data.UserLabel), nil
}
