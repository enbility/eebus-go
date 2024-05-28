package internal

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type SmartEnergyManagementPsCommon struct {
	featureLocal  spineapi.FeatureLocalInterface
	featureRemote spineapi.FeatureRemoteInterface
}

func NewLocalSmartEnergyManagementPs(featureLocal spineapi.FeatureLocalInterface) *SmartEnergyManagementPsCommon {
	return &SmartEnergyManagementPsCommon{
		featureLocal: featureLocal,
	}
}

func NewRemoteSmartEnergyManagementPs(featureRemote spineapi.FeatureRemoteInterface) *SmartEnergyManagementPsCommon {
	return &SmartEnergyManagementPsCommon{
		featureRemote: featureRemote,
	}
}

var _ api.FeatureServerInterface = (*SmartEnergyManagementPsCommon)(nil)

// return current data for FunctionTypeSmartEnergyManagementPsData
func (i *SmartEnergyManagementPsCommon) GetData() (*model.SmartEnergyManagementPsDataType, error) {
	function := model.FunctionTypeSmartEnergyManagementPsData

	data, err := featureDataCopyOfType[model.SmartEnergyManagementPsDataType](i.featureLocal, i.featureRemote, function)
	if err != nil || data == nil {
		return nil, api.ErrDataNotAvailable
	}

	return data, nil
}
