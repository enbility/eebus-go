package cevc

import (
	"errors"
	"time"

	"github.com/enbility/eebus-go/api"
	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/ship-go/logging"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
)

// returns the minimum and maximum number of incentive slots allowed
func (e *CEVC) IncentiveConstraints(entity spineapi.EntityRemoteInterface) (ucapi.IncentiveSlotConstraints, error) {
	result := ucapi.IncentiveSlotConstraints{}

	if !e.IsCompatibleEntityType(entity) {
		return result, api.ErrNoCompatibleEntity
	}

	evIncentiveTable, err := client.NewIncentiveTable(e.LocalEntity, entity)
	if err != nil {
		return result, api.ErrDataNotAvailable
	}

	constraints, err := evIncentiveTable.GetConstraints()
	if err != nil {
		return result, err
	}

	// only use the first constraint
	constraint := constraints[0]

	if constraint.IncentiveSlotConstraints.SlotCountMin != nil {
		result.MinSlots = uint(*constraint.IncentiveSlotConstraints.SlotCountMin)
	}
	if constraint.IncentiveSlotConstraints.SlotCountMax != nil {
		result.MaxSlots = uint(*constraint.IncentiveSlotConstraints.SlotCountMax)
	}

	return result, nil
}

// inform the EVSE about used currency and boundary units
//
// SPINE UC CoordinatedEVCharging 2.4.3
func (e *CEVC) WriteIncentiveTableDescriptions(entity spineapi.EntityRemoteInterface, data []ucapi.IncentiveTariffDescription) error {
	if !e.IsCompatibleEntityType(entity) {
		return api.ErrNoCompatibleEntity
	}

	evIncentiveTable, err := client.NewIncentiveTable(e.LocalEntity, entity)
	if err != nil {
		logging.Log().Error("incentivetable feature not found")
		return err
	}

	filter := model.TariffDescriptionDataType{
		ScopeType: util.Ptr(model.ScopeTypeTypeSimpleIncentiveTable),
	}
	descriptions, err := evIncentiveTable.GetDescriptionsForFilter(filter)
	if err != nil {
		logging.Log().Error(err)
		return err
	}

	// default tariff
	//
	// - tariff, min 1
	//   each tariff has
	//   - tiers: min 1, max 3
	//     each tier has:
	//     - boundaries: min 1, used for different power limits, e.g. 0-1kW x€, 1-3kW y€, ...
	//     - incentives: min 1, max 3
	//       - price/costs (absolute or relative)
	//       - renewable energy percentage
	//       - CO2 emissions
	//
	// limit this to
	// - 1 tariff
	//   - 1 tier
	//     - 1 boundary
	//     - 1 incentive (price)
	//       incentive type has to be the same for all sent power limits!
	descData := []model.IncentiveTableDescriptionType{
		{
			TariffDescription: descriptions[0].TariffDescription,
			Tier: []model.IncentiveTableDescriptionTierType{
				{
					TierDescription: &model.TierDescriptionDataType{
						TierId:   util.Ptr(model.TierIdType(0)),
						TierType: util.Ptr(model.TierTypeTypeDynamicCost),
					},
					BoundaryDescription: []model.TierBoundaryDescriptionDataType{
						{
							BoundaryId:   util.Ptr(model.TierBoundaryIdType(0)),
							BoundaryType: util.Ptr(model.TierBoundaryTypeTypePowerBoundary),
							BoundaryUnit: util.Ptr(model.UnitOfMeasurementTypeW),
						},
					},
					IncentiveDescription: []model.IncentiveDescriptionDataType{
						{
							IncentiveId:   util.Ptr(model.IncentiveIdType(0)),
							IncentiveType: util.Ptr(model.IncentiveTypeTypeAbsoluteCost),
							Currency:      util.Ptr(model.CurrencyTypeEur),
						},
					},
				},
			},
		},
	}

	if len(data) > 0 && len(data[0].Tiers) > 0 {
		newDescData := []model.IncentiveTableDescriptionType{}
		allDataPresent := false

		for index, tariff := range data {
			tariffDesc := descriptions[0].TariffDescription
			if len(descriptions) > index {
				tariffDesc = descriptions[index].TariffDescription
			}

			newTariff := model.IncentiveTableDescriptionType{
				TariffDescription: tariffDesc,
			}

			tierData := []model.IncentiveTableDescriptionTierType{}
			for _, tier := range tariff.Tiers {
				newTier := model.IncentiveTableDescriptionTierType{}

				newTier.TierDescription = &model.TierDescriptionDataType{
					TierId:   util.Ptr(model.TierIdType(tier.Id)),
					TierType: util.Ptr(tier.Type),
				}

				boundaryDescription := []model.TierBoundaryDescriptionDataType{}
				for _, boundary := range tier.Boundaries {
					newBoundary := model.TierBoundaryDescriptionDataType{
						BoundaryId:   util.Ptr(model.TierBoundaryIdType(boundary.Id)),
						BoundaryType: util.Ptr(boundary.Type),
						BoundaryUnit: util.Ptr(boundary.Unit),
					}
					boundaryDescription = append(boundaryDescription, newBoundary)
				}
				newTier.BoundaryDescription = boundaryDescription

				incentiveDescription := []model.IncentiveDescriptionDataType{}
				for _, incentive := range tier.Incentives {
					newIncentive := model.IncentiveDescriptionDataType{
						IncentiveId:   util.Ptr(model.IncentiveIdType(incentive.Id)),
						IncentiveType: util.Ptr(incentive.Type),
					}
					if incentive.Currency != "" {
						newIncentive.Currency = util.Ptr(incentive.Currency)
					}
					incentiveDescription = append(incentiveDescription, newIncentive)
				}
				newTier.IncentiveDescription = incentiveDescription

				if len(newTier.BoundaryDescription) > 0 &&
					len(newTier.IncentiveDescription) > 0 {
					allDataPresent = true
				}
				tierData = append(tierData, newTier)
			}

			newTariff.Tier = tierData

			newDescData = append(newDescData, newTariff)
		}

		if allDataPresent {
			descData = newDescData
		}
	}

	_, err = evIncentiveTable.WriteDescriptions(descData)
	if err != nil {
		logging.Log().Error(err)
		return err
	}

	return nil
}

// send incentives to the EV
// if no data is provided, default incentives with the same price for 7 days will be sent
func (e *CEVC) WriteIncentives(entity spineapi.EntityRemoteInterface, data []ucapi.DurationSlotValue) error {
	if !e.IsCompatibleEntityType(entity) {
		return api.ErrNoCompatibleEntity
	}

	evIncentiveTable, err := client.NewIncentiveTable(e.LocalEntity, entity)
	if err != nil {
		return api.ErrDataNotAvailable
	}

	if len(data) == 0 {
		// send default incentives for the maximum timeframe
		// to fullfill spec, as there is no data provided
		logging.Log().Info("Fallback sending default incentives")
		data = []ucapi.DurationSlotValue{
			{Duration: 7 * time.Hour * 24, Value: 0.30},
		}
	}

	constraints, err := e.IncentiveConstraints(entity)
	if err != nil {
		return err
	}

	if constraints.MinSlots != 0 && constraints.MinSlots > uint(len(data)) {
		return errors.New("too few charge slots provided")
	}

	if constraints.MaxSlots != 0 && constraints.MaxSlots < uint(len(data)) {
		return errors.New("too many charge slots provided")
	}

	incentiveSlots := []model.IncentiveTableIncentiveSlotType{}
	var totalDuration time.Duration
	for index, slot := range data {
		relativeStart := totalDuration

		timeInterval := &model.TimeTableDataType{
			StartTime: &model.AbsoluteOrRecurringTimeType{
				Relative: model.NewDurationType(relativeStart),
			},
		}

		// the last slot also needs an End Time
		if index == len(data)-1 {
			relativeEndTime := relativeStart + slot.Duration
			timeInterval.EndTime = &model.AbsoluteOrRecurringTimeType{
				Relative: model.NewDurationType(relativeEndTime),
			}
		}

		incentiveSlot := model.IncentiveTableIncentiveSlotType{
			TimeInterval: timeInterval,
			Tier: []model.IncentiveTableTierType{
				{
					Tier: &model.TierDataType{
						TierId: util.Ptr(model.TierIdType(0)),
					},
					Boundary: []model.TierBoundaryDataType{
						{
							BoundaryId:         util.Ptr(model.TierBoundaryIdType(0)), // only 1 boundary exists
							LowerBoundaryValue: model.NewScaledNumberType(0),
						},
					},
					Incentive: []model.IncentiveDataType{
						{
							IncentiveId: util.Ptr(model.IncentiveIdType(0)), // always use price
							Value:       model.NewScaledNumberType(slot.Value),
						},
					},
				},
			},
		}
		incentiveSlots = append(incentiveSlots, incentiveSlot)

		totalDuration += slot.Duration
	}

	incentiveData := model.IncentiveTableType{
		Tariff: &model.TariffDataType{
			TariffId: util.Ptr(model.TariffIdType(0)),
		},
		IncentiveSlot: incentiveSlots,
	}

	_, err = evIncentiveTable.WriteValues([]model.IncentiveTableType{incentiveData})

	return err
}
