package mot

import "github.com/enbility/eebus-go/api"

const (
	// Update of the list of remote entities supporting the Use Case
	//
	// Use `RemoteEntities` to get the current data
	UseCaseSupportUpdate api.EventType = "ma-mot-UseCaseSupportUpdate"

	// Outdoor temperature data updated
	//
	// Use `Temperature` to get the current data
	//
	// Use Case MOT, Scenario 1
	DataUpdateTemperature api.EventType = "ma-mot-DataUpdateTemperature"
)
