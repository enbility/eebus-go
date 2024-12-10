package internal

import (
	"slices"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
)

// return the phase specific measurement data
func MeasurementPhaseSpecificDataForFilter(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	measurementFilter model.MeasurementDescriptionDataType,
	energyDirection model.EnergyDirectionType,
	validPhaseNameTypes []model.ElectricalConnectionPhaseNameType,
) ([]float64, error) {
	measurement, err := client.NewMeasurement(localEntity, remoteEntity)
	electricalConnection, err1 := client.NewElectricalConnection(localEntity, remoteEntity)
	if err != nil || err1 != nil {
		return nil, api.ErrMetadataNotAvailable
	}

	data, err := measurement.GetDataForFilter(measurementFilter)
	if err != nil || len(data) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	var result []float64

	for _, item := range data {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		if validPhaseNameTypes != nil {
			filter := model.ElectricalConnectionParameterDescriptionDataType{
				MeasurementId: item.MeasurementId,
			}
			param, err := electricalConnection.GetParameterDescriptionsForFilter(filter)
			if err != nil || len(param) == 0 ||
				param[0].AcMeasuredPhases == nil ||
				!slices.Contains(validPhaseNameTypes, *param[0].AcMeasuredPhases) {
				continue
			}
		}

		if energyDirection != "" {
			filter := model.ElectricalConnectionParameterDescriptionDataType{
				MeasurementId: item.MeasurementId,
			}
			desc, err := electricalConnection.GetDescriptionForParameterDescriptionFilter(filter)
			if err != nil || desc == nil {
				continue
			}

			// if energy direction is not consume
			if desc.PositiveEnergyDirection == nil || *desc.PositiveEnergyDirection != energyDirection {
				return nil, err
			}
		}

		// if the value state is set and not normal, the value is not valid and should be ignored
		// therefore we return an error
		if item.ValueState != nil && *item.ValueState != model.MeasurementValueStateTypeNormal {
			return nil, api.ErrDataInvalid
		}

		value := item.Value.GetValue()

		result = append(result, value)
	}

	return result, nil
}
