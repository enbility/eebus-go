package cdt

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	"github.com/enbility/ship-go/util"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// HandleEvent handles events for the CDT use case.
func (e *CDT) HandleEvent(payload spineapi.EventPayload) {
	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityConnected(payload) {
		e.dhwCircuitconnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.SetpointDescriptionListDataType:
		e.setpointDescriptionsUpdate(payload)

	case *model.SetpointConstraintsListDataType:
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateSetpointConstraints)

	case *model.SetpointListDataType:
		e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateSetpoints)
	}
}

// setpointDescriptionsUpdate processes the necessary steps when setpoint descriptions are updated.
func (e *CDT) setpointDescriptionsUpdate(payload spineapi.EventPayload) {
	setPoint, err := client.NewSetpoint(e.LocalEntity, payload.Entity)
	if err != nil {
		logging.Log().Debug(err)
		return
	}

	// The setpointConstraintsListData and setpointListData reads should
	// be partial, using setpointId from setpointDescriptionListData.
	setpointDescriptions := payload.Data.(*model.SetpointDescriptionListDataType).SetpointDescriptionData
	for _, setpointDescription := range setpointDescriptions {
		constraintsSelector := &model.SetpointConstraintsListDataSelectorsType{
			SetpointId: setpointDescription.SetpointId,
		}
		if _, err := setPoint.RequestSetPointConstraints(constraintsSelector, nil); err != nil {
			logging.Log().Debug(err)
		}

		setpointSelector := &model.SetpointListDataSelectorsType{
			SetpointId: setpointDescription.SetpointId,
		}
		if _, err := setPoint.RequestSetPoints(setpointSelector, nil); err != nil {
			logging.Log().Debug(err)
		}
	}

	// Request setpoint relations to map operation modes to setpoints.
	if hvac, err := client.NewHvac(e.LocalEntity, payload.Entity); err == nil {
		if _, err := hvac.RequestHvacSystemFunctionSetPointRelations(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}

// dhwCircuitconnected processes required steps when a DHW Circuit is connected.
func (e *CDT) dhwCircuitconnected(entity spineapi.EntityRemoteInterface) {
	if hvac, err := client.NewHvac(e.LocalEntity, entity); err == nil {
		if !hvac.HasSubscription() {
			if _, err := hvac.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}
	}

	if setPoint, err := client.NewSetpoint(e.LocalEntity, entity); err == nil {
		if !setPoint.HasSubscription() {
			if _, err := setPoint.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		selector := &model.SetpointDescriptionListDataSelectorsType{
			ScopeType: util.Ptr(model.ScopeTypeTypeDhwTemperature),
		}
		if _, err := setPoint.RequestSetPointDescriptions(selector, nil); err != nil {
			logging.Log().Debug(err)
		}
	}
}
