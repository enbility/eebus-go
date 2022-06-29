package model

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/util"
)

var _ Updater[ElectricalConnectionPermittedValueSetListDataType] = (*ElectricalConnectionPermittedValueSetListDataType)(nil)

func (r *ElectricalConnectionPermittedValueSetListDataType) Update(s *ElectricalConnectionPermittedValueSetListDataType, filterPartial *FilterType, filterDelete *FilterType) {
	if s == nil {
		return
	}

	// TODO: consider filterPartial and filterDelete
	newList := util.Union(r.ElectricalConnectionPermittedValueSetData, s.ElectricalConnectionPermittedValueSetData, CalcHash)

	r.ElectricalConnectionPermittedValueSetData = newList
}

func CalcHash(s *ElectricalConnectionPermittedValueSetDataType) string {
	return fmt.Sprintf("%d|%d", *s.ElectricalConnectionId, *s.ParameterId)
}
