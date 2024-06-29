package oscev

import (
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *OSCEV) HandleEvent(payload spineapi.EventPayload) {
	// most of the events are identical to OPEV, and OPEV is required to be used,
	// we don't handle the same events in here

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.ElectricalConnectionPermittedValueSetListDataType:
		e.evElectricalPermittedValuesUpdate(payload)

	case *model.LoadControlLimitListDataType:
		e.evLoadControlLimitDataUpdate(payload)
	}
}

// the load control limit data of an EV was updated
func (e *OSCEV) evLoadControlLimitDataUpdate(payload spineapi.EventPayload) {
	lc, err := client.NewLoadControl(e.LocalEntity, payload.Entity)
	if err != nil {
		return
	}

	filter := model.LoadControlLimitDescriptionDataType{
		LimitType:     util.Ptr(model.LoadControlLimitTypeTypeMaxValueLimit),
		LimitCategory: util.Ptr(model.LoadControlCategoryTypeRecommendation),
		ScopeType:     util.Ptr(model.ScopeTypeTypeSelfConsumption),
	}

	if lc.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateLimit)
	}
}

// the electrical connection permitted value sets data of an EV was updated
func (e *OSCEV) evElectricalPermittedValuesUpdate(payload spineapi.EventPayload) {
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
