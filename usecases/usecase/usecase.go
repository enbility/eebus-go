package usecase

import (
	"slices"

	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type UseCaseBase struct {
	LocalEntity spineapi.EntityLocalInterface

	UseCaseActor              model.UseCaseActorType
	UseCaseName               model.UseCaseNameType
	useCaseVersion            model.SpecificationVersionType
	useCaseDocumentSubVersion string
	useCaseScenarios          []model.UseCaseScenarioSupportType

	EventCB api.EntityEventCallback

	validEntityTypes []model.EntityTypeType
}

var _ api.UseCaseBaseInterface = (*UseCaseBase)(nil)

func NewUseCaseBase(
	localEntity spineapi.EntityLocalInterface,
	usecaseActor model.UseCaseActorType,
	usecaseName model.UseCaseNameType,
	useCaseVersion string,
	useCaseDocumentSubVersion string,
	useCaseScenarios []model.UseCaseScenarioSupportType,
	eventCB api.EntityEventCallback,
	validEntityTypes []model.EntityTypeType,
) *UseCaseBase {
	return &UseCaseBase{
		LocalEntity:               localEntity,
		UseCaseActor:              usecaseActor,
		UseCaseName:               usecaseName,
		useCaseVersion:            model.SpecificationVersionType(useCaseVersion),
		useCaseDocumentSubVersion: useCaseDocumentSubVersion,
		useCaseScenarios:          useCaseScenarios,
		EventCB:                   eventCB,
		validEntityTypes:          validEntityTypes,
	}
}

func (u *UseCaseBase) AddUseCase() {
	u.LocalEntity.AddUseCaseSupport(
		u.UseCaseActor,
		u.UseCaseName,
		u.useCaseVersion,
		u.useCaseDocumentSubVersion,
		true,
		u.useCaseScenarios)
}

func (u *UseCaseBase) UpdateUseCaseAvailability(available bool) {
	u.LocalEntity.SetUseCaseAvailability(u.UseCaseActor, u.UseCaseName, available)
}

func (u *UseCaseBase) IsCompatibleEntity(entity spineapi.EntityRemoteInterface) bool {
	if entity == nil {
		return false
	}

	return slices.Contains(u.validEntityTypes, entity.EntityType())
}
