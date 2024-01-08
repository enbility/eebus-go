package model

import (
	"reflect"
	"sync"
)

var nmMux sync.Mutex

// NodeManagementDestinationListDataType

var _ Updater = (*NodeManagementDestinationListDataType)(nil)

func (r *NodeManagementDestinationListDataType) UpdateList(newList any, filterPartial, filterDelete *FilterType) {
	var newData []NodeManagementDestinationDataType
	if newList != nil {
		newData = newList.(*NodeManagementDestinationListDataType).NodeManagementDestinationData
	}

	r.NodeManagementDestinationData = UpdateList(r.NodeManagementDestinationData, newData, filterPartial, filterDelete)
}

// NodeManagementUseCaseDataType

// find the matching UseCaseInformation index for
// a given FeatureAddressType, UseCaseActorType and UseCaseNameType
//
// if UseCaseActorType and UseCaseNameType are empty they are ignored,
// and the first matching UseCaseInformation item is returned
func (n *NodeManagementUseCaseDataType) useCaseInformationIndex(
	address FeatureAddressType,
	actor UseCaseActorType,
	useCaseName UseCaseNameType,
) (int, bool) {

	// get the element with the same entity
	for index, item := range n.UseCaseInformation {
		if item.Address.Device == nil ||
			item.Address.Entity == nil ||
			!reflect.DeepEqual(item.Address.Device, address.Device) ||
			!reflect.DeepEqual(item.Address.Entity, address.Entity) {
			continue
		}

		if len(actor) == 0 && len(useCaseName) == 0 {
			return index, true
		}

		if len(actor) > 0 {
			if item.Actor == nil || *item.Actor != actor {
				continue
			}

		}

		if len(useCaseName) == 0 {
			return index, true
		}

		if _, ok := item.useCaseSupportIndex(useCaseName); ok {
			return index, true
		}
	}

	return -1, false
}

// add a new UseCaseSupportType
func (n *NodeManagementUseCaseDataType) AddUseCaseSupport(
	address FeatureAddressType,
	actor UseCaseActorType,
	useCaseName UseCaseNameType,
	useCaseVersion SpecificationVersionType,
	useCaseDocumemtSubRevision string,
	useCaseAvailable bool,
	scenarios []UseCaseScenarioSupportType,
) {
	nmMux.Lock()
	defer nmMux.Unlock()

	useCaseSupport := UseCaseSupportType{
		UseCaseName:                &useCaseName,
		UseCaseVersion:             &useCaseVersion,
		UseCaseAvailable:           &useCaseAvailable,
		ScenarioSupport:            scenarios,
		UseCaseDocumentSubRevision: &useCaseDocumemtSubRevision,
	}

	// is there an entry for the entity address and actor
	usecaseIndex, ok := n.useCaseInformationIndex(address, actor, "")

	if ok {
		n.UseCaseInformation[usecaseIndex].Add(useCaseSupport)
	} else {
		// create a new element for this entity
		useCaseInformation := UseCaseInformationDataType{
			Address: &FeatureAddressType{
				Device: address.Device,
				Entity: address.Entity,
			},
			Actor:          &actor,
			UseCaseSupport: []UseCaseSupportType{useCaseSupport},
		}
		n.UseCaseInformation = append(n.UseCaseInformation, useCaseInformation)
	}
}

// Remove a UseCaseSupportType with
// a provided FeatureAddressType, UseCaseActorType and UseCaseNameType
func (n *NodeManagementUseCaseDataType) RemoveUseCaseSupport(
	address FeatureAddressType,
	actor UseCaseActorType,
	useCaseName UseCaseNameType,
) {
	nmMux.Lock()
	defer nmMux.Unlock()

	// is there an entry for the entity address, actor and usecase name
	usecaseIndex, ok := n.useCaseInformationIndex(address, actor, useCaseName)
	if !ok {
		return
	}

	var usecaseInfo []UseCaseInformationDataType

	for index, item := range n.UseCaseInformation {
		if index != usecaseIndex {
			usecaseInfo = append(usecaseInfo, item)
			continue
		}

		item.Remove(useCaseName)

		// only add the item if there are any usecases left
		if len(item.UseCaseSupport) == 0 {
			continue
		}

		usecaseInfo = append(usecaseInfo, item)
	}

	n.UseCaseInformation = usecaseInfo
}
