package features

import (
	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type IdentificationType struct {
	Identifier string
	Type       model.IdentificationTypeType
}

type Identification struct {
	*FeatureImpl
}

func NewIdentification(localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*Identification, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeIdentification, localRole, remoteRole, spineLocalDevice, entity)
	if err != nil {
		return nil, err
	}

	i := &Identification{
		FeatureImpl: feature,
	}

	return i, nil
}

// request FunctionTypeIdentificationListData from a remote entity
func (i *Identification) Request() (*model.MsgCounterType, error) {
	msgCounter, err := i.requestData(model.FunctionTypeIdentificationListData)
	if err != nil {
		logging.Log.Error(err)
		return nil, err
	}

	return msgCounter, nil
}

// return current values for Identification
func (i *Identification) GetValues() ([]IdentificationType, error) {
	if i.featureRemote == nil {
		return nil, ErrDataNotAvailable
	}

	rData := i.featureRemote.Data(model.FunctionTypeIdentificationListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.IdentificationListDataType)
	var resultSet []IdentificationType

	for _, item := range data.IdentificationData {
		if item.IdentificationValue == nil {
			continue
		}

		result := IdentificationType{
			Identifier: string(*item.IdentificationValue),
		}
		if item.IdentificationType != nil {
			result.Type = *item.IdentificationType
		}

		resultSet = append(resultSet, result)
	}

	return resultSet, nil
}
