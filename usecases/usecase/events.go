package usecase

import (
	"slices"

	"github.com/enbility/eebus-go/usecases/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// handle SPINE events
func (u *UseCaseBase) HandleEvent(payload spineapi.EventPayload) {
	if u.deviceOrEntityRemoved(payload) {
		return
	}

	switch payload.Data.(type) {
	case *model.NodeManagementUseCaseDataType,
		*model.NodeManagementDetailedDiscoveryDataType:
		u.useCaseDataUpdate(payload)
	default:
		return
	}
}

func (u *UseCaseBase) deviceOrEntityRemoved(payload spineapi.EventPayload) bool {
	if payload.EventType == spineapi.EventTypeDeviceChange &&
		payload.ChangeType == spineapi.ElementChangeRemove &&
		payload.Device != nil &&
		payload.Entity == nil {
		// device was disconnected, remove all usecases related to this device
		u.removeDeviceFromAvailableEntityScenarios(payload.Device)
		return true
	}

	if internal.IsEntityRemoved(payload) {
		// entity was disconnected, remove all usecases related to this entity
		u.removeEntityFromAvailableEntityScenarios(payload.Entity)
		return true
	}

	return false
}

func (u *UseCaseBase) useCaseDataUpdate(
	payload spineapi.EventPayload,
) {
	remoteDevice := payload.Device

	// go over the use cases and check which entity of the remote device supports the usecase
	ucs := remoteDevice.UseCases()
	for _, uc := range ucs {
		if uc.Actor == nil || !slices.Contains(u.validActorTypes, *uc.Actor) {
			continue
		}

		for _, support := range uc.UseCaseSupport {
			// UseCaseAvailable should only be checkd for client functionality
			// of a use acse. But as there are devices (e.g. Porsche PMCC) that
			// also use it for the server side, we need to check it always.
			if support.UseCaseName == nil ||
				*support.UseCaseName != u.UseCaseName ||
				support.ScenarioSupport == nil &&
					len(support.ScenarioSupport) == 0 ||
				(support.UseCaseAvailable != nil && !*support.UseCaseAvailable) {
				continue
			}

			// TODO: check the version and compare them for compatibility, it may be required to provide a minimmal compatible version
			// but then also the question is what to do if the version is newer

			entitiesToCheck := []spineapi.EntityRemoteInterface{}

			// if the address is given, use that for further checks
			if uc.Address != nil {
				ucEntity := remoteDevice.Entity(uc.Address.Entity)
				// the PMCP EVSE reports EV use cases with the address of the EVSE
				if *uc.Actor == model.UseCaseActorTypeEV && len(uc.Address.Entity) == 1 {
					// add the EV subentity to the address
					evAddress := append(uc.Address.Entity, 1)
					ucEntity = remoteDevice.Entity(evAddress)
				}
				if ucEntity != nil {
					entitiesToCheck = append(entitiesToCheck, ucEntity)
				}
			}

			// if no entity was defined/found, go over all known entities and search for one or multiple
			if len(entitiesToCheck) == 0 {
				entitiesToCheck = remoteDevice.Entities()
			}

			for _, entity := range entitiesToCheck {
				if !slices.Contains(u.validEntityTypes, entity.EntityType()) {
					continue
				}

				supportedScenarios := []model.UseCaseScenarioSupportType{}

				// go over each scenario this use case supports and check if the required server features are available
				for _, scenario := range u.useCaseScenarios {
					if !slices.Contains(support.ScenarioSupport, scenario.Scenario) {
						continue
					}

					// check if the required server features are available
					requiredServerFeatures := u.requiredServerFeaturesForScenario(scenario.Scenario)
					foundMatchingServerFeatures := []model.FeatureTypeType{}
					for _, feature := range entity.Features() {
						if feature.Role() != model.RoleTypeServer ||
							!slices.Contains(requiredServerFeatures, feature.Type()) {
							continue
						}

						foundMatchingServerFeatures = append(foundMatchingServerFeatures, feature.Type())
					}

					// check if the minimum required server features are available
					slices.Sort(foundMatchingServerFeatures)
					slices.Sort(requiredServerFeatures)

					if slices.Compare(requiredServerFeatures, foundMatchingServerFeatures) != 0 {
						continue
					}

					supportedScenarios = append(supportedScenarios, scenario.Scenario)
				}

				u.updateRemoteEntityScenarios(entity, supportedScenarios)
			}
		}
	}
}
