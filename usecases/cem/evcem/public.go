package evcem

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// return the number of ac connected phases of the EV or 0 if it is unknown
func (e *EVCEM) PhasesConnected(entity spineapi.EntityRemoteInterface) (uint, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	evElectricalConnection, err := client.NewElectricalConnection(e.LocalEntity, entity)
	if err != nil {
		return 0, api.ErrDataNotAvailable
	}

	data, err := evElectricalConnection.GetDescriptionsForFilter(model.ElectricalConnectionDescriptionDataType{})
	if err != nil || len(data) == 0 {
		return 0, api.ErrDataNotAvailable
	}

	for _, item := range data {
		if item.ElectricalConnectionId != nil && item.AcConnectedPhases != nil {
			return *item.AcConnectedPhases, nil
		}
	}

	// default to 0 if the value is not available
	return 0, nil
}

// return the last current measurement for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EVCEM) CurrentPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	evMeasurement, err := client.NewMeasurement(e.LocalEntity, entity)
	evElectricalConnection, err2 := client.NewElectricalConnection(e.LocalEntity, entity)
	if err != nil || err2 != nil {
		return nil, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeCurrent),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACCurrent),
	}
	data, err := evMeasurement.GetDataForFilter(filter)
	if err != nil || len(data) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	var result []float64

	for _, phase := range ucapi.PhaseNameMapping {
		for _, item := range data {
			if item.Value == nil {
				continue
			}

			filter := model.ElectricalConnectionParameterDescriptionDataType{
				MeasurementId: item.MeasurementId,
			}
			elParam, err := evElectricalConnection.GetParameterDescriptionsForFilter(filter)
			if err != nil || len(elParam) == 0 ||
				elParam[0].AcMeasuredPhases == nil || *elParam[0].AcMeasuredPhases != phase {
				continue
			}

			phaseValue := item.Value.GetValue()
			result = append(result, phaseValue)
		}
	}

	return result, nil
}

// return the last power measurement for each phase of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EVCEM) PowerPerPhase(entity spineapi.EntityRemoteInterface) ([]float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return nil, api.ErrNoCompatibleEntity
	}

	evMeasurement, err := client.NewMeasurement(e.LocalEntity, entity)
	evElectricalConnection, err2 := client.NewElectricalConnection(e.LocalEntity, entity)
	if err != nil || err2 != nil {
		return nil, err
	}

	var data []model.MeasurementDataType

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPower),
	}
	data, err = evMeasurement.GetDataForFilter(filter)
	// Elli Charger Connect/Pro (Gen1) returns power descriptions, but only measurements without actual values, see test case Test_EVPowerPerPhase_Current
	if err != nil || len(data) == 0 || data[0].Value == nil {
		return nil, api.ErrDataNotAvailable
	}

	var result []float64

	for _, phase := range ucapi.PhaseNameMapping {
		for _, item := range data {
			if item.Value == nil {
				continue
			}

			filter := model.ElectricalConnectionParameterDescriptionDataType{
				MeasurementId: item.MeasurementId,
			}
			elParam, err := evElectricalConnection.GetParameterDescriptionsForFilter(filter)
			if err != nil || len(elParam) == 0 ||
				*elParam[0].AcMeasuredPhases != phase {
				continue
			}

			phaseValue := item.Value.GetValue()
			result = append(result, phaseValue)
		}
	}

	return result, nil
}

// return the charged energy measurement in Wh of the connected EV
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func (e *EVCEM) EnergyCharged(entity spineapi.EntityRemoteInterface) (float64, error) {
	if !e.IsCompatibleEntityType(entity) {
		return 0, api.ErrNoCompatibleEntity
	}

	evMeasurement, err := client.NewMeasurement(e.LocalEntity, entity)
	if err != nil {
		return 0, err
	}

	filter := model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypeEnergy),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		ScopeType:       util.Ptr(model.ScopeTypeTypeCharge),
	}
	data, err := evMeasurement.GetDataForFilter(filter)
	if err != nil || len(data) == 0 {
		return 0, api.ErrDataNotAvailable
	}

	// we assume there is only one result
	value := data[0].Value
	if value == nil {
		return 0, api.ErrDataNotAvailable
	}

	return value.GetValue(), err
}
