package api

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// Entity event callback
//
// Used by Use Case implementations
type EntityEventCallback func(ski string, device spineapi.DeviceRemoteInterface, entity spineapi.EntityRemoteInterface, event EventType)

type UseCaseBaseInterface interface {
	// add the use case
	AddUseCase()

	// remove the use case
	RemoveUseCase()

	// update availability of the use case
	//
	// NOTE: only allowed to be used for client side implementations
	// of a use case! Otherwise use `RemoveUseCase` and `AddUseCase`.
	UpdateUseCaseAvailability(available bool)

	// check if the entity is compatible with the use case
	IsCompatibleEntityType(entity spineapi.EntityRemoteInterface) bool

	// get the supported use case scenarios of the remote entity
	SupportedUseCaseScenarios(
		entity spineapi.EntityRemoteInterface,
	) []model.UseCaseScenarioSupportType

	// check if the entity supports all provided use case scenarios
	HasSupportForUseCaseScenarios(
		entity spineapi.EntityRemoteInterface,
		scenarios []model.UseCaseScenarioSupportType,
	) bool

	// return the current list of compatible remote entities and their scenarios
	RemoteEntities() []RemoteEntityScenarios
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
