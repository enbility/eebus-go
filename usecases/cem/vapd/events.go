package vapd

import (
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/usecases/internal"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// handle SPINE events
func (e *VAPD) HandleEvent(payload spineapi.EventPayload) {
	// only about events from an SGMW entity or device changes for this remote device

	if !e.IsCompatibleEntityType(payload.Entity) {
		return
	}

	if internal.IsEntityAdded(payload) {
		e.inverterConnected(payload.Entity)
		return
	}

	if payload.EventType != spineapi.EventTypeDataChange ||
		payload.ChangeType != spineapi.ElementChangeUpdate {
		return
	}

	switch payload.Data.(type) {
	case *model.DeviceConfigurationKeyValueDescriptionListDataType:
		e.inverterConfigurationDescriptionDataUpdate(payload.Entity)

	case *model.DeviceConfigurationKeyValueListDataType:
		e.inverterConfigurationDataUpdate(payload)

	case *model.MeasurementDescriptionListDataType:
		e.inverterMeasurementDescriptionDataUpdate(payload.Entity)

	case *model.MeasurementListDataType:
		e.inverterMeasurementDataUpdate(payload)
	}
}

// process required steps when a grid device is connected
func (e *VAPD) inverterConnected(entity spineapi.EntityRemoteInterface) {
	if deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity); err == nil {
		if !deviceConfiguration.HasSubscription() {
			if _, err := deviceConfiguration.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		// get configuration data
		if _, err := deviceConfiguration.RequestKeyValueDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}

	if electricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity); err == nil {
		if !electricalConnection.HasSubscription() {
			if _, err := electricalConnection.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		// get electrical connection parameter
		if _, err := electricalConnection.RequestDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := electricalConnection.RequestParameterDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}

	if measurement, err := client.NewMeasurement(e.LocalEntity, entity); err == nil {
		if !measurement.HasSubscription() {
			if _, err := measurement.Subscribe(); err != nil {
				logging.Log().Error(err)
			}
		}

		// get measurement parameters
		if _, err := measurement.RequestDescriptions(nil, nil); err != nil {
			logging.Log().Error(err)
		}

		if _, err := measurement.RequestConstraints(nil, nil); err != nil {
			logging.Log().Error(err)
		}
	}
}

// the configuration key description data of an SMGW was updated
func (e *VAPD) inverterConfigurationDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, entity); err == nil {
		// key value descriptions received, now get the data
		if _, err := deviceConfiguration.RequestKeyValues(nil, nil); err != nil {
			logging.Log().Error("Error getting configuration key values:", err)
		}
	}
}

// the measurement data of an SMGW was updated
func (e *VAPD) inverterConfigurationDataUpdate(payload spineapi.EventPayload) {
	// Scenario 1
	if deviceConfiguration, err := client.NewDeviceConfiguration(e.LocalEntity, payload.Entity); err == nil {
		filter := model.DeviceConfigurationKeyValueDescriptionDataType{
			KeyName:   util.Ptr(model.DeviceConfigurationKeyNameTypePeakPowerOfPVSystem),
			ValueType: util.Ptr(model.DeviceConfigurationKeyValueTypeTypeScaledNumber),
		}
		if deviceConfiguration.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePowerNominalPeak)
		}
	}
}

// the measurement descriptiondata of an SMGW was updated
func (e *VAPD) inverterMeasurementDescriptionDataUpdate(entity spineapi.EntityRemoteInterface) {
	if measurement, err := client.NewMeasurement(e.LocalEntity, entity); err == nil {
		// measurement descriptions received, now get the data
		if _, err := measurement.RequestData(nil, nil); err != nil {
			logging.Log().Error("Error getting measurement list values:", err)
		}
	}
}

// the measurement data of an SMGW was updated
func (e *VAPD) inverterMeasurementDataUpdate(payload spineapi.EventPayload) {
	if measurement, err := client.NewMeasurement(e.LocalEntity, payload.Entity); err == nil {
		// Scenario 2
		filter := model.MeasurementDescriptionDataType{
			ScopeType: util.Ptr(model.ScopeTypeTypeACPowerTotal),
		}
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePower)
		}

		// Scenario 3
		filter.ScopeType = util.Ptr(model.ScopeTypeTypeACYieldTotal)
		if measurement.CheckEventPayloadDataForFilter(payload.Data, filter) && e.EventCB != nil {
			e.EventCB(payload.Ski, payload.Device, payload.Entity, DataUpdatePVYieldTotal)
		}
	}
}
