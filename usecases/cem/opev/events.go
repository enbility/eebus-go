package opev

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *OPEV) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		e.evConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.ElectricalConnectionPermittedValueSetListDataType:
		e.evElectricalPermittedValuesUpdate(payload)

	case *model.LoadControlLimitDescriptionListDataType:
		e.evLoadControlLimitDescriptionDataUpdate(payload.Entity)

	case *model.LoadControlLimitListDataType:
		e.evLoadControlLimitDataUpdate(payload)
	}
}

// TODOS:
//   - wie gehen wir damit um, wenn subscriptions nicht erfolgreich sind
//   - wie gehen wir damit um, wenn bindings nicht erfolgreich sind, z.b. write implementierungen müssen prüfen ob ein binding exisiteirt
//     könnte prüfen ob ein binding exisitert, wenn nicht gleich einen fehler zurückliefern
//   - wie gehen wir damit um, wenn descriptions nicht abgefragt werden können
//   - wie, wenn die einen Fehler zurückgeben

// an EV was connected
func (e *OPEV) evConnected(entity spineapi.EntityRemoteInterface) {
	// initialise features, e.g. subscriptions, descriptions
	if evLoadControl, err := client.NewLoadControl(e.LocalEntity, entity); err == nil {
		if !evLoadControl.HasSubscription() {
			if _, err := evLoadControl.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if !evLoadControl.HasBinding() {
			if _, err := evLoadControl.Bind(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get descriptions
		if _, err := evLoadControl.RequestLimitDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}

		// get constraints
		if _, err := evLoadControl.RequestLimitConstraints(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit description data of an EV was updated
func (e *OPEV) evLoadControlLimitDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if evLoadControl, err := client.NewLoadControl(e.LocalEntity, entity); err == nil {
		// get values
		if _, err := evLoadControl.RequestLimitData(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// the load control limit data of an EV was updated
func (e *OPEV) evLoadControlLimitDataUpdate(payload spineapi.EventPayload) {
	lc, err := client.NewLoadControl(e.LocalEntity, payload.Entity)
	if err != nil {
		return
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeObligation),
		ScopeType:     util.Ptr(model.ScopeTypeTypeOverloadProtection),
	}
	if lc.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
	}
}

// the electrical connection permitted value sets data of an EV was updated
func (e *OPEV) evElectricalPermittedValuesUpdate(payload spineapi.EventPayload) {
	if ec, err := client.NewElectricalConnection(e.LocalEntity, payload.Entity); err == nil {
		filter := model.ElectricalConnectionParameterDescriptionDataType{
			AcMeasuredPhases: util.Ptr(model.ElectricalConnectionPhaseNameTypeA),
		}
		data, err := ec.GetParameterDescriptionsForFilter(filter)
		if err != nil || len(data) == 0 || data[0].ParameterId == nil {
			return
		}

		filter = model.ElectricalConnectionParameterDescriptionDataType{
			ParameterId: data[0].ParameterId,
		}
		values, err := ec.GetParameterDescriptionsForFilter(filter)
		if err != nil || values == nil {
			return
		}

		// Scenario 6
		filter1 := model.ElectricalConnectionParameterDescriptionDataType{
			ElectricalConnectionId: values[0].ElectricalConnectionId,
			ParameterId:            values[0].ParameterId,
		}
		if ec.CheckEventPayloadDataForFilter(payload.Data, filter1) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateCurrentLimits)
		}
	}
}
