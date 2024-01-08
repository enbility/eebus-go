package model

import "sync"

var uciMux sync.Mutex

// UseCaseInformationDataType

// find the matching UseCaseSupport index for a UseCaseNameType
func (u *UseCaseInformationDataType) useCaseSupportIndex(useCaseName UseCaseNameType) (int, bool) {
	// get the element with the same entity
	for index, item := range u.UseCaseSupport {
		if item.UseCaseName != nil && *item.UseCaseName == useCaseName {
			return index, true
		}
	}

	return -1, false
}

// add a new UseCaseSupportType
func (u *UseCaseInformationDataType) Add(useCase UseCaseSupportType) {
	uciMux.Lock()
	defer uciMux.Unlock()

	if useCase.UseCaseName == nil {
		return
	}

	// only add it if it does not exist yet
	if _, ok := u.useCaseSupportIndex(*useCase.UseCaseName); ok {
		return
	}

	u.UseCaseSupport = append(u.UseCaseSupport, useCase)
}

// remove a UseCaseSupportType with a given UseCaseNameType
func (u *UseCaseInformationDataType) Remove(useCaseName UseCaseNameType) {
	uciMux.Lock()
	defer uciMux.Unlock()

	var usecases []UseCaseSupportType

	for _, item := range u.UseCaseSupport {
		if item.UseCaseName != nil && *item.UseCaseName != useCaseName {
			continue
		}

		usecases = append(usecases, item)
	}

	u.UseCaseSupport = usecases
}
