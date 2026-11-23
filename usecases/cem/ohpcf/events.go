package ohpcf

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (o *OHPCF) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !o.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		o.connected(payload.Entity)
	}

	if payload.Data == nil {
		return
	}

	switch payload.Data.(type) {
	case *model.SmartEnergyManagementPsDataType:
		o.loadSmartEnergyManagementPsDataType(payload)
	}
}

func (o *OHPCF) connected(entity spineapi.EntityRemoteInterface) {
	if semp, err := client.NewSmartEnergyManagementPs(o.LocalEntity, entity); err == nil {
		if !semp.HasSubscription() {
			if _, err := semp.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		if !semp.HasBinding() {
			if _, err := semp.Bind(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// read the current data as a subscription only delivers future updates
		if _, err := semp.RequestData(); err != nil {
			logging.Log().Debug(err)
		}
	}
}

func (o *OHPCF) loadSmartEnergyManagementPsDataType(payload spineapi.EventPayload) {
	data := payload.Data.(*model.SmartEnergyManagementPsDataType)

	if o.EventCB == nil {
		return
	}

	if len(data.Alternatives) == 1 {
		alternative := data.Alternatives[0]

		if len(alternative.PowerSequence) != 1 {
			return
		}

		request := alternative.PowerSequence[0]

		if len(request.PowerTimeSlot) == 1 &&
			request.PowerTimeSlot[0].ValueList != nil &&
			len(request.PowerTimeSlot[0].ValueList.Value) > 0 {
			for _, value := range request.PowerTimeSlot[0].ValueList.Value {
				if value.Value != nil &&
					value.ValueType != nil {
					if *value.ValueType == model.PowerTimeSlotValueTypeTypePower {
						o.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateRequestedPowerEstimate)
					} else if *value.ValueType == model.PowerTimeSlotValueTypeTypePowerMax {
						o.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateRequestedPowerMax)
					}
				}
			}
		}

		if request.OperatingConstraintsInterrupt != nil {
			if request.OperatingConstraintsInterrupt.IsStoppable != nil {
				// [OHPCF-011/5]
				// [OHPCF-012/3]
				o.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateConsumptionIsStoppable)
			}

			if request.OperatingConstraintsInterrupt.IsPausable != nil {
				// [OHPCF-011/6]
				// [OHPCF-012/3]
				o.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateConsumptionIsPausable)
			}
		}

		if request.Schedule != nil &&
			request.Schedule.StartTime != nil {
			// [OHPCF-012/1]
			o.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateConsumptionStartTime)
		}

		if request.State != nil &&
			request.State.State != nil {
			// [OHPCF-006], [OHPCF-011/1], [OHPCF-012/2], [OHPCF-022]
			o.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateConsumptionState)
		}

		if request.OperatingConstraintsDuration != nil {
			if request.OperatingConstraintsDuration.ActiveDurationMin != nil {
				// [OHPCF-008]
				o.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateMinimalRunDuration)
			}
			if request.OperatingConstraintsDuration.PauseDurationMin != nil {
				// [OHPCF-009]
				o.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateMinimalPauseDuration)
			}
		}

	} else if len(data.Alternatives) == 0 {
		// [OHPCF-003], [OHPCF-006/2]
		o.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateConsumptionState)
	}
}
