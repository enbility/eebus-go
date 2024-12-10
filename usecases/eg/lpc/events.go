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
	if internal.IsEntityAdded(payload) {
		e.connected(payload.Entity)
		return
	}

	if internal.IsHeartbeat(payload) && e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateHeartbeat)
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
func (e *LPC) connected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if loadControl, err := client.NewLoadControl(e.LocalEntity, entity); err == nil {
		if !loadControl.HasSubscription() {
			if _, err := loadControl.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if !loadControl.HasBinding() {
			if _, err := loadControl.Bind(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get descriptions
		selector := &model.LoadControlLimitDescriptionListDataSelectorsType{
			LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
			LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
		}
		if _, err := loadControl.RequestLimitDescriptions(selector, nil); err != nil {
			logging.Log().Debug(err)
		}
	}

	if deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity); err == nil {
		if !deviceConfiguration.HasSubscription() {
			if _, err := deviceConfiguration.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if !deviceConfiguration.HasBinding() {
			if _, err := deviceConfiguration.Bind(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get descriptions
		// don't use selectors yet, as we would have to query 2 which could result in 2 full reads
		if _, err := deviceConfiguration.RequestKeyValueDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}

	if deviceDiagnosis, err := client.NewDeviceDiagnosis(e.LocalEntity, entity); err == nil {
		if !deviceDiagnosis.HasSubscription() {
			if _, err := deviceDiagnosis.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if _, err := deviceDiagnosis.RequestHeartbeat(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit description data was updated
func (e *LPC) loadControlLimitDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if loadControl, err := client.NewLoadControl(e.LocalEntity, entity); err == nil {
		// get values
		var selector *model.LoadControlLimitListDataSelectorsType
		filter := model.LoadControlLimitDescriptionDataType{
			LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
			LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
		}
		if descs, err := loadControl.GetLimitDescriptionsForFilter(filter); err == nil && len(descs) > 0 {
			selector = &model.LoadControlLimitListDataSelectorsType{
				LimitId: descs[0].LimitId,
			}
		}
		if _, err := loadControl.RequestLimitData(selector, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit data was updated
func (e *LPC) loadControlLimitDataUpdate(payload spineapi.EventPayload) {
	if lc, err := client.NewLoadControl(e.LocalEntity, payload.Entity); err == nil {
		filter := model.LoadControlLimitDescriptionDataType{
			LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
			LimitDirection: util.Ptr(model.EnergyDirectionTypeConsume),
			ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
		}
		if lc.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
		}
	}
}

// the configuration key description data was updated
func (e *LPC) configurationDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity); err == nil {
		// key value descriptions received, now get the data
		if _, err := deviceConfiguration.RequestKeyValues(nil, nil); err != nil {
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
		if dc.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeConsumptionActivePowerLimit)
		}
		filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum)
		if dc.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeDurationMinimum)
		}
	}
}
