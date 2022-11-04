package features

import (
	"fmt"

	"github.com/DerAndereAndi/eebus-go/service"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type IncentiveTable struct {
	*FeatureImpl
}

func NewIncentiveTable(service *service.EEBUSService, entity *spine.EntityRemoteImpl) (*IncentiveTable, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeIncentiveTable, service, entity)
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
	_, err := i.requestData(model.FunctionTypeIncentiveTableDescriptionData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// request FunctionTypeIncentiveTableConstraintsData from a remote entity
func (i *IncentiveTable) RequestConstraints() error {
	_, err := i.requestData(model.FunctionTypeIncentiveTableConstraintsData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// request FunctionTypeIncentiveTableData from a remote entity
func (i *IncentiveTable) RequestValues() error {
	_, err := i.requestData(model.FunctionTypeIncentiveTableData)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}
