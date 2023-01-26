package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

type HVAC struct {
	*FeatureImpl
}

func NewHVAC(localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*HVAC, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeHvac, localRole, remoteRole, spineLocalDevice, entity)
	if err != nil {
		return nil, err
	}

	i := &HVAC{
		FeatureImpl: feature,
	}

	return i, nil
}

// request FunctionTypeHvacOverrunDescriptionListData from a remote entity
func (i *HVAC) RequestOverrunDescriptions() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeHvacOverrunDescriptionListData, nil, nil)
}

// return current values for Hvac Overrun Descriptions
func (i *HVAC) GetOverrunDescriptions() ([]model.HvacOverrunDescriptionDataType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeHvacOverrunDescriptionListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.HvacOverrunDescriptionListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.HvacOverrunDescriptionData, nil
}

// request FunctionTypeHvacOverrunListData from a remote entity
func (i *HVAC) RequestOverrunValues() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeHvacOverrunListData, nil, nil)
}

// return current values for HvacOverrun
func (i *HVAC) GetValues() ([]model.HvacOverrunDataType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeHvacOverrunListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.HvacOverrunListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.HvacOverrunData, nil
}

// request FunctionTypeHvacSystemFunctionDescriptionListData from a remote entity
func (i *HVAC) RequestSystemFunctionDescriptions() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeHvacSystemFunctionDescriptionListData, nil, nil)
}

// return current values for Hvac System Function Descriptions
func (i *HVAC) GetSystemFunctionDescriptions() ([]model.HvacSystemFunctionDescriptionDataType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeHvacSystemFunctionDescriptionListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.HvacSystemFunctionDescriptionListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.HvacSystemFunctionDescriptionData, nil
}

// request FunctionTypeHvacSystemFunctionListData from a remote entity
func (i *HVAC) RequestSystemFunctionValues() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeHvacSystemFunctionListData, nil, nil)
}

// return current values for Hvac System Function Values
func (i *HVAC) GetSystemFunctionValues() ([]model.HvacSystemFunctionDataType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeHvacSystemFunctionListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.HvacSystemFunctionListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.HvacSystemFunctionData, nil
}

// request FunctionTypeHvacSystemFunctionOperationModeRelationListData from a remote entity
func (i *HVAC) RequestSystemFunctionOperationModeRelationValues() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeHvacSystemFunctionOperationModeRelationListData, nil, nil)
}

// return current values for Hvac System Function Operation Mode Relation Values
func (i *HVAC) GetSystemFunctionOperationModeRelationValues() ([]model.HvacSystemFunctionOperationModeRelationDataType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeHvacSystemFunctionOperationModeRelationListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.HvacSystemFunctionOperationModeRelationListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.HvacSystemFunctionOperationModeRelationData, nil
}

// request FunctionTypeHvacSystemFunctionPowerSequenceRelationListData from a remote entity
func (i *HVAC) RequestSystemFunctionPowerSequenceRelationValues() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeHvacSystemFunctionPowerSequenceRelationListData, nil, nil)
}

// return current values for Hvac System Function Power Sequence Relation Values
func (i *HVAC) GetSystemFunctionPowerSequenceRelationValues() ([]model.HvacSystemFunctionPowerSequenceRelationDataType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeHvacSystemFunctionPowerSequenceRelationListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.HvacSystemFunctionPowerSequenceRelationListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.HvacSystemFunctionPowerSequenceRelationData, nil
}

// request FunctionTypeHvacSystemFunctionSetPointRelationListData from a remote entity
func (i *HVAC) RequestSystemFunctionSetpointRelationValues() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeHvacSystemFunctionSetPointRelationListData, nil, nil)
}

// return current values for Hvac System Function Power Sequence Relation Values
func (i *HVAC) GetSystemFunctionSetpointRelationValues() ([]model.HvacSystemFunctionSetpointRelationDataType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeHvacSystemFunctionSetPointRelationListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.HvacSystemFunctionSetpointRelationListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.HvacSystemFunctionSetpointRelationData, nil
}

// request FunctionTypeHvacOperationModeDescriptionListData from a remote entity
func (i *HVAC) RequestOperationModeDescriptions() (*model.MsgCounterType, error) {
	return i.requestData(model.FunctionTypeHvacOperationModeDescriptionListData, nil, nil)
}

// return current values for Hvac System Function Power Sequence Relation Values
func (i *HVAC) GetOperationModeDescriptions() ([]model.HvacOperationModeDescriptionDataType, error) {
	rData := i.featureRemote.Data(model.FunctionTypeHvacOperationModeDescriptionListData)
	if rData == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData.(*model.HvacOperationModeDescriptionListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	return data.HvacOperationModeDescriptionData, nil
}
