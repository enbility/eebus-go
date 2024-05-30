package lpc

import (
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *LPC) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		e.connected(payload.Entity)
		return
	}

	if internal.IsEntityDisconnected(payload) {
		e.UseCaseDataUpdate(payload, e.EventCB, UseCaseSupportUpdate)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.NodeManagementUseCaseDataType:
		e.UseCaseDataUpdate(payload, e.EventCB, UseCaseSupportUpdate)

	case *model.LoadControlLimitDescriptionListDataType:
		e.loadControlLimitDescriptionDataUpdate(payload.Entity)

	case *model.LoadControlLimitListDataType:
		e.loadControlLimitDataUpdate(payload)

	case *model.DeviceConfigurationKeyValueDescriptionListDataType:
		e.configurationDescriptionDataUpdate(payload.Entity)

	case *model.DeviceConfigurationKeyValueListDataType:
		e.configurationDataUpdate(payload)
	}
}

// the remote entity was connected
func (e *LPC) connected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if loadControl, err := client.NewLoadControl(e.LocalEntity, entity); err == nil {
		if _, err := loadControl.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}

		if _, err := loadControl.Bind(); err != nil {
			logging.Log().Debug(err)
		}

		// get descriptions
		if _, err := loadControl.RequestLimitDescriptions(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if localDeviceDiag, err := client.NewDeviceDiagnosis(e.LocalEntity, entity); err == nil {
		if _, err := localDeviceDiag.Subscribe(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// // use case data updated or entity was removed,
// // update the supported remote entities list and send event if neccessary
// func (e *EgLPC) useCaseDataUpdate(payload spineapi.EventPayload) {
// 	// entity was removed, so remove it from the list
// 	if internal.IsEntityDisconnected(payload) {
// 		if e.HasRemoteEntity(payload.Entity) {
// 			e.RemoveRemoteEntity(payload.Entity)
// 		}

// 		return
// 	}

// 	// entity changed usecase data
// 	if result, err := e.IsUseCaseSupported(payload.Entity); err == nil && result {
// 		scenarios := e.SupportedUseCaseScenarios(payload.Entity)
// 		e.SetRemoteEntityScenarios(payload.Entity, scenarios)

// 		return
// 	}

// 	// entity does not support the use case, maybe support was removed
// 	e.RemoveRemoteEntity(payload.Entity)
// }

// the load control limit description data was updated
func (e *LPC) loadControlLimitDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if loadControl, err := client.NewLoadControl(e.LocalEntity, entity); err == nil {
		// get values
		if _, err := loadControl.RequestLimitData(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit data was updated
func (e *LPC) loadControlLimitDataUpdate(payload spineapi.EventPayload) {
	if lc, err := client.NewLoadControl(e.LocalEntity, payload.Entity); err == nil {
		filter := model.LoadControlLimitDescriptionDataType{
			LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
			LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
			LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
		}
		if lc.CheckEventPayloadDataForFilter(payload.Data, filter) {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
		}
	}
}

// the configuration key description data was updated
func (e *LPC) configurationDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity); err == nil {
		// key value descriptions received, now get the data
		if _, err := deviceConfiguration.RequestKeyValues(); err != nil {
			logging.Log().Error("Error getting configuration key values:", err)
		}
	}
}

// the configuration key data was updated
func (e *LPC) configurationDataUpdate(payload spineapi.EventPayload) {
	if dc, err := client.NewDeviceConfiguration(e.LocalEntity, payload.Entity); err == nil {
		filter := model.DeviceConfigurationKeyValueDescriptionDataType{
			KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeConsumptionActivePowerLimit),
		}
		if dc.CheckEventPayloadDataForFilter(payload.Data, filter) {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeConsumptionActivePowerLimit)
		}
		filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum)
		if dc.CheckEventPayloadDataForFilter(payload.Data, filter) {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeDurationMinimum)
		}
	}
}
