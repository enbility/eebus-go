package evcc

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/spine-go/model"
)

// value if the UCEVCC communication standard is unknown
const (
	EVCCCommunicationStandardUnknown model.DeviceConfigurationKeyValueStringType = "unknown"
)

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-evcc-UseCaseSupportUpdate"

	// An EV was connected
	//
	// Use Case EVCC, Scenario 1
	EvConnected api.EventType = "cem-evcc-EvConnected"

	// An EV was disconnected
	//
	// Note: The ev entity is no longer connected to the device!
	//
	// Use Case EVCC, Scenario 8
	EvDisconnected api.EventType = "cem-evcc-EvDisconnected"

	// EV charge state data was updated
	//
	// Use `ChargeState` to get the current data
	DataUpdateChargeState api.EventType = "cem-evcc-DataUpdateChargeState"

	// EV communication standard data was updated
	//
	// Use `CommunicationStandard` to get the current data
	//
	// Use Case EVCC, Scenario 2
	DataUpdateCommunicationStandard api.EventType = "cem-evcc-DataUpdateCommunicationStandard"

	// EV asymmetric charging data was updated
	//
	// Use `AsymmetricChargingSupport` to get the current data
	DataUpdateAsymmetricChargingSupport api.EventType = "cem-evcc-DataUpdateAsymmetricChargingSupport"

	// EV identificationdata was updated
	//
	// Use `Identifications` to get the current data
	//
	// Use Case EVCC, Scenario 4
	DataUpdateIdentifications api.EventType = "cem-evcc-DataUpdateIdentifications"

	// EV manufacturer data was updated
	//
	// Use `ManufacturerData` to get the current data
	//
	// Use Case EVCC, Scenario 5
	DataUpdateManufacturerData api.EventType = "cem-evcc-DataUpdateManufacturerData"

	// EV charging power limits
	//
	// Use `ChargingPowerLimits` to get the current data
	//
	// Use Case EVCC, Scenario 6
	DataUpdateCurrentLimits api.EventType = "cem-evcc-DataUpdateCurrentLimits"

	// EV permitted power limits updated
	//
	// Use `IsInSleepMode` to get the current data
	//
	// Use Case EVCC, Scenario 7
	DataUpdateIsInSleepMode api.EventType = "cem-evcc-DataUpdateIsInSleepMode"
)
