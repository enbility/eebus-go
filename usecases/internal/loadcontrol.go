package internal

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// return the current loadcontrol limits for a categoriy
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
func LoadControlLimits(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	filter model.LoadControlLimitDescriptionDataType,
) (limits []ucapi.LoadLimitsPhase, resultErr error) {
	limits = nil
	resultErr = api.ErrNoCompatibleEntity

	evLoadControl, err := client.NewLoadControl(localEntity, remoteEntity)
	evElectricalConnection, err2 := client.NewElectricalConnection(localEntity, remoteEntity)
	if err != nil || err2 != nil {
		return
	}

	resultErr = api.ErrDataNotAvailable
	// find out the appropriate limitId for each phase value
	// limitDescription contains the measurementId for each limitId
	limitDescriptions, err := evLoadControl.GetLimitDescriptionsForFilter(filter)
	if err != nil || limitDescriptions == nil {
		return
	}

	var result []ucapi.LoadLimitsPhase

	for i := 0; i < len(PhaseNameMapping); i++ {
		phaseName := PhaseNameMapping[i]

		// electricalParameterDescription contains the measured phase for each measurementId
		filter := model.ElectricalConnectionParameterDescriptionDataType{
			AcMeasuredPhases: &phaseName,
		}
		elParamDesc, err := evElectricalConnection.GetParameterDescriptionsForFilter(filter)
		if err != nil || len(elParamDesc) == 0 || elParamDesc[0].MeasurementId == nil {
			// there is no data for this phase, the phase may not exist
			result = append(result, ucapi.LoadLimitsPhase{Phase: phaseName})
			continue
		}

		var limitDesc *model.LoadControlLimitDescriptionDataType
		for _, desc := range limitDescriptions {
			if desc.MeasurementId != nil &&
				elParamDesc[0].MeasurementId != nil &&
				*desc.MeasurementId == *elParamDesc[0].MeasurementId {
				safeDesc := desc
				limitDesc = &safeDesc
				break
			}
		}

		if limitDesc == nil || limitDesc.LimitId == nil {
			return
		}

		limitIdData, err := evLoadControl.GetLimitDataForId(*limitDesc.LimitId)
		if err != nil {
			return
		}

		var limitValue float64
		if limitIdData.Value == nil || (limitIdData.IsLimitActive != nil && !*limitIdData.IsLimitActive) {
			// report maximum possible if no limit is available or the limit is not active
			filter := model.ElectricalConnectionPermittedValueSetDataType{
				ParameterId: elParamDesc[0].ParameterId,
			}
			_, dataMax, _, err := evElectricalConnection.GetPermittedValueDataForFilter(filter)
			if err != nil {
				return
			}

			limitValue = dataMax
		} else {
			limitValue = limitIdData.Value.GetValue()
		}

		newLimit := ucapi.LoadLimitsPhase{
			Phase:        phaseName,
			IsChangeable: (limitIdData.IsLimitChangeable != nil && *limitIdData.IsLimitChangeable),
			IsActive:     (limitIdData.IsLimitActive != nil && *limitIdData.IsLimitActive),
			Value:        limitValue,
		}

		result = append(result, newLimit)
	}

	return result, nil
}

// generic helper to be used in UCOPEV & UCOSCEV
// send new LoadControlLimits to the remote EV
//
// parameters:
//   - limits: a set of limits for a  given limit category containing phase specific limit data
//
// category obligations:
// Sets a maximum A limit for each phase that the EV may not exceed.
// Mainly used for implementing overload protection of the site or limiting the
// maximum charge power of EVs when the EV and EVSE communicate via IEC61851
// and with ISO15118 if the EV does not support the Optimization of Self Consumption
// usecase.
//
// category recommendations:
// Sets a recommended charge power in A for each phase. This is mainly
// used if the EV and EVSE communicate via ISO15118 to support charging excess solar power.
// The EV either needs to support the Optimization of Self Consumption usecase or
// the EVSE needs to be able map the recommendations into oligation limits which then
// works for all EVs communication either via IEC61851 or ISO15118.
//
// notes:
//   - For obligations to work for optimizing solar excess power, the EV needs to have an energy demand.
//   - Recommendations work even if the EV does not have an active energy demand, given it communicated with the EVSE via ISO15118 and supports the usecase.
//   - In ISO15118-2 the usecase is only supported via VAS extensions which are vendor specific and needs to have specific EVSE support for the specific EV brand.
//   - In ISO15118-20 this is a standard feature which does not need special support on the EVSE.
//   - Min power data is only provided via IEC61851 or using VAS in ISO15118-2.
func WriteLoadControlLimits(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	category model.LoadControlCategoryType,
	limits []ucapi.LoadLimitsPhase) (*model.MsgCounterType, error) {
	loadControl, err := client.NewLoadControl(localEntity, remoteEntity)
	electricalConnection, err2 := client.NewElectricalConnection(localEntity, remoteEntity)
	if err != nil || err2 != nil {
		return nil, api.ErrNoCompatibleEntity
	}

	var limitData []model.LoadControlLimitDataType

	for _, phaseLimit := range limits {
		// find out the appropriate limitId for each phase value
		// limitDescription contains the measurementId for each limitId
		filter := model.LoadControlLimitDescriptionDataType{
			LimitCategory: &category,
		}
		limitDescriptions, err := loadControl.GetLimitDescriptionsForFilter(filter)
		if err != nil || limitDescriptions == nil {
			continue
		}

		// electricalParameterDescription contains the measured phase for each measurementId
		filter2 := model.ElectricalConnectionParameterDescriptionDataType{
			AcMeasuredPhases: util.Ptr(phaseLimit.Phase),
		}
		elParamDesc, err := electricalConnection.GetParameterDescriptionsForFilter(filter2)
		if err != nil || len(elParamDesc) == 0 || elParamDesc[0].MeasurementId == nil {
			continue
		}

		var limitDesc *model.LoadControlLimitDescriptionDataType
		for _, desc := range limitDescriptions {
			if desc.MeasurementId != nil &&
				elParamDesc[0].MeasurementId != nil &&
				*desc.MeasurementId == *elParamDesc[0].MeasurementId {
				safeDesc := desc
				limitDesc = &safeDesc
				break
			}
		}

		if limitDesc == nil || limitDesc.LimitId == nil {
			continue
		}

		limitIdData, err := loadControl.GetLimitDataForId(*limitDesc.LimitId)
		if err != nil {
			continue
		}

		// EEBus_UC_TS_OverloadProtectionByEvChargingCurrentCurtailment V1.01b 3.2.1.2.2.2
		// If omitted or set to "true", the timePeriod, value and isLimitActive element SHALL be writeable by a client.
		if limitIdData.IsLimitChangeable != nil && !*limitIdData.IsLimitChangeable {
			continue
		}

		// electricalPermittedValueSet contains the allowed min, max and the default values per phase
		limit := electricalConnection.AdjustValueToBeWithinPermittedValuesForParameterId(
			phaseLimit.Value, *elParamDesc[0].ParameterId)

		newLimit := model.LoadControlLimitDataType{
			LimitId:       limitDesc.LimitId,
			IsLimitActive: util.Ptr(phaseLimit.IsActive),
			Value:         model.NewScaledNumberType(limit),
		}
		limitData = append(limitData, newLimit)
	}

	msgCounter, err := loadControl.WriteLimitData(limitData, nil, nil)

	return msgCounter, err
}
