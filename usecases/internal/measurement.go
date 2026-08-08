package internal

import (
	"slices"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	"github.com/enbility/eebus-go/features/server"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// According to the LPC Installation Guide, this value is relatively high to make sure it doesn't conflict with other IDs on this entity
var defaultPowerTotalMeasurementId = model.MeasurementIdType(50)

// return the phase specific measurement data
func MeasurementPhaseSpecificDataForFilter(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	measurementFilter model.MeasurementDescriptionDataType,
	energyDirection model.EnergyDirectionType,
	validPhaseNameTypes []model.ElectricalConnectionPhaseNameType,
) (map[model.ElectricalConnectionPhaseNameType]float64, error) {
	measurement, err := client.NewMeasurement(localEntity, remoteEntity)
	electricalConnection, err1 := client.NewElectricalConnection(localEntity, remoteEntity)
	if err != nil || err1 != nil {
		return nil, api.ErrMetadataNotAvailable
	}

	data, err := measurement.GetDataForFilter(measurementFilter)
	if err != nil || len(data) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	result := make(map[model.ElectricalConnectionPhaseNameType]float64, len(validPhaseNameTypes))

	for _, item := range data {
		if item.Value == nil || item.MeasurementId == nil {
			continue
		}

		filter := model.ElectricalConnectionParameterDescriptionDataType{
			MeasurementId: item.MeasurementId,
		}
		param, err := electricalConnection.GetParameterDescriptionsForFilter(filter)
		if err != nil || len(param) == 0 {
			// error getting parameter description
			continue
		}

		phaseName, ok := phaseNameFromParameterDescription(param[0], validPhaseNameTypes == nil)
		if !ok {
			// error getting parameter description
			continue
		}

		if validPhaseNameTypes != nil &&
			!slices.Contains(validPhaseNameTypes, phaseName) {
			// ignore phase measurements not specified in validPhaseNameTypes
			continue
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

		result[phaseName] = value
	}

	return result, nil
}

// GetPowerTotalMeasurementId returns the MeasurementId for the AC Power Total measurement
func GetPowerTotalMeasurementId(localEntity spineapi.EntityLocalInterface) model.MeasurementIdType {
	if localEntity == nil {
		return defaultPowerTotalMeasurementId
	}
	measurementFeat := localEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeMeasurement, model.RoleTypeServer)
	if measurementFeat == nil {
		return defaultPowerTotalMeasurementId
	}
	measurement, err := server.NewMeasurement(localEntity)
	if err != nil || measurement == nil {
		return defaultPowerTotalMeasurementId
	}
	MeasurementDescriptionData, err := measurement.GetDescriptionsForFilter(model.MeasurementDescriptionDataType{
		MeasurementType: util.Ptr(model.MeasurementTypeTypePower),
		CommodityType:   util.Ptr(model.CommodityTypeTypeElectricity),
		Unit:            util.Ptr(model.UnitOfMeasurementTypeW),
		ScopeType:       util.Ptr(model.ScopeTypeTypeACPowerTotal),
	})
	if err != nil || len(MeasurementDescriptionData) != 1 || MeasurementDescriptionData[0].MeasurementId == nil {
		return defaultPowerTotalMeasurementId
	}

	return *MeasurementDescriptionData[0].MeasurementId
}

func phaseNameFromParameterDescription(
	param model.ElectricalConnectionParameterDescriptionDataType,
	allowUnset bool,
) (model.ElectricalConnectionPhaseNameType, bool) {
	if param.AcMeasuredPhases == nil {
		if allowUnset {
			return model.ElectricalConnectionPhaseNameTypeNone, true
		}
		return "", false
	}

	phaseName := *param.AcMeasuredPhases
	if param.AcMeasuredInReferenceTo == nil {
		return phaseName, true
	}

	referencePhaseName := *param.AcMeasuredInReferenceTo
	if referencePhaseName == phaseName ||
		referencePhaseName == model.ElectricalConnectionPhaseNameTypeNeutral ||
		referencePhaseName == model.ElectricalConnectionPhaseNameTypeGround ||
		referencePhaseName == model.ElectricalConnectionPhaseNameTypeNone {
		return phaseName, true
	}

	if combinedPhaseName, ok := phasePairName(phaseName, referencePhaseName); ok {
		return combinedPhaseName, true
	}

	return phaseName, true
}

func phasePairName(
	phaseName model.ElectricalConnectionPhaseNameType,
	referencePhaseName model.ElectricalConnectionPhaseNameType,
) (model.ElectricalConnectionPhaseNameType, bool) {
	switch {
	case isPhasePair(phaseName, referencePhaseName, model.ElectricalConnectionPhaseNameTypeA, model.ElectricalConnectionPhaseNameTypeB):
		return model.ElectricalConnectionPhaseNameTypeAb, true
	case isPhasePair(phaseName, referencePhaseName, model.ElectricalConnectionPhaseNameTypeB, model.ElectricalConnectionPhaseNameTypeC):
		return model.ElectricalConnectionPhaseNameTypeBc, true
	case isPhasePair(phaseName, referencePhaseName, model.ElectricalConnectionPhaseNameTypeA, model.ElectricalConnectionPhaseNameTypeC):
		return model.ElectricalConnectionPhaseNameTypeAc, true
	default:
		return "", false
	}
}

func isPhasePair(
	phaseName model.ElectricalConnectionPhaseNameType,
	referencePhaseName model.ElectricalConnectionPhaseNameType,
	first model.ElectricalConnectionPhaseNameType,
	second model.ElectricalConnectionPhaseNameType,
) bool {
	return (phaseName == first && referencePhaseName == second) ||
		(phaseName == second && referencePhaseName == first)
}
