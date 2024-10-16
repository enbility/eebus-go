package ohpcf

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
)

// handle SPINE events
func (e *OHPCF) HandleEvent(payload spineapi.EventPayload) {
	// only about events from a compressor entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		// get the smart energy management data from the remote entity
		e.connected(payload.Entity)
	}
}

func (e *OHPCF) connected(entity spineapi.EntityRemoteInterface) {
	smartEnergyManagement, err := client.NewSmartEnergyManagementPs(e.LocalEntity, entity)
	if err != nil || smartEnergyManagement == nil {
		return
	}
	if _, err := smartEnergyManagement.RequestData(); err != nil {
		logging.Log().Debug(err)
	}
}
