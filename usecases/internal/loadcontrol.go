package internal

import (
	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// return the current loadcontrol limits for a categoriy
//
// possible errors:
//   - ErrDataNotAvailable if no such measurement is (yet) available
//   - and others
//
// Notes:
//   - If a limit of phase is not active, the value returned is set to the maximum permitted value (if available), otherwise the phase value is not returned
//   - If at least one phase value is available, no error is returned
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

	for i := 0; i < len(ucapi.PhaseNameMapping); i++ {
		phaseName := ucapi.PhaseNameMapping[i]

		// electricalParameterDescription contains the measured phase for each measurementId
		filter := model.ElectricalConnectionParameterDescriptionDataType{
			AcMeasuredPhases: &phaseName,
		}
		elParamDesc, err := evElectricalConnection.GetParameterDescriptionsForFilter(filter)
		if err != nil || len(elParamDesc) == 0 || elParamDesc[0].MeasurementId == nil {
			// there is no data for this phase, the phase may not exist
			// so do not add id to the result
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
				// the device provides no limits for this phase (e.g. Bender 1 Phase)
				// according to OpEV Spec 1.0.1b, page 30: "At least one set of permitted values SHALL be stated."
				// which is not the case for all elements for this EVSE setup
				// but we have to ignore this case and continue with the next phase

				// In addition Elli Connect/Pro (Gen1) returns no or empty PermittedValueSet data
				continue
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

	if len(result) == 0 {
		return nil, api.ErrDataNotAvailable
	}

	return result, nil
}

// generic helper to be used in UCLPC & UCLPP
// send new LoadControlLimit to the remote entity
//
// parameters:
//   - limits: a set of limits for a  given limit category containing phase specific limit data
//   - resultCB: callback function for handling the result response
func WriteLoadControlLimit(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	filter model.LoadControlLimitDescriptionDataType,
	limit ucapi.LoadLimit,
	resultCB func(result model.ResultDataType),
) (*model.MsgCounterType, error) {
	loadControl, err := client.NewLoadControl(localEntity, remoteEntity)
	if err != nil {
		return nil, api.ErrNoCompatibleEntity
	}

	var limitData []model.LoadControlLimitDataType

	currentLimits, err := loadControl.GetLimitDataForFilter(filter)
	if err != nil || len(currentLimits) != 1 || currentLimits[0].LimitId == nil {
		return nil, api.ErrDataNotAvailable
	}

	item := currentLimits[0]

	// EEBus_UC_TS_LimitationOfPowerConsumption V1.0.0 3.2.2.2.2.2
	// If set to "true", the timePeriod, value and isLimitActive Elements SHALL be writeable by a client.
	if item.IsLimitChangeable != nil && !*item.IsLimitChangeable {
		return nil, api.ErrNotSupported
	}

	var deleteSelectors *model.LoadControlLimitListDataSelectorsType
	var deleteElements *model.LoadControlLimitDataElementsType

	newLimit := model.LoadControlLimitDataType{
		LimitId:       item.LimitId,
		IsLimitActive: util.Ptr(limit.IsActive),
		Value:         model.NewScaledNumberType(limit.Value),
	}
	if limit.Duration > 0 {
		newLimit.TimePeriod = &model.TimePeriodType{
			EndTime: model.NewAbsoluteOrRelativeTimeTypeFromDuration(limit.Duration),
		}
	}

	// should we delete the TimePeriod value?
	if limit.DeleteDuration {
		deleteSelectors = &model.LoadControlLimitListDataSelectorsType{
			LimitId: currentLimits[0].LimitId,
		}
		deleteElements = &model.LoadControlLimitDataElementsType{
			TimePeriod: &model.TimePeriodElementsType{},
		}
	}

	limitData = append(limitData, newLimit)

	msgCounter, err := loadControl.WriteLimitData(limitData, deleteSelectors, deleteElements)

	if resultCB != nil && msgCounter != nil {
		cb := func(msg spineapi.ResponseMessage) {
			response, ok := msg.Data.(*model.ResultDataType)
			if ok {
				resultCB(*response)
			}
		}
		if errCB := loadControl.AddResponseCallback(*msgCounter, cb); errCB != nil {
			logging.Log().Debug("Failed to add response callback for msgCounter %v: %v", msgCounter, errCB)
		}
	}

	return msgCounter, err
}

// generic helper to be used in UCOPEV & UCOSCEV
// send new LoadControlLimits to the remote EV
//
// parameters:
//   - limits: a set of limits for a  given limit category containing phase specific limit data
//   - resultCB: callback function for handling the result response
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
func WriteLoadControlPhaseLimits(
	localEntity spineapi.EntityLocalInterface,
	remoteEntity spineapi.EntityRemoteInterface,
	filter model.LoadControlLimitDescriptionDataType,
	limits []ucapi.LoadLimitsPhase,
	resultCB func(result model.ResultDataType),
) (*model.MsgCounterType, error) {
	loadControl, err := client.NewLoadControl(localEntity, remoteEntity)
	electricalConnection, err2 := client.NewElectricalConnection(localEntity, remoteEntity)
	if err != nil || err2 != nil {
		return nil, api.ErrNoCompatibleEntity
	}

	var limitData []model.LoadControlLimitDataType

	for _, phaseLimit := range limits {
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

	if resultCB != nil && msgCounter != nil {
		cb := func(msg spineapi.ResponseMessage) {
			response, ok := msg.Data.(*model.ResultDataType)
			if ok {
				resultCB(*response)
			}
		}

		if errCB := loadControl.AddResponseCallback(*msgCounter, cb); errCB != nil {
			logging.Log().Debug("Failed to add response callback for msgCounter %v: %v", msgCounter, errCB)
		}
	}

	return msgCounter, err
}
