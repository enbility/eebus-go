package api

import (
	spineapi "github.com/enbility/spine-go/api"
)

// Entity event callback
//
// Used by Use Case implementations
type EntityEventCallback func(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event EventType)

type UseCaseBaseInterface interface {
	// add the use case
	AddUseCase()

	// update availability of the use case
	UpdateUseCaseAvailability(available bool)

	// check if the entity is compatible with the use case
	IsCompatibleEntity(entity spineapi.EntityRemoteInterface) bool
}

// Implemented by each Use Case
type UseCaseInterface interface {
	UseCaseBaseInterface

	// add the features
	AddFeatures()

	// returns if the entity supports the usecase
	//
	// possible errors:
	//   - ErrDataNotAvailable if that information is not (yet) available
	//   - and others
	IsUseCaseSupported(remoteEntity spineapi.EntityRemoteInterface) (bool, error)
}

type ManufacturerData struct {
	DeviceName                     string `json:"deviceName,omitempty"`
	DeviceCode                     string `json:"deviceCode,omitempty"`
	SerialNumber                   string `json:"serialNumber,omitempty"`
	SoftwareRevision               string `json:"softwareRevision,omitempty"`
	HardwareRevision               string `json:"hardwareRevision,omitempty"`
	VendorName                     string `json:"vendorName,omitempty"`
	VendorCode                     string `json:"vendorCode,omitempty"`
	BrandName                      string `json:"brandName,omitempty"`
	PowerSource                    string `json:"powerSource,omitempty"`
	ManufacturerNodeIdentification string `json:"manufacturerNodeIdentification,omitempty"`
	ManufacturerLabel              string `json:"manufacturerLabel,omitempty"`
	ManufacturerDescription        string `json:"manufacturerDescription,omitempty"`
}
