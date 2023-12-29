package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type Identification struct {
	*FeatureImpl
}

func NewIdentification(localRole, remoteRole model.RoleType, localEntity *spine.EntityLocalImpl, remoteEntity *spine.EntityRemoteImpl) (*Identification, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeIdentification, localRole, remoteRole, localEntity, remoteEntity)
	if err != nil {
		return nil, err
	}

	i := &Identification{
		FeatureImpl: feature,
	}

	return i, nil
}

// request FunctionTypeIdentificationListData from a remote entity
func (i *Identification) RequestValues() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeIdentificationListData, nil, nil)
}

// return current values for Identification
func (i *Identification) GetValues() ([]model.IdentificationDataType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeIdentificationListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.IdentificationListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.IdentificationData, nil
}
