package spine

import (
	"github.com/enbility/eebus-go/spine/model"
)

// manages the supported usecases for a device
// each device has its own UseCaseManager
type UseCaseManager struct {
	useCaseInformationMap map[model.UseCaseActorType][]model.UseCaseSupportType

	localDevice *DeviceImpl
}

// return a new UseCaseManager
func NewUseCaseManager(localDevice *DeviceImpl) *UseCaseManager {
	return &UseCaseManager{
		useCaseInformationMap: make(map[model.UseCaseActorType][]model.UseCaseSupportType),
		localDevice:           localDevice,
	}
}

// add a usecase
func (r *UseCaseManager) Add(
	actor model.UseCaseActorType,
	useCaseName model.UseCaseNameType,
	useCaseVersion model.SpecificationVersionType,
	useCaseDocumemtSubRevision string,
	useCaseAvailable bool,
	scenarios []model.UseCaseScenarioSupportType,
) {
	useCaseSupport := model.UseCaseSupportType{
		UseCaseVersion:   &useCaseVersion,
		UseCaseName:      &useCaseName,
		UseCaseAvailable: &useCaseAvailable,
		ScenarioSupport:  scenarios,
	}

	if len(useCaseDocumemtSubRevision) > 0 {
		useCaseSupport.UseCaseDocumentSubRevision = &useCaseDocumemtSubRevision
	}

	useCaseInfo, exists := r.useCaseInformationMap[actor]
	if !exists {
		useCaseInfo = make([]model.UseCaseSupportType, 0)
	}
	useCaseInfo = append(useCaseInfo, useCaseSupport)

	r.useCaseInformationMap[actor] = useCaseInfo
}

// this needs to be called when a new notification or reply will provide a new set of UseCases
func (r *UseCaseManager) RemoveAll() {
	r.useCaseInformationMap = make(map[model.UseCaseActorType][]model.UseCaseSupportType)
}

// return all actors and their supported usecases
func (r *UseCaseManager) UseCaseInformation() []model.UseCaseInformationDataType {
	var result []model.UseCaseInformationDataType

	for actor, useCaseSupport := range r.useCaseInformationMap {
		thisActor := actor
		// according to ProtocolSpecification Version 1.3.0 chapter 7.5.2
		// the address is mandatory. At least the device address should be shown,
		// preferably the entity and features as well
		deviceAddress := &model.FeatureAddressType{
			Device: r.localDevice.Address(),
		}
		useCaseInfo := model.UseCaseInformationDataType{
			Address:        deviceAddress,
			Actor:          &thisActor,
			UseCaseSupport: useCaseSupport,
		}
		result = append(result, useCaseInfo)
	}

	return result
}
