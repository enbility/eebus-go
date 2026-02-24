package lpp

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/features/server"
	internal "github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *LPP) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	// subscribe to heartbeat when a remote entity binds to our loadControl server
	if payload.EventType == spineapi.EventTypeBindingChange &&
		payload.ChangeType == spineapi.ElementChangeAdd &&
		payload.LocalFeature != nil &&
		payload.LocalFeature.Type() == model.FeatureTypeTypeLoadControl &&
		payload.LocalFeature.Role() == model.RoleTypeServer &&
		e.IsScenarioAvailableAtEntity(payload.Entity, 1) {
		e.subscribeHeartbeat(payload.Entity)
		return
	}

	if internal.IsHeartbeat(payload) && e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateHeartbeat)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate ||
		payload.CmdClassifier == nil ||
		*payload.CmdClassifier != model.CmdClassifierTypeWrite {
		return
	}

	if !e.IsScenarioAvailableAtEntity(payload.Entity, 1) {
		return
	}

	switch payload.Data.(type) {
	case *model.LoadControlLimitListDataType:
		serverF := e.LocalEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)

		if payload.Function != model.FunctionTypeLoadControlLimitListData ||
			payload.LocalFeature != serverF {
			return
		}

		e.loadControlLimitDataUpdate(payload)
	case *model.DeviceConfigurationKeyValueListDataType:
		serverF := e.LocalEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceConfiguration, model.RoleTypeServer)

		if payload.Function != model.FunctionTypeDeviceConfigurationKeyValueListData ||
			payload.LocalFeature != serverF {
			return
		}

		e.configurationDataUpdate(payload)
	}
}

// subscribe to the DeviceDiagnosis of the entity that created a binding
func (e *LPP) subscribeHeartbeat(entity spineapi.EntityRemoteInterface) {
	if localDeviceDiag, err := client.NewDeviceDiagnosis(e.LocalEntity, entity); err == nil {
		e.heartbeatDiag = localDeviceDiag
		if !localDeviceDiag.HasSubscription() {
			if _, err := localDeviceDiag.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if _, err := localDeviceDiag.RequestHeartbeat(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit data was updated
func (e *LPP) loadControlLimitDataUpdate(payload spineapi.EventPayload) {
	if lc, err := server.NewLoadControl(e.LocalEntity); err == nil {
		filter := model.LoadControlLimitDescriptionDataType{
			LimitType:      util.Ptr(model.LoadControlLimitTypeTypeSignDependentAbsValueLimit),
			LimitCategory:  util.Ptr(model.LoadControlCategoryTypeObligation),
			ScopeType:      util.Ptr(model.ScopeTypeTypeActivePowerLimit),
			LimitDirection: util.Ptr(model.EnergyDirectionTypeProduce),
		}
		if lc.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
		}
	}
}

// the configuration key data was updated
func (e *LPP) configurationDataUpdate(payload spineapi.EventPayload) {
	if dc, err := server.NewDeviceConfiguration(e.LocalEntity); err == nil {
		filter := model.DeviceConfigurationKeyValueDescriptionDataType{
			KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeProductionActivePowerLimit),
		}
		if dc.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeProductionActivePowerLimit)
		}
		filter = model.DeviceConfigurationKeyValueDescriptionDataType{
			KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeFailsafeDurationMinimum),
		}
		if dc.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateFailsafeDurationMinimum)
		}
	}
}
