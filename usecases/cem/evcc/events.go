package evcc

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *EVCC) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an EV entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		e.evConnected(payload)
		return
	} else if internal.IsEntityRemoved(payload) {
		e.evDisconnected(payload)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.DeviceConfigurationKeyValueDescriptionListDataType:
		e.evConfigurationDescriptionDataUpdate(payload.Entity)

	case *model.DeviceConfigurationKeyValueListDataType:
		e.evConfigurationDataUpdate(payload)

	case *model.DeviceDiagnosisOperatingStateType:
		e.evOperatingStateDataUpdate(payload)

	case *model.DeviceClassificationManufacturerDataType:
		e.evManufacturerDataUpdate(payload)

	case *model.ElectricalConnectionParameterDescriptionListDataType:
		e.evElectricalParamerDescriptionUpdate(payload.Entity)

	case *model.ElectricalConnectionPermittedValueSetListDataType:
		e.evElectricalPermittedValuesUpdate(payload)

	case *model.IdentificationListDataType:
		e.evIdentificationDataUpdate(payload)
	}
}

// an EV was connected
func (e *EVCC) evConnected(payload spineapi.EventPayload) {
	// initialise features, e.g. subscriptions, descriptions
	if evDeviceClassification, err := client.NewDeviceClassification(e.LocalEntity, payload.Entity); err == nil {
		if !evDeviceClassification.HasSubscription() {
			if _, err := evDeviceClassification.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get manufacturer details
		if _, err := evDeviceClassification.RequestManufacturerDetails(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evDeviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, payload.Entity); err == nil {
		if !evDeviceConfiguration.HasSubscription() {
			if _, err := evDeviceConfiguration.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get ev configuration data
		if _, err := evDeviceConfiguration.RequestKeyValueDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evDeviceDiagnosis, err := client.NewDeviceDiagnosis(e.LocalEntity, payload.Entity); err == nil {
		if !evDeviceDiagnosis.HasSubscription() {
			if _, err := evDeviceDiagnosis.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get device diagnosis state
		if _, err := evDeviceDiagnosis.RequestState(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evElectricalConnection, err := client.NewElectricalConnection(e.LocalEntity, payload.Entity); err == nil {
		if !evElectricalConnection.HasSubscription() {
			if _, err := evElectricalConnection.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get electrical connection parameter descriptions
		if _, err := evElectricalConnection.RequestParameterDescriptions(nil, nil); err != nil {
			logging.Log().Debug(err)
		}

		// get electrical permitted values descriptions
		if _, err := evElectricalConnection.RequestPermittedValueSets(nil, nil); err != nil {
			logging.Log().Debug(err)
		}
	}

	if evIdentification, err := client.NewIdentification(e.LocalEntity, payload.Entity); err == nil {
		if !evIdentification.HasSubscription() {
			if _, err := evIdentification.Subscribe(); err != nil {
				logging.Log().Debug(err)
			}
		}

		// get identification
		if _, err := evIdentification.RequestValues(); err != nil {
			logging.Log().Debug(err)
		}
	}

	if e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, EvConnected)
	}
}

// an EV was disconnected
func (e *EVCC) evDisconnected(payload spineapi.EventPayload) {
	if e.EventCB != nil {
		e.EventCB(payload.Ski, payload.Device, payload.Entity, EvDisconnected)
	}
}

// the configuration key description data of an EV was updated
func (e *EVCC) evConfigurationDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if evDeviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity); err == nil {
		// key value descriptions received, now get the data
		if _, err := evDeviceConfiguration.RequestKeyValues(nil, nil); err != nil {
			logging.Log().Error("Error getting configuration key values:", err)
		}
	}
}

// the configuration key data of an EV was updated
func (e *EVCC) evConfigurationDataUpdate(payload spineapi.EventPayload) {
	if dc, err := client.NewDeviceConfiguration(e.LocalEntity, payload.Entity); err == nil {
		// Scenario 2
		filter := model.DeviceConfigurationKeyValueDescriptionDataType{
			KeyName: util.Ptr(model.DeviceConfigurationKeyNameTypeCommunicationsStandard),
		}
		if dc.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateCommunicationStandard)
		}

		// Scenario 3
		filter.KeyName = util.Ptr(model.DeviceConfigurationKeyNameTypeAsymmetricChargingSupported)
		if dc.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateAsymmetricChargingSupport)
		}
	}
}

// the operating state of an EV was updated
func (e *EVCC) evOperatingStateDataUpdate(payload spineapi.EventPayload) {
	if deviceDiagnosis, err := client.NewDeviceDiagnosis(e.LocalEntity, payload.Entity); err == nil {
		if _, err := deviceDiagnosis.GetState(); err == nil && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateIdentifications)
		}
	}
}

// the identification data of an EV was updated
func (e *EVCC) evIdentificationDataUpdate(payload spineapi.EventPayload) {
	if evIdentification, err := client.NewIdentification(e.LocalEntity, payload.Entity); err == nil {
		// Scenario 4
		if evIdentification.CheckEventPayloadDataForFilter(payload.Data) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateIdentifications)
		}
	}
}

// the manufacturer data of an EV was updated
func (e *EVCC) evManufacturerDataUpdate(payload spineapi.EventPayload) {
	if evDeviceClassification, err := client.NewDeviceClassification(e.LocalEntity, payload.Entity); err == nil {
		// Scenario 5
		if _, err := evDeviceClassification.GetManufacturerDetails(); err == nil && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateManufacturerData)
		}
	}
}

// the electrical connection parameter description data of an EV was updated
func (e *EVCC) evElectricalParamerDescriptionUpdate(entity spineapi.EntityRemoteInterface) {
	if evElectricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity); err == nil {
		if _, err := evElectricalConnection.RequestPermittedValueSets(nil, nil); err != nil {
			logging.Log().Error("Error getting electrical permitted values:", err)
		}
	}
}

// the electrical connection permitted value sets data of an EV was updated
func (e *EVCC) evElectricalPermittedValuesUpdate(payload spineapi.EventPayload) {
	if evElectricalConnection, err := client.NewElectricalConnection(e.LocalEntity, payload.Entity); err == nil {
		filter := model.ElectricalConnectionParameterDescriptionDataType{
			ScopeType: util.Ptr(model.ScopeTypeTypeACPowerTotal),
		}
		if evElectricalConnection.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			// Scenario 6
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdateCurrentLimits)
		}
	}
}
