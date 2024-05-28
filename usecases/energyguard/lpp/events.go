package lpp

import (
	"github.com/enbility/eebus-go/features/client"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *LPP) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntity(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		e.connected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
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
func (e *LPP) connected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if loadControl, err := client.NewLoadControl(e.LocalEntity, entity); err == nil {
		if _, err := loadControl.Subscribe(); err != nil {
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

// the load control limit description data was updated
func (e *LPP) loadControlLimitDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if loadControl, err := client.NewLoadControl(e.LocalEntity, entity); err == nil {
		// get values
		if _, err := loadControl.RequestLimitData(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit data was updated
func (e *LPP) loadControlLimitDataUpdate(payload spineapi.EventPayload) {
	if lc, err := client.NewLoadControl(e.LocalEntity, payload.Entity); err == nil {
		filter := model.LoadControlLimitDescriptionDataType{
			LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
			LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
			LimitDirection: util.Ptr(model.EnergyDirectionTypeProduce),
			ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
		}
		if lc.CheckEventPayloadDataForFilter(payload.Data, filter) {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
		}
	}
}

// the configuration key description data was updated
func (e *LPP) configurationDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity); err == nil {
		// key value descriptions received, now get the data
		if _, err := deviceConfiguration.RequestKeyValues(); err != nil {
			logging.Log().Error("Error getting configuration key values:", err)
		}
	}
}

// the configuration key data was updated
func (e *LPP) configurationDataUpdate(payload spineapi.EventPayload) {
	if dc, err := client.NewDeviceConfiguration(e.LocalEntity, payload.Entity); err == nil {
		filter := model.DeviceConfigurationKeyValueDescriptionDataType{
			KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
		}
		if dc.CheckEventPayloadDataForFilter(payload.Data, filter) {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeProductionActivePowerLimit)
		}
		filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum)
		if dc.CheckEventPayloadDataForFilter(payload.Data, filter) {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeDurationMinimum)
		}
	}
}
