package evsoc

import "github.com/enbility/eebus-go/api"

const (
	// EV state of charge data was updated
	//
	// Use `StateOfCharge` to get the current data
	//
	// Use Case EVSOC, Scenario 1
	DataUpdateStateOfCharge api.EventType = "ucevsoc-DataUpdateStateOfCharge"
)
