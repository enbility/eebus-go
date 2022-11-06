package features

import (
	"time"

	"github.com/DerAndereAndi/eebus-go/logging"
	"github.com/DerAndereAndi/eebus-go/spine"
	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type MeasurementType struct {
	MeasurementId uint
	Value         float64
	ValueMin      float64
	ValueMax      float64
	ValueStep     float64
	Unit          model.UnitOfMeasurementType
	Scope         model.ScopeTypeType
	Timestamp     time.Time
}

type Measurement struct {
	*FeatureImpl
}

func NewMeasurement(localRole, remoteRole model.RoleType, spineLocalDevice *spine.DeviceLocalImpl, entity *spine.EntityRemoteImpl) (*Measurement, error) {
	feature, err := NewFeatureImpl(model.FeatureTypeTypeMeasurement, localRole, remoteRole, spineLocalDevice, entity)
	if err != nil {
		return nil, err
	}

	m := &Measurement{
		FeatureImpl: feature,
	}

	return m, nil
}

// request FunctionTypeMeasurementDescriptionListData from a remote device
func (m *Measurement) RequestDescription() error {
	if _, err := m.requestData(model.FunctionTypeMeasurementDescriptionListData); err != nil {
		logging.Log.Error(err)
		return err
	}

	return nil
}

// request FunctionTypeMeasurementConstraintsListData from a remote entity
func (m *Measurement) RequestConstraints() error {
	if _, err := m.requestData(model.FunctionTypeMeasurementConstraintsListData); err != nil {
		logging.Log.Error(err)
		return err
	}

	return nil
}

// request FunctionTypeMeasurementListData from a remote entity
func (m *Measurement) Request() (*model.MsgCounterType, error) {
	msgCounter, err := m.requestData(model.FunctionTypeMeasurementListData)
	if err != nil {
		logging.Log.Error(err)
		return nil, err
	}

	return msgCounter, nil
}

// return current value of a defined scope
func (m *Measurement) GetValueForScope(scope model.ScopeTypeType, electricalConnection *ElectricalConnection) (float64, error) {
	descRef, err := m.GetDescription()
	if err != nil {
		return 0.0, ErrMetadataNotAvailable
	}

	data := m.featureRemote.Data(model.FunctionTypeMeasurementListData).(*model.MeasurementListDataType)
	if data == nil {
		return 0.0, ErrDataNotAvailable
	}

	var result float64
	for _, item := range data.MeasurementData {
		if item.MeasurementId == nil || item.Value == nil {
			continue
		}

		desc, exists := descRef[*item.MeasurementId]
		if !exists {
			continue
		}

		if desc.ScopeType == nil {
			continue
		}

		if *desc.ScopeType == scope {
			return item.Value.GetValue(), nil
		}
	}

	return result, nil
}

// return current values of a defined scope per phase
//
// returns a map with the phase ("a", "b", "c") as a key
func (m *Measurement) GetValuesPerPhaseForScope(scope model.ScopeTypeType, electricalConnection *ElectricalConnection) (map[string]float64, error) {
	descRef, err := m.GetDescription()
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	paramRef, _, err := electricalConnection.GetParamDescriptionListData()
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	data := m.featureRemote.Data(model.FunctionTypeMeasurementListData).(*model.MeasurementListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	resultSet := make(map[string]float64)
	for _, item := range data.MeasurementData {
		if item.MeasurementId == nil || item.Value == nil {
			continue
		}

		param, exists := paramRef[*item.MeasurementId]
		if !exists {
			continue
		}

		desc, exists := descRef[*item.MeasurementId]
		if !exists {
			continue
		}

		if desc.ScopeType == nil || param.AcMeasuredPhases == nil {
			continue
		}

		if *desc.ScopeType == scope {
			resultSet[string(*param.AcMeasuredPhases)] = item.Value.GetValue()
		}
	}
	if len(resultSet) == 0 {
		return nil, ErrDataNotAvailable
	}

	return resultSet, nil
}

type measurementDescriptionMap map[model.MeasurementIdType]model.MeasurementDescriptionDataType

// return a map of MeasurementDescriptionListDataType with measurementId as key
// returns an error if no description data is available yet
func (m *Measurement) GetDescription() (measurementDescriptionMap, error) {
	data := m.featureRemote.Data(model.FunctionTypeMeasurementDescriptionListData).(*model.MeasurementDescriptionListDataType)
	if data == nil {
		return nil, ErrMetadataNotAvailable
	}
	ref := make(measurementDescriptionMap)
	for _, item := range data.MeasurementDescriptionData {
		if item.MeasurementId == nil {
			continue
		}
		ref[*item.MeasurementId] = item
	}
	return ref, nil
}

// return a map of MeasurementDescriptionListDataType with measurementId as key for a given scope
// returns an error if no description data is available yet
func (m *Measurement) GetDescriptionForScope(scope model.ScopeTypeType) (measurementDescriptionMap, error) {
	data, err := m.GetDescription()
	if err != nil {
		return nil, err
	}

	ref := make(measurementDescriptionMap)
	for _, item := range data {
		if item.MeasurementId == nil || item.ScopeType == nil {
			continue
		}
		if *item.ScopeType == scope {
			ref[*item.MeasurementId] = item
		}
	}

	return ref, ErrDataNotAvailable
}

// return current SoC for measurements
func (m *Measurement) GetSoC() (float64, error) {
	descRef, err := m.GetDescription()
	if err != nil {
		return 0.0, ErrMetadataNotAvailable
	}

	data := m.featureRemote.Data(model.FunctionTypeMeasurementListData).(*model.MeasurementListDataType)
	if data == nil {
		return 0.0, ErrDataNotAvailable
	}

	for _, item := range data.MeasurementData {
		if item.MeasurementId == nil || item.Value == nil {
			continue
		}

		desc, exists := descRef[*item.MeasurementId]
		if !exists {
			continue
		}

		if desc.ScopeType == nil {
			continue
		}

		if *desc.ScopeType == model.ScopeTypeTypeStateOfCharge {
			return item.Value.GetValue(), nil
		}
	}

	return 0.0, ErrDataNotAvailable
}

type measurementConstraintMap map[model.MeasurementIdType]model.MeasurementConstraintsDataType

// return a map of MeasurementDescriptionListDataType with measurementId as key
func (m *Measurement) GetConstraints() (measurementConstraintMap, error) {
	data := m.featureRemote.Data(model.FunctionTypeMeasurementConstraintsListData).(*model.MeasurementConstraintsListDataType)
	if data == nil {
		return nil, ErrMetadataNotAvailable
	}
	ref := make(measurementConstraintMap)
	for _, item := range data.MeasurementConstraintsData {
		if item.MeasurementId == nil {
			continue
		}
		ref[*item.MeasurementId] = item
	}
	return ref, nil
}

// return current values for measurements
func (m *Measurement) GetValues() ([]MeasurementType, error) {
	// constraints can be optional
	constraintsRef, _ := m.GetConstraints()

	descRef, err := m.GetDescription()
	if err != nil {
		return nil, ErrMetadataNotAvailable
	}

	data := m.featureRemote.Data(model.FunctionTypeMeasurementListData).(*model.MeasurementListDataType)
	if data == nil {
		return nil, ErrDataNotAvailable
	}

	var resultSet []MeasurementType
	for _, item := range data.MeasurementData {
		if item.MeasurementId == nil {
			continue
		}

		desc, exists := descRef[*item.MeasurementId]
		if !exists {
			continue
		}

		result := MeasurementType{
			MeasurementId: uint(*item.MeasurementId),
		}

		if item.Value != nil {
			result.Value = item.Value.GetValue()
		}

		if item.Timestamp != nil {
			if value, err := time.Parse(time.RFC3339, *item.Timestamp); err == nil {
				result.Timestamp = value
			}
		}

		if desc.ScopeType != nil {
			result.Scope = *desc.ScopeType
		}
		if desc.Unit != nil {
			result.Unit = *desc.Unit
		}

		constraint, exists := constraintsRef[*item.MeasurementId]
		if exists {
			if constraint.ValueRangeMin != nil {
				result.ValueMin = constraint.ValueRangeMin.GetValue()
			}
			if constraint.ValueRangeMax != nil {
				result.ValueMax = constraint.ValueRangeMax.GetValue()
			}
			if constraint.ValueStepSize != nil {
				result.ValueStep = constraint.ValueStepSize.GetValue()
			}
		}

		resultSet = append(resultSet, result)
	}

	return resultSet, nil
}
