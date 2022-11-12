package features

import (
	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
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
func (i *IncentiveTable) RequestDescription() error {
	_, err := i.requestData(model.FunctionTypeIncentiveTableDescriptionData, nil, nil)
	if err != nil {
		logging.Log.Error(err)
		return err
	}

	return nil
}

// request FunctionTypeIncentiveTableConstraintsData from a remote entity
func (i *IncentiveTable) RequestConstraints() error {
	_, err := i.requestData(model.FunctionTypeIncentiveTableConstraintsData, nil, nil)
	if err != nil {
		logging.Log.Error(err)
		return err
	}

	return nil
}

// request FunctionTypeIncentiveTableData from a remote entity
func (i *IncentiveTable) RequestValues() error {
	_, err := i.requestData(model.FunctionTypeIncentiveTableData, nil, nil)
	if err != nil {
		logging.Log.Error(err)
		return err
	}

	return nil
}
