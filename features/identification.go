package features

import (
	"github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type Identification struct {
	*Feature
}

func NewIdentification(
	localRole, remoteRole model.RoleType,
	localEntity api.EntityLocalInterface,
	remoteEntity api.EntityRemoteInterface) (*Identification, error) {
	feature, err := NewFeature(model.FeatureTypeTypeIdentification, localRole, remoteRole, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	i := &Identification{
		Feature: feature,
	}

	return i, nil
}

// request FunctionTypeIdentificationListData from a remote entity
func (i *Identification) RequestValues() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeIdentificationListData, nil, nil)
}

// return current values for Identification
func (i *Identification) GetValues() ([]model.IdentificationDataType, error) {
	data, err := spine.RemoteFeatureDataCopyOfType[*model.IdentificationListDataType](i.featureRemote, model.FunctionTypeIdentificationListData)
	if err != nil {
		return nil, ErrDataNotAvailable
	}

	return data.IdentificationData, nil
}
