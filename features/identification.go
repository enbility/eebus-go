package features

import (
	"github.com/enbility/eebus-go/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/spine"
)

type Identification struct {
	*Feature
}

// Get a new Identification features helper
//
// - The feature on the local entity has to be of role client
// - The feature on the remote entity has to be of role server
func NewIdentification(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface) (*Identification, error) {
	feature, err := NewFeature(model.FeatureTypeTypeIdentification, localEntity, remoteEntity)
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
		return nil, api.ErrDataNotAvailable
	}

	return data.IdentificationData, nil
}
