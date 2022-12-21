package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type IncentiveTable struct {
	*FeatureImpl
}

func NewIncentiveTable(localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*IncentiveTable, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeIncentiveTable, localRole, remoteRole, spineLocalDevice, entity)
	if err != nil {
		return nil, err
	}

	i := &IncentiveTable{
		FeatureImpl: feature,
	}

	return i, nil
}

// request FunctionTypeIncentiveTableDescriptionData from a remote entity
func (i *IncentiveTable) RequestDescriptions() error {
	_, err := i.requestData(model.FunctionTypeIncentiveTableDescriptionData, nil, nil)
	return err
}

// request FunctionTypeIncentiveTableConstraintsData from a remote entity
func (i *IncentiveTable) RequestConstraints() error {
	_, err := i.requestData(model.FunctionTypeIncentiveTableConstraintsData, nil, nil)
	return err
}

// request FunctionTypeIncentiveTableData from a remote entity
func (i *IncentiveTable) RequestValues() error {
	_, err := i.requestData(model.FunctionTypeIncentiveTableData, nil, nil)
	return err
}
