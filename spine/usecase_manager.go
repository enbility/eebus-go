package spine

import "github.com/enbility/eebus-go/spine/model"

// manages the supported usecases for a device
// each device has its own UseCaseManager
type UseCaseManager struct {
	useCaseInformationMap map[model.UseCaseActorType][]model.UseCaseSupportType
}

// return a new UseCaseManager
func NewUseCaseManager() *UseCaseManager {
	return &UseCaseManager{
		useCaseInformationMap: make(map[model.UseCaseActorType][]model.UseCaseSupportType),
	}
}

// add a usecase
func (r *UseCaseManager) Add(actor model.UseCaseActorType, useCaseName model.UseCaseNameType, useCaseVersion model.SpecificationVersionType, useCaseAvailable bool, scenarios []model.UseCaseScenarioSupportType) {
	useCaseSupport := model.UseCaseSupportType{
		UseCaseVersion:   &useCaseVersion,
		UseCaseName:      &useCaseName,
		UseCaseAvailable: &useCaseAvailable,
		ScenarioSupport:  scenarios,
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
		useCaseInfo := model.UseCaseInformationDataType{
			//Address:        r.address, // TODO: which address ???
			Actor:          &thisActor,
			UseCaseSupport: useCaseSupport,
		}
		result = append(result, useCaseInfo)
	}

	return result
}
