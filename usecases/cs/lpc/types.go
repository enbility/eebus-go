package lpc

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cs-lpc-UseCaseSupportUpdate"

	// Load control obligation limit data update received
	//
	// Use `ConsumptionLimit` to get the current data
	//
	// Use Case LPC, Scenario 1
	DataUpdateLimit api.EventType = "cs-lpc-DataUpdateLimit"

	// An incoming load control obligation limit needs to be approved or denied
	//
	// Use `PendingConsumptionLimits` to get the currently pending consumption limits
	// awaiting approval and invoke `ApproveOrDenyConsumptionLimit` for each
	//
	// Use Case LPC, Scenario 1
	LimitWriteApprovalRequired api.EventType = "cs-lpc-LimitWriteApprovalRequired"

	// An incoming device configuration write needs to be approved or denied
	//
	// Use `PendingDeviceConfigurations` to get the currently pending device configurations
	// awaiting approval and invoke `ApproveOrDenyDeviceConfiguration` for each
	//
	// Use Case LPC, Scenario 1
	ConfigurationWriteApprovalRequired api.EventType = "cs-lpc-ConfigurationWriteApprovalRequired"

	// Failsafe limit for the consumed active (real) power of the
	// Controllable System data update received
	//
	// Use `FailsafeConsumptionActivePowerLimit` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeConsumptionActivePowerLimit api.EventType = "cs-lpc-DataUpdateFailsafeConsumptionActivePowerLimit"

	// Minimum time the Controllable System remains in "failsafe state" unless conditions
	// specified in this Use Case permit leaving the "failsafe state" data update received
	//
	// Use `FailsafeDurationMinimum` to get the current data
	//
	// Use Case LPC, Scenario 2
	DataUpdateFailsafeDurationMinimum api.EventType = "cs-lpc-DataUpdateFailsafeDurationMinimum"

	// Indicates a notify heartbeat event the application should care of.
	// E.g. going into or out of the Failsafe state
	//
	// Use Case LPC, Scenario 3
	DataUpdateHeartbeat api.EventType = "cs-lpc-DataUpdateHeartbeat"
)
