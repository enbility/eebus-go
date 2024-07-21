package evsecc

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remotei entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "cem-evsecc-UseCaseSupportUpdate"

	// An EVSE was connected
	EvseConnected api.EventType = "cem-evsecc-EvseConnected"

	// An EVSE was disconnected
	EvseDisconnected api.EventType = "cem-evsecc-EvseDisconnected"

	// EVSE manufacturer data was updated
	//
	// Use `ManufacturerData` to get the current data
	//
	// Use Case EVSECC, Scenario 1
	//
	// The entity of the message is the entity of the EVSE
	DataUpdateManufacturerData api.EventType = "cem-evsecc-DataUpdateManufacturerData"

	// EVSE operation state was updated
	//
	// Use `OperatingState` to get the current data
	//
	// Use Case EVSECC, Scenario 2
	//
	// The entity of the message is the entity of the EVSE
	DataUpdateOperatingState api.EventType = "cem-evsecc-DataUpdateOperatingState"
)
