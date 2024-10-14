package api

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type CemOHPCFInterface interface {
	api.UseCaseInterface

	SmartEnergyManagementData(entity spineapi.EntityRemoteInterface) (
		smartEnergyManagementData model.SmartEnergyManagementPsDataType, resultErr error)
}
