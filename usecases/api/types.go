package api

import (
	"time"

	"github.com/enbility/spine-go/model"
)

type EVChargeStateType string

const (
	EVChargeStateTypeUnknown   EVChargeStateType = "Unknown"
	EVChargeStateTypeUnplugged EVChargeStateType = "unplugged"
	EVChargeStateTypeError     EVChargeStateType = "error"
	EVChargeStateTypePaused    EVChargeStateType = "paused"
	EVChargeStateTypeActive    EVChargeStateType = "active"
	EVChargeStateTypeFinished  EVChargeStateType = "finished"
)

// manufacturer data type
type ManufacturerData struct {
	DeviceName                     string
	DeviceCode                     string
	SerialNumber                   string
	SoftwareRevision               string
	HardwareRevision               string
	VendorName                     string
	VendorCode                     string
	BrandName                      string
	PowerSource                    string
	ManufacturerNodeIdentification string
	ManufacturerLabel              string
	ManufacturerDescription        string
}

// Defines a phase specific limit data set
type LoadLimitsPhase struct {
	Phase        model.ElectricalConnectionPhaseNameType // the phase
	IsChangeable bool                                    // if the value can be changed via write, ignored when writing data
	IsActive     bool                                    // if the limit is active
	Value        float64                                 // the limit
}

// Defines a limit data set
type LoadLimit struct {
	Duration       time.Duration // the duration of the limit,
	IsChangeable   bool          // if the value can be changed via write, ignored when writing data
	IsActive       bool          // if the limit is active
	Value          float64       // the limit
	DeleteDuration bool          // if the Duration (TimePeriod in SPINE) should be deleted (only used for write commands. Relevant for LPC & LPP only)
}

// identification
type IdentificationItem struct {
	// the identification value
	Value string

	// the type of the identification value, e.g.
	ValueType model.IdentificationTypeType
}

type EVChargeStrategyType string

const (
	EVChargeStrategyTypeUnknown        EVChargeStrategyType = "unknown"
	EVChargeStrategyTypeNoDemand       EVChargeStrategyType = "nodemand"
	EVChargeStrategyTypeDirectCharging EVChargeStrategyType = "directcharging"
	EVChargeStrategyTypeMinSoC         EVChargeStrategyType = "minsoc"
	EVChargeStrategyTypeTimedCharging  EVChargeStrategyType = "timedcharging"
)

// Contains details about the actual demands from the EV
//
// General:
//   - If duration and energy is 0, charge mode is EVChargeStrategyTypeNoDemand
//   - If duration is 0, charge mode is EVChargeStrategyTypeDirectCharging and the slots should cover at least 48h
//   - If both are != 0, charge mode is EVChargeStrategyTypeTimedCharging and the slots should cover at least the duration, but at max 168h (7d)
type Demand struct {
	MinDemand          float64 // minimum demand in Wh to reach the minSoC setting, 0 if not set
	OptDemand          float64 // demand in Wh to reach the timer SoC setting
	MaxDemand          float64 // the maximum possible demand until the battery is full
	DurationUntilStart float64 // the duration in s from now until charging will start, this could be in the future but usualy is now
	DurationUntilEnd   float64 // the duration in s from now until minDemand or optDemand has to be reached, 0 if direct charge strategy is active
}

// Contains details about an EV generated charging plan
type ChargePlan struct {
	Slots []ChargePlanSlotValue // Individual charging slot details
}

// Contains details about a charging plan slot
type ChargePlanSlotValue struct {
	Start    time.Time // The start time of the slot
	End      time.Time // The duration of the slot
	Value    float64   // planned power value
	MinValue float64   // minimum power value
	MaxValue float64   // maximum power value
}

// Details about the time slot constraints
type TimeSlotConstraints struct {
	MinSlots             uint          // the minimum number of slots, no minimum if 0
	MaxSlots             uint          // the maximum number of slots, unlimited if 0
	MinSlotDuration      time.Duration // the minimum duration of a slot, no minimum if 0
	MaxSlotDuration      time.Duration // the maximum duration of a slot, unlimited if 0
	SlotDurationStepSize time.Duration // the duration has to be a multiple of this value if != 0
}

// Details about the incentive slot constraints
type IncentiveSlotConstraints struct {
	MinSlots uint // the minimum number of slots, no minimum if 0
	MaxSlots uint // the maximum number of slots, unlimited if 0
}

// details about the boundary
type TierBoundaryDescription struct {
	// the id of the boundary
	Id uint

	// the type of the boundary
	Type model.TierBoundaryTypeType

	// the unit of the boundary
	Unit model.UnitOfMeasurementType
}

// details about incentive
type IncentiveDescription struct {
	// the id of the incentive
	Id uint

	// the type of the incentive
	Type model.IncentiveTypeType

	// the currency of the incentive, if it is price based
	Currency model.CurrencyType
}

// Contains about one tier in a tariff
type IncentiveTableDescriptionTier struct {
	// the id of the tier
	Id uint

	// the tiers type
	Type model.TierTypeType

	// each tear has 1 to 3 boundaries
	// used for different power limits, e.g. 0-1kW x€, 1-3kW y€, ...
	Boundaries []TierBoundaryDescription

	// each tier has 1 to 3 incentives
	//   - price/costs (absolute or relative)
	//   - renewable energy percentage
	//   - CO2 emissions
	Incentives []IncentiveDescription
}

// Contains details about a tariff
type IncentiveTariffDescription struct {
	// each tariff can have 1 to 3 tiers
	Tiers []IncentiveTableDescriptionTier
}

// Contains details about power limits or incentives for a defined timeframe
type DurationSlotValue struct {
	Duration time.Duration // Duration of this slot
	Value    float64       // Energy Cost or Power Limit
}
