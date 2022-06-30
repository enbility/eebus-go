package model

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/util"
)

var _ Updater[ElectricalConnectionPermittedValueSetListDataType] = (*ElectricalConnectionPermittedValueSetListDataType)(nil)
var _ util.HashKeyer = (*ElectricalConnectionPermittedValueSetDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) Update(s *ElectricalConnectionPermittedValueSetListDataType, filterPartial *FilterType, filterDelete *FilterType) {
	if s == nil {
		return
	}

	// TODO: consider filterPartial and filterDelete
	// TODO: consider items without identifiers
	// TODO: Check if only single fields should be considered here
	newList := util.Merge(r.ElectricalConnectionPermittedValueSetData, s.ElectricalConnectionPermittedValueSetData)

	r.ElectricalConnectionPermittedValueSetData = newList
}

func (r ElectricalConnectionPermittedValueSetDataType) HashKey() string {
	return fmt.Sprintf("%d|%d", *r.ElectricalConnectionId, *r.ParameterId)
}
