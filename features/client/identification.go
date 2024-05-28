package client

import (
	internal "github.com/enbility/eebus-go/features/internal"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

type Identification struct {
	*Feature

	*internal.IdentificationCommon
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
		Feature:              feature,
		IdentificationCommon: internal.NewRemoteIdentification(feature.featureRemote),
	}

	return i, nil
}

// request FunctionTypeIdentificationListData from a remote entity
func (i *Identification) RequestValues() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeIdentificationListData, nil, nil)
}
