package features

import (
	"github.com/enbility/eebus-go/spine"
	"github.com/enbility/eebus-go/spine/model"
)

// Details about the electrical connection
type ElectricalDescriptionType struct {
	ConnectionID            uint
	PowerSupplyType         model.ElectricalConnectionVoltageTypeType
	AcConnectedPhases       uint
	PositiveEnergyDirection model.EnergyDirectionType
}

// Details about the limits of an electrical connection
type ElectricalLimitType struct {
	ConnectionID uint
	Min          float64
	Max          float64
	Default      float64
	Phase        model.ElectricalConnectionPhaseNameType
	Scope        model.ScopeTypeType
}

type ElectricalConnection struct {
	*FeatureImpl
}

func NewElectricalConnection(localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*ElectricalConnection, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeElectricalConnection, localRole, remoteRole, spineLocalDevice, entity)
	if err != nil {
		return nil, err
	}

	e := &ElectricalConnection{
		FeatureImpl: feature,
	}

	return e, nil
}

// request ElectricalConnectionDescriptionListDataType from a remote entity
func (e *ElectricalConnection) RequestDescription() error {
	if _, err := e.requestData(model.FunctionTypeElectricalConnectionDescriptionListData, nil, nil); err != nil {
		return err
	}

	return nil
}

// request FunctionTypeElectricalConnectionParameterDescriptionListData from a remote entity
func (e *ElectricalConnection) RequestParameterDescription() error {
	if _, err := e.requestData(model.FunctionTypeElectricalConnectionParameterDescriptionListData, nil, nil); err != nil {
		return err
	}

	return nil
}

// request FunctionTypeElectricalConnectionPermittedValueSetListData from a remote entity
func (e *ElectricalConnection) RequestPermittedValueSet() (*model.MsgCounterType, error) {
	msgCounter, err := e.requestData(model.FunctionTypeElectricalConnectionPermittedValueSetListData, nil, nil)
	if err != nil {
		return nil, err
	}

	return msgCounter, nil
}

type electricalParamDescriptionMapMeasurementId map[model.MeasurementIdType]model.ElectricalConnectionParameterDescriptionDataType
type electricatlParamDescriptionMapParamId map[model.ElectricalConnectionParameterIdType]model.ElectricalConnectionParameterDescriptionDataType

// return a map of ElectricalConnectionParameterDescriptionListDataType with measurementId as key and
// ElectricalConnectionParameterDescriptionListDataType with parameterId as key
func (e *ElectricalConnection) GetParamDescriptionListData() (electricalParamDescriptionMapMeasurementId, electricatlParamDescriptionMapParamId, error) {
	if e.featureRemote == nil {
		return nil, nil, ErrDataNotAvailable
	}

	rData := e.featureRemote.Data(model.FunctionTypeElectricalConnectionParameterDescriptionListData)
	if rData == nil {
		return nil, nil, ErrDataNotAvailable
	}
	data := rData.(*model.ElectricalConnectionParameterDescriptionListDataType)
	if data == nil {
		return nil, nil, ErrDataNotAvailable
	}

	refMeasurement := make(electricalParamDescriptionMapMeasurementId)
	refElectrical := make(electricatlParamDescriptionMapParamId)
	for _, item := range data.ElectricalConnectionParameterDescriptionData {
		if item.MeasurementId == nil || item.ElectricalConnectionId == nil {
			continue
		}
		refMeasurement[*item.MeasurementId] = item
		refElectrical[*item.ParameterId] = item
	}

	return refMeasurement, refElectrical, nil
}

// return current values for Electrical Description
func (e *ElectricalConnection) GetDescription() ([]ElectricalDescriptionType, error) {
	if e.featureRemote == nil {
		return nil, ErrDataNotAvailable
	}

	rData := e.featureRemote.Data(model.FunctionTypeElectricalConnectionDescriptionListData)
	if rData == nil {
		return nil, ErrMetadataNotAvailable
	}
	data := rData.(*model.ElectricalConnectionDescriptionListDataType)
	if data == nil {
		return nil, ErrMetadataNotAvailable
	}

	var resultSet []ElectricalDescriptionType

	for _, item := range data.ElectricalConnectionDescriptionData {
		if item.ElectricalConnectionId == nil {
			continue
		}

		result := ElectricalDescriptionType{}

		if item.PowerSupplyType != nil {
			result.PowerSupplyType = *item.PowerSupplyType
		}
		if item.AcConnectedPhases != nil {
			result.AcConnectedPhases = *item.AcConnectedPhases
		}
		if item.PositiveEnergyDirection != nil {
			result.PositiveEnergyDirection = *item.PositiveEnergyDirection
		}

		resultSet = append(resultSet, result)
	}

	return resultSet, nil
}

// return number of phases the device is connected with
func (e *ElectricalConnection) GetConnectedPhases() (uint, error) {
	if e.featureRemote == nil {
		return 0, ErrDataNotAvailable
	}

	rData := e.featureRemote.Data(model.FunctionTypeElectricalConnectionDescriptionListData)
	if rData == nil {
		return 0, ErrDataNotAvailable
	}

	data := rData.(*model.ElectricalConnectionDescriptionListDataType)
	if data == nil {
		return 0, ErrDataNotAvailable
	}

	for _, item := range data.ElectricalConnectionDescriptionData {
		if item.ElectricalConnectionId == nil {
			continue
		}

		if item.AcConnectedPhases != nil {
			return *item.AcConnectedPhases, nil
		}
	}

	// default to 3 if the value is not available
	return 3, nil
}

// return current current limit values
//
// returns a map with the phase ("a", "b", "c") as a key for
// minimum, maximum, default/pause values
func (e *ElectricalConnection) GetCurrentsLimits() (map[string]float64, map[string]float64, map[string]float64, error) {
	if e.featureRemote == nil {
		return nil, nil, nil, ErrDataNotAvailable
	}

	_, paramRef, err := e.GetParamDescriptionListData()
	if err != nil {
		return nil, nil, nil, ErrMetadataNotAvailable
	}

	rData := e.featureRemote.Data(model.FunctionTypeElectricalConnectionPermittedValueSetListData)
	if rData == nil {
		return nil, nil, nil, ErrDataNotAvailable
	}

	data := rData.(*model.ElectricalConnectionPermittedValueSetListDataType)
	if data == nil {
		return nil, nil, nil, ErrDataNotAvailable
	}

	resultSetMin := make(map[string]float64)
	resultSetMax := make(map[string]float64)
	resultSetDefault := make(map[string]float64)
	for _, item := range data.ElectricalConnectionPermittedValueSetData {
		if item.ElectricalConnectionId == nil || item.PermittedValueSet == nil {
			continue
		}

		param, exists := paramRef[*item.ParameterId]
		if !exists {
			continue
		}

		if param.AcMeasuredPhases == nil {
			continue
		}

		for _, set := range item.PermittedValueSet {
			if set.Value != nil && len(set.Value) > 0 {
				resultSetDefault[string(*param.AcMeasuredPhases)] = set.Value[0].GetValue()
			}
			if set.Range != nil {
				for _, rangeItem := range set.Range {
					if rangeItem.Min != nil {
						resultSetMin[string(*param.AcMeasuredPhases)] = rangeItem.Min.GetValue()
					}
					if rangeItem.Max != nil {
						resultSetMax[string(*param.AcMeasuredPhases)] = rangeItem.Max.GetValue()
					}
				}
			}
		}
	}

	if len(resultSetMin) == 0 && len(resultSetMax) == 0 && len(resultSetMax) == 0 {
		return nil, nil, nil, ErrDataNotAvailable
	}

	return resultSetMin, resultSetMax, resultSetDefault, nil
}

// return current values for Electrical Limits
//
// EV only: Min power data is only provided via IEC61851 or using VAS in ISO15118-2.
func (e *ElectricalConnection) GetEVLimitValues() ([]ElectricalLimitType, error) {
	if e.featureRemote == nil {
		return nil, ErrDataNotAvailable
	}

	rData := e.featureRemote.Data(model.FunctionTypeElectricalConnectionParameterDescriptionListData)
	if rData == nil {
		return nil, ErrMetadataNotAvailable
	}

	paramDescriptionData := rData.(*model.ElectricalConnectionParameterDescriptionListDataType)
	if paramDescriptionData == nil {
		return nil, ErrMetadataNotAvailable
	}

	paramRef := make(map[model.ElectricalConnectionParameterIdType]model.ElectricalConnectionParameterDescriptionDataType)
	for _, item := range paramDescriptionData.ElectricalConnectionParameterDescriptionData {
		if item.ParameterId == nil {
			continue
		}
		paramRef[*item.ParameterId] = item
	}

	rData2 := e.featureRemote.Data(model.FunctionTypeElectricalConnectionPermittedValueSetListData)
	if rData2 == nil {
		return nil, ErrDataNotAvailable
	}

	data := rData2.(*model.ElectricalConnectionPermittedValueSetListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	var resultSet []ElectricalLimitType

	for _, item := range data.ElectricalConnectionPermittedValueSetData {
		if item.ParameterId == nil || item.ElectricalConnectionId == nil {
			continue
		}
		param, exists := paramRef[*item.ParameterId]
		if !exists {
			continue
		}

		if len(item.PermittedValueSet) == 0 {
			continue
		}

		var value, minValue, maxValue float64
		hasValue := false
		hasRange := false

		for _, element := range item.PermittedValueSet {
			// is a value set
			if element.Value != nil && len(element.Value) > 0 {
				value = element.Value[0].GetValue()
				hasValue = true
			}
			// is a range set
			if element.Range != nil && len(element.Range) > 0 {
				if element.Range[0].Min != nil {
					minValue = element.Range[0].Min.GetValue()
				}
				if element.Range[0].Max != nil {
					maxValue = element.Range[0].Max.GetValue()
				}
				hasRange = true
			}
		}

		switch {
		// AC Total Power Limits
		case param.ScopeType != nil && *param.ScopeType == model.ScopeTypeTypeACPowerTotal && hasRange:
			result := ElectricalLimitType{
				ConnectionID: uint(*item.ElectricalConnectionId),
				Min:          minValue,
				Max:          maxValue,
				Scope:        model.ScopeTypeTypeACPowerTotal,
			}
			resultSet = append(resultSet, result)

		case param.AcMeasuredPhases != nil && hasRange && hasValue:
			// AC Phase Current Limits
			result := ElectricalLimitType{
				ConnectionID: uint(*item.ElectricalConnectionId),
				Min:          minValue,
				Max:          maxValue,
				Default:      value,
				Phase:        *param.AcMeasuredPhases,
				Scope:        model.ScopeTypeTypeACCurrent,
			}
			resultSet = append(resultSet, result)
		}
	}

	return resultSet, nil
}
