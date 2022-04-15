package spine

import (
	"math"
	"strconv"
	"strings"
)

type ElementTagType struct{}

type LabelType string

type DescriptionType string

type SpecificationVersionType string

type TimePeriodType struct {
	StartTime *string `json:"startTime,omitempty"`
	EndTime   *string `json:"endTime,omitempty"`
}

type TimePeriodElementsType struct {
	StartTime *ElementTagType `json:"startTime,omitempty"`
	EndTime   *ElementTagType `json:"endTime,omitempty"`
}

type TimestampIntervalType struct {
	StartTime *string `json:"startTime,omitempty"`
	EndTime   *string `json:"endTime,omitempty"`
}

type AbsoluteOrRelativeTimeType string

type RecurringIntervalType RecurringIntervalEnumType

type RecurringIntervalEnumType string

const (
	RecurringIntervalEnumTypeYearly      RecurringIntervalEnumType = "yearly"
	RecurringIntervalEnumTypeMonthly     RecurringIntervalEnumType = "monthly"
	RecurringIntervalEnumTypeWeekly      RecurringIntervalEnumType = "weekly"
	RecurringIntervalEnumTypeDaily       RecurringIntervalEnumType = "daily"
	RecurringIntervalEnumTypeHourly      RecurringIntervalEnumType = "hourly"
	RecurringIntervalEnumTypeEveryminute RecurringIntervalEnumType = "everyMinute"
	RecurringIntervalEnumTypeEverysecond RecurringIntervalEnumType = "everySecond"
)

type MonthType string

const (
	MonthTypeJanuary   MonthType = "january"
	MonthTypeFebruary  MonthType = "february"
	MonthTypeMarch     MonthType = "march"
	MonthTypeApril     MonthType = "april"
	MonthTypeMay       MonthType = "may"
	MonthTypeJune      MonthType = "june"
	MonthTypeJuly      MonthType = "july"
	MonthTypeAugust    MonthType = "august"
	MonthTypeSeptember MonthType = "september"
	MonthTypeOctober   MonthType = "october"
	MonthTypeNovember  MonthType = "november"
	MonthTypeDecember  MonthType = "december"
)

type DayOfMonthType uint8

type CalendarWeekType uint8

type DayOfWeekType string

const (
	DayOfWeekTypeMonday    DayOfWeekType = "monday"
	DayOfWeekTypeTuesday   DayOfWeekType = "tuesday"
	DayOfWeekTypeWednesday DayOfWeekType = "wednesday"
	DayOfWeekTypeThursday  DayOfWeekType = "thursday"
	DayOfWeekTypeFriday    DayOfWeekType = "friday"
	DayOfWeekTypeSaturday  DayOfWeekType = "saturday"
	DayOfWeekTypeSunday    DayOfWeekType = "sunday"
)

type DaysOfWeekType struct {
	Monday    *ElementTagType `json:"monday,omitempty"`
	Tuesday   *ElementTagType `json:"tuesday,omitempty"`
	Wednesday *ElementTagType `json:"wednesday,omitempty"`
	Thursday  *ElementTagType `json:"thursday,omitempty"`
	Friday    *ElementTagType `json:"friday,omitempty"`
	Saturday  *ElementTagType `json:"saturday,omitempty"`
	Sunday    *ElementTagType `json:"sunday,omitempty"`
}

type OccurrenceType string

type OccurrenceEnumType string

const (
	OccurrenceEnumTypeFirst  OccurrenceEnumType = "first"
	OccurrenceEnumTypeSecond OccurrenceEnumType = "second"
	OccurrenceEnumTypeThird  OccurrenceEnumType = "third"
	OccurrenceEnumTypeFourth OccurrenceEnumType = "fourth"
	OccurrenceEnumTypeLast   OccurrenceEnumType = "last"
)

type AbsoluteOrRecurringTimeType struct {
	DateTime            *string           `json:"dateTime,omitempty"`
	Month               *MonthType        `json:"month,omitempty"`
	DayOfMonth          *DayOfMonthType   `json:"dayOfMonth,omitempty"`
	CalendarWeek        *CalendarWeekType `json:"calendarWeek,omitempty"`
	DayOfWeekOccurrence *OccurrenceType   `json:"dayOfWeekOccurrence,omitempty"`
	DaysOfWeek          *DaysOfWeekType   `json:"daysOfWeek,omitempty"`
	Time                *string           `json:"time,omitempty"`
	Relative            *string           `json:"relative,omitempty"`
}

type AbsoluteOrRecurringTimeElementsType struct {
	DateTime            *ElementTagType `json:"dateTime,omitempty"`
	Month               *ElementTagType `json:"month,omitempty"`
	DayOfMonth          *ElementTagType `json:"dayOfMonth,omitempty"`
	CalendarWeek        *ElementTagType `json:"calendarWeek,omitempty"`
	DayOfWeekOccurrence *ElementTagType `json:"dayOfWeekOccurrence,omitempty"`
	DaysOfWeek          *ElementTagType `json:"daysOfWeek,omitempty"`
	Time                *ElementTagType `json:"time,omitempty"`
	Relative            *ElementTagType `json:"relative,omitempty"`
}

type RecurrenceInformationType struct {
	RecurringInterval     *string `json:"recurringInterval,omitempty"`
	RecurringIntervalStep *uint   `json:"recurringIntervalStep,omitempty"`
	FirstExecution        *string `json:"firstExecution,omitempty"`
	ExecutionCount        *uint   `json:"executionCount,omitempty"`
	LastExecution         *string `json:"lastExecution,omitempty"`
}

type RecurrenceInformationElementsType struct {
	RecurringInterval     *ElementTagType `json:"recurringInterval,omitempty"`
	RecurringIntervalStep *ElementTagType `json:"recurringIntervalStep,omitempty"`
	FirstExecution        *ElementTagType `json:"firstExecution,omitempty"`
	ExecutionCount        *ElementTagType `json:"executionCount,omitempty"`
	LastExecution         *ElementTagType `json:"lastExecution,omitempty"`
}

type ScaledNumberRangeType struct {
	Min *ScaledNumberType `json:"min,omitempty"`
	Max *ScaledNumberType `json:"max,omitempty"`
}

type ScaledNumberRangeElementsType struct {
	Min *ElementTagType `json:"min,omitempty"`
	Max *ElementTagType `json:"max,omitempty"`
}

type ScaledNumberSetType struct {
	Value []ScaledNumberType      `json:"value,omitempty"`
	Range []ScaledNumberRangeType `json:"range,omitempty"`
}

type ScaledNumberSetElementsType struct {
	Value *ElementTagType `json:"value,omitempty"`
	Range *ElementTagType `json:"range,omitempty"`
}

type NumberType int64

type ScaleType int8

type ScaledNumberType struct {
	Number *NumberType `json:"number,omitempty"`
	Scale  *ScaleType  `json:"scale,omitempty"`
}

func (m ScaledNumberType) GetValue() float64 {
	if m.Number == nil {
		return 0
	}
	var scale float64 = 0
	if m.Scale != nil {
		scale = float64(*m.Scale)
	}
	return float64(*m.Number) * math.Pow(10, scale)
}

func NewScaledNumberType(value float64) *ScaledNumberType {
	m := &ScaledNumberType{}

	numberOfDecimals := 0
	temp := strconv.FormatFloat(value, 'f', -1, 64)
	index := strings.IndexByte(temp, '.')
	if index > -1 {
		numberOfDecimals = len(temp) - index - 1
	}

	if numberOfDecimals > 4 {
		numberOfDecimals = 4
	}

	numberValue := NumberType(math.Trunc(value * math.Pow(10, float64(numberOfDecimals))))
	m.Number = &numberValue

	if numberValue != 0 {
		scaleValue := ScaleType(-numberOfDecimals)
		m.Scale = &scaleValue
	}

	return m
}

type ScaledNumberElementsType struct {
	Number *ElementTagType `json:"number,omitempty"`
	Scale  *ElementTagType `json:"scale,omitempty"`
}

type MaxResponseDelayType string

type CommodityTypeType CommodityTypeEnumType

type CommodityTypeEnumType string

const (
	CommodityTypeEnumTypeElectricity      CommodityTypeEnumType = "electricity"
	CommodityTypeEnumTypeGas              CommodityTypeEnumType = "gas"
	CommodityTypeEnumTypeOil              CommodityTypeEnumType = "oil"
	CommodityTypeEnumTypeWater            CommodityTypeEnumType = "water"
	CommodityTypeEnumTypeWastewater       CommodityTypeEnumType = "wasteWater"
	CommodityTypeEnumTypeDomestichotwater CommodityTypeEnumType = "domesticHotWater"
	CommodityTypeEnumTypeHeatingwater     CommodityTypeEnumType = "heatingWater"
	CommodityTypeEnumTypeSteam            CommodityTypeEnumType = "steam"
	CommodityTypeEnumTypeHeat             CommodityTypeEnumType = "heat"
	CommodityTypeEnumTypeCoolingload      CommodityTypeEnumType = "coolingLoad"
	CommodityTypeEnumTypeAir              CommodityTypeEnumType = "air"
)

type EnergyDirectionType EnergyDirectionEnumType

type EnergyDirectionEnumType string

const (
	EnergyDirectionEnumTypeConsume EnergyDirectionEnumType = "consume"
	EnergyDirectionEnumTypeProduce EnergyDirectionEnumType = "produce"
)

type EnergyModeType EnergyModeEnumType

type EnergyModeEnumType string

const (
	EnergyModeEnumTypeConsume EnergyModeEnumType = "consume"
	EnergyModeEnumTypeProduce EnergyModeEnumType = "produce"
	EnergyModeEnumTypeIdle    EnergyModeEnumType = "idle"
	EnergyModeEnumTypeAuto    EnergyModeEnumType = "auto"
)

type UnitOfMeasurementType UnitOfMeasurementEnumType

type UnitOfMeasurementEnumType string

const (
	UnitOfMeasurementEnumTypeUnknown   UnitOfMeasurementEnumType = "unknown"
	UnitOfMeasurementEnumType1         UnitOfMeasurementEnumType = "1"
	UnitOfMeasurementEnumTypem         UnitOfMeasurementEnumType = "m"
	UnitOfMeasurementEnumTypekg        UnitOfMeasurementEnumType = "kg"
	UnitOfMeasurementEnumTypes         UnitOfMeasurementEnumType = "s"
	UnitOfMeasurementEnumTypeA         UnitOfMeasurementEnumType = "A"
	UnitOfMeasurementEnumTypeK         UnitOfMeasurementEnumType = "K"
	UnitOfMeasurementEnumTypemol       UnitOfMeasurementEnumType = "mol"
	UnitOfMeasurementEnumTypecd        UnitOfMeasurementEnumType = "cd"
	UnitOfMeasurementEnumTypeV         UnitOfMeasurementEnumType = "V"
	UnitOfMeasurementEnumTypeW         UnitOfMeasurementEnumType = "W"
	UnitOfMeasurementEnumTypeWh        UnitOfMeasurementEnumType = "Wh"
	UnitOfMeasurementEnumTypeVA        UnitOfMeasurementEnumType = "VA"
	UnitOfMeasurementEnumTypeVAh       UnitOfMeasurementEnumType = "VAh"
	UnitOfMeasurementEnumTypevar       UnitOfMeasurementEnumType = "var"
	UnitOfMeasurementEnumTypevarh      UnitOfMeasurementEnumType = "varh"
	UnitOfMeasurementEnumTypedegC      UnitOfMeasurementEnumType = "degC"
	UnitOfMeasurementEnumTypedegF      UnitOfMeasurementEnumType = "degF"
	UnitOfMeasurementEnumTypeLm        UnitOfMeasurementEnumType = "Lm"
	UnitOfMeasurementEnumTypelx        UnitOfMeasurementEnumType = "lx"
	UnitOfMeasurementEnumTypeOhm       UnitOfMeasurementEnumType = "Ohm"
	UnitOfMeasurementEnumTypeHz        UnitOfMeasurementEnumType = "Hz"
	UnitOfMeasurementEnumTypedB        UnitOfMeasurementEnumType = "dB"
	UnitOfMeasurementEnumTypedBm       UnitOfMeasurementEnumType = "dBm"
	UnitOfMeasurementEnumTypepct       UnitOfMeasurementEnumType = "pct"
	UnitOfMeasurementEnumTypeppm       UnitOfMeasurementEnumType = "ppm"
	UnitOfMeasurementEnumTypel         UnitOfMeasurementEnumType = "l"
	UnitOfMeasurementEnumTypels        UnitOfMeasurementEnumType = "l/s"
	UnitOfMeasurementEnumTypelh        UnitOfMeasurementEnumType = "l/h"
	UnitOfMeasurementEnumTypedeg       UnitOfMeasurementEnumType = "deg"
	UnitOfMeasurementEnumTyperad       UnitOfMeasurementEnumType = "rad"
	UnitOfMeasurementEnumTyperads      UnitOfMeasurementEnumType = "rad/s"
	UnitOfMeasurementEnumTypesr        UnitOfMeasurementEnumType = "sr"
	UnitOfMeasurementEnumTypeGy        UnitOfMeasurementEnumType = "Gy"
	UnitOfMeasurementEnumTypeBq        UnitOfMeasurementEnumType = "Bq"
	UnitOfMeasurementEnumTypeBqm3      UnitOfMeasurementEnumType = "Bq/m^3"
	UnitOfMeasurementEnumTypeSv        UnitOfMeasurementEnumType = "Sv"
	UnitOfMeasurementEnumTypeRd        UnitOfMeasurementEnumType = "Rd"
	UnitOfMeasurementEnumTypeC         UnitOfMeasurementEnumType = "C"
	UnitOfMeasurementEnumTypeF         UnitOfMeasurementEnumType = "F"
	UnitOfMeasurementEnumTypeH         UnitOfMeasurementEnumType = "H"
	UnitOfMeasurementEnumTypeJ         UnitOfMeasurementEnumType = "J"
	UnitOfMeasurementEnumTypeN         UnitOfMeasurementEnumType = "N"
	UnitOfMeasurementEnumTypeNm        UnitOfMeasurementEnumType = "N_m"
	UnitOfMeasurementEnumTypeNs        UnitOfMeasurementEnumType = "N_s"
	UnitOfMeasurementEnumTypeWb        UnitOfMeasurementEnumType = "Wb"
	UnitOfMeasurementEnumTypeT         UnitOfMeasurementEnumType = "T"
	UnitOfMeasurementEnumTypePa        UnitOfMeasurementEnumType = "Pa"
	UnitOfMeasurementEnumTypebar       UnitOfMeasurementEnumType = "bar"
	UnitOfMeasurementEnumTypeatm       UnitOfMeasurementEnumType = "atm"
	UnitOfMeasurementEnumTypepsi       UnitOfMeasurementEnumType = "psi"
	UnitOfMeasurementEnumTypemmHg      UnitOfMeasurementEnumType = "mmHg"
	UnitOfMeasurementEnumTypem2        UnitOfMeasurementEnumType = "m^2"
	UnitOfMeasurementEnumTypem3        UnitOfMeasurementEnumType = "m^3"
	UnitOfMeasurementEnumTypem3h       UnitOfMeasurementEnumType = "m^3/h"
	UnitOfMeasurementEnumTypems        UnitOfMeasurementEnumType = "m/s"
	UnitOfMeasurementEnumTypems2       UnitOfMeasurementEnumType = "m/s^2"
	UnitOfMeasurementEnumTypem3s       UnitOfMeasurementEnumType = "m^3/s"
	UnitOfMeasurementEnumTypemm3       UnitOfMeasurementEnumType = "m/m^3"
	UnitOfMeasurementEnumTypekgm3      UnitOfMeasurementEnumType = "kg/m^3"
	UnitOfMeasurementEnumTypekgm       UnitOfMeasurementEnumType = "kg_m"
	UnitOfMeasurementEnumTypem2s       UnitOfMeasurementEnumType = "m^2/s"
	UnitOfMeasurementEnumTypewmk       UnitOfMeasurementEnumType = "W/m_K"
	UnitOfMeasurementEnumTypeJK        UnitOfMeasurementEnumType = "J/K"
	UnitOfMeasurementEnumType1s        UnitOfMeasurementEnumType = "1/s"
	UnitOfMeasurementEnumTypeWm2       UnitOfMeasurementEnumType = "W/m^2"
	UnitOfMeasurementEnumTypeJm2       UnitOfMeasurementEnumType = "J/m^2"
	UnitOfMeasurementEnumTypeS         UnitOfMeasurementEnumType = "S"
	UnitOfMeasurementEnumTypeSm        UnitOfMeasurementEnumType = "S/m"
	UnitOfMeasurementEnumTypeKs        UnitOfMeasurementEnumType = "K/s"
	UnitOfMeasurementEnumTypePas       UnitOfMeasurementEnumType = "Pa/s"
	UnitOfMeasurementEnumTypeJkgK      UnitOfMeasurementEnumType = "J/kg_K"
	UnitOfMeasurementEnumTypeVs        UnitOfMeasurementEnumType = "Vs"
	UnitOfMeasurementEnumTypeVm        UnitOfMeasurementEnumType = "V/m"
	UnitOfMeasurementEnumTypeVHz       UnitOfMeasurementEnumType = "V/Hz"
	UnitOfMeasurementEnumTypeAs        UnitOfMeasurementEnumType = "As"
	UnitOfMeasurementEnumTypeAm        UnitOfMeasurementEnumType = "A/m"
	UnitOfMeasurementEnumTypeHzs       UnitOfMeasurementEnumType = "Hz/s"
	UnitOfMeasurementEnumTypekgs       UnitOfMeasurementEnumType = "kg/s"
	UnitOfMeasurementEnumTypekgm2      UnitOfMeasurementEnumType = "kg_m^2"
	UnitOfMeasurementEnumTypeJWh       UnitOfMeasurementEnumType = "J/Wh"
	UnitOfMeasurementEnumTypeWs        UnitOfMeasurementEnumType = "W/s"
	UnitOfMeasurementEnumTypeft3       UnitOfMeasurementEnumType = "ft^3"
	UnitOfMeasurementEnumTypeft3h      UnitOfMeasurementEnumType = "ft^3/h"
	UnitOfMeasurementEnumTypeccf       UnitOfMeasurementEnumType = "ccf"
	UnitOfMeasurementEnumTypeccfh      UnitOfMeasurementEnumType = "ccf/h"
	UnitOfMeasurementEnumTypeUSliqgal  UnitOfMeasurementEnumType = "US.liq.gal"
	UnitOfMeasurementEnumTypeUSliqgalh UnitOfMeasurementEnumType = "US.liq.gal/h"
	UnitOfMeasurementEnumTypeImpgal    UnitOfMeasurementEnumType = "Imp.gal"
	UnitOfMeasurementEnumTypeImpgalh   UnitOfMeasurementEnumType = "Imp.gal/h"
	UnitOfMeasurementEnumTypeBtu       UnitOfMeasurementEnumType = "Btu"
	UnitOfMeasurementEnumTypeBtuh      UnitOfMeasurementEnumType = "Btu/h"
	UnitOfMeasurementEnumTypeAh        UnitOfMeasurementEnumType = "Ah"
	UnitOfMeasurementEnumTypekgWh      UnitOfMeasurementEnumType = "kg/Wh"
)

type CurrencyType CurrencyEnumType

type CurrencyEnumType string

const (
	CurrencyEnumTypeAed CurrencyEnumType = "AED"
	CurrencyEnumTypeAfn CurrencyEnumType = "AFN"
	CurrencyEnumTypeAll CurrencyEnumType = "ALL"
	CurrencyEnumTypeAmd CurrencyEnumType = "AMD"
	CurrencyEnumTypeAng CurrencyEnumType = "ANG"
	CurrencyEnumTypeAoa CurrencyEnumType = "AOA"
	CurrencyEnumTypeArs CurrencyEnumType = "ARS"
	CurrencyEnumTypeAud CurrencyEnumType = "AUD"
	CurrencyEnumTypeAwg CurrencyEnumType = "AWG"
	CurrencyEnumTypeAzn CurrencyEnumType = "AZN"
	CurrencyEnumTypeBam CurrencyEnumType = "BAM"
	CurrencyEnumTypeBbd CurrencyEnumType = "BBD"
	CurrencyEnumTypeBdt CurrencyEnumType = "BDT"
	CurrencyEnumTypeBgn CurrencyEnumType = "BGN"
	CurrencyEnumTypeBhd CurrencyEnumType = "BHD"
	CurrencyEnumTypeBif CurrencyEnumType = "BIF"
	CurrencyEnumTypeBmd CurrencyEnumType = "BMD"
	CurrencyEnumTypeBnd CurrencyEnumType = "BND"
	CurrencyEnumTypeBob CurrencyEnumType = "BOB"
	CurrencyEnumTypeBov CurrencyEnumType = "BOV"
	CurrencyEnumTypeBrl CurrencyEnumType = "BRL"
	CurrencyEnumTypeBsd CurrencyEnumType = "BSD"
	CurrencyEnumTypeBtn CurrencyEnumType = "BTN"
	CurrencyEnumTypeBwp CurrencyEnumType = "BWP"
	CurrencyEnumTypeByr CurrencyEnumType = "BYR"
	CurrencyEnumTypeBzd CurrencyEnumType = "BZD"
	CurrencyEnumTypeCad CurrencyEnumType = "CAD"
	CurrencyEnumTypeCdf CurrencyEnumType = "CDF"
	CurrencyEnumTypeChe CurrencyEnumType = "CHE"
	CurrencyEnumTypeChf CurrencyEnumType = "CHF"
	CurrencyEnumTypeChw CurrencyEnumType = "CHW"
	CurrencyEnumTypeClf CurrencyEnumType = "CLF"
	CurrencyEnumTypeClp CurrencyEnumType = "CLP"
	CurrencyEnumTypeCny CurrencyEnumType = "CNY"
	CurrencyEnumTypeCop CurrencyEnumType = "COP"
	CurrencyEnumTypeCou CurrencyEnumType = "COU"
	CurrencyEnumTypeCrc CurrencyEnumType = "CRC"
	CurrencyEnumTypeCuc CurrencyEnumType = "CUC"
	CurrencyEnumTypeCup CurrencyEnumType = "CUP"
	CurrencyEnumTypeCve CurrencyEnumType = "CVE"
	CurrencyEnumTypeCzk CurrencyEnumType = "CZK"
	CurrencyEnumTypeDjf CurrencyEnumType = "DJF"
	CurrencyEnumTypeDkk CurrencyEnumType = "DKK"
	CurrencyEnumTypeDop CurrencyEnumType = "DOP"
	CurrencyEnumTypeDzd CurrencyEnumType = "DZD"
	CurrencyEnumTypeEgp CurrencyEnumType = "EGP"
	CurrencyEnumTypeErn CurrencyEnumType = "ERN"
	CurrencyEnumTypeEtb CurrencyEnumType = "ETB"
	CurrencyEnumTypeEur CurrencyEnumType = "EUR"
	CurrencyEnumTypeFjd CurrencyEnumType = "FJD"
	CurrencyEnumTypeFkp CurrencyEnumType = "FKP"
	CurrencyEnumTypeGbp CurrencyEnumType = "GBP"
	CurrencyEnumTypeGel CurrencyEnumType = "GEL"
	CurrencyEnumTypeGhs CurrencyEnumType = "GHS"
	CurrencyEnumTypeGip CurrencyEnumType = "GIP"
	CurrencyEnumTypeGmd CurrencyEnumType = "GMD"
	CurrencyEnumTypeGnf CurrencyEnumType = "GNF"
	CurrencyEnumTypeGtq CurrencyEnumType = "GTQ"
	CurrencyEnumTypeGyd CurrencyEnumType = "GYD"
	CurrencyEnumTypeHkd CurrencyEnumType = "HKD"
	CurrencyEnumTypeHnl CurrencyEnumType = "HNL"
	CurrencyEnumTypeHrk CurrencyEnumType = "HRK"
	CurrencyEnumTypeHtg CurrencyEnumType = "HTG"
	CurrencyEnumTypeHuf CurrencyEnumType = "HUF"
	CurrencyEnumTypeIdr CurrencyEnumType = "IDR"
	CurrencyEnumTypeIls CurrencyEnumType = "ILS"
	CurrencyEnumTypeInr CurrencyEnumType = "INR"
	CurrencyEnumTypeIqd CurrencyEnumType = "IQD"
	CurrencyEnumTypeIrr CurrencyEnumType = "IRR"
	CurrencyEnumTypeIsk CurrencyEnumType = "ISK"
	CurrencyEnumTypeJmd CurrencyEnumType = "JMD"
	CurrencyEnumTypeJod CurrencyEnumType = "JOD"
	CurrencyEnumTypeJpy CurrencyEnumType = "JPY"
	CurrencyEnumTypeKes CurrencyEnumType = "KES"
	CurrencyEnumTypeKgs CurrencyEnumType = "KGS"
	CurrencyEnumTypeKhr CurrencyEnumType = "KHR"
	CurrencyEnumTypeKmf CurrencyEnumType = "KMF"
	CurrencyEnumTypeKpw CurrencyEnumType = "KPW"
	CurrencyEnumTypeKrw CurrencyEnumType = "KRW"
	CurrencyEnumTypeKwd CurrencyEnumType = "KWD"
	CurrencyEnumTypeKyd CurrencyEnumType = "KYD"
	CurrencyEnumTypeKzt CurrencyEnumType = "KZT"
	CurrencyEnumTypeLak CurrencyEnumType = "LAK"
	CurrencyEnumTypeLbp CurrencyEnumType = "LBP"
	CurrencyEnumTypeLkr CurrencyEnumType = "LKR"
	CurrencyEnumTypeLrd CurrencyEnumType = "LRD"
	CurrencyEnumTypeLsl CurrencyEnumType = "LSL"
	CurrencyEnumTypeLyd CurrencyEnumType = "LYD"
	CurrencyEnumTypeMad CurrencyEnumType = "MAD"
	CurrencyEnumTypeMdl CurrencyEnumType = "MDL"
	CurrencyEnumTypeMga CurrencyEnumType = "MGA"
	CurrencyEnumTypeMkd CurrencyEnumType = "MKD"
	CurrencyEnumTypeMmk CurrencyEnumType = "MMK"
	CurrencyEnumTypeMnt CurrencyEnumType = "MNT"
	CurrencyEnumTypeMop CurrencyEnumType = "MOP"
	CurrencyEnumTypeMro CurrencyEnumType = "MRO"
	CurrencyEnumTypeMur CurrencyEnumType = "MUR"
	CurrencyEnumTypeMvr CurrencyEnumType = "MVR"
	CurrencyEnumTypeMwk CurrencyEnumType = "MWK"
	CurrencyEnumTypeMxn CurrencyEnumType = "MXN"
	CurrencyEnumTypeMxv CurrencyEnumType = "MXV"
	CurrencyEnumTypeMyr CurrencyEnumType = "MYR"
	CurrencyEnumTypeMzn CurrencyEnumType = "MZN"
	CurrencyEnumTypeNad CurrencyEnumType = "NAD"
	CurrencyEnumTypeNgn CurrencyEnumType = "NGN"
	CurrencyEnumTypeNio CurrencyEnumType = "NIO"
	CurrencyEnumTypeNok CurrencyEnumType = "NOK"
	CurrencyEnumTypeNpr CurrencyEnumType = "NPR"
	CurrencyEnumTypeNzd CurrencyEnumType = "NZD"
	CurrencyEnumTypeOmr CurrencyEnumType = "OMR"
	CurrencyEnumTypePab CurrencyEnumType = "PAB"
	CurrencyEnumTypePen CurrencyEnumType = "PEN"
	CurrencyEnumTypePgk CurrencyEnumType = "PGK"
	CurrencyEnumTypePhp CurrencyEnumType = "PHP"
	CurrencyEnumTypePkr CurrencyEnumType = "PKR"
	CurrencyEnumTypePln CurrencyEnumType = "PLN"
	CurrencyEnumTypePyg CurrencyEnumType = "PYG"
	CurrencyEnumTypeQar CurrencyEnumType = "QAR"
	CurrencyEnumTypeRon CurrencyEnumType = "RON"
	CurrencyEnumTypeRsd CurrencyEnumType = "RSD"
	CurrencyEnumTypeRub CurrencyEnumType = "RUB"
	CurrencyEnumTypeRwf CurrencyEnumType = "RWF"
	CurrencyEnumTypeSar CurrencyEnumType = "SAR"
	CurrencyEnumTypeSbd CurrencyEnumType = "SBD"
	CurrencyEnumTypeScr CurrencyEnumType = "SCR"
	CurrencyEnumTypeSdg CurrencyEnumType = "SDG"
	CurrencyEnumTypeSek CurrencyEnumType = "SEK"
	CurrencyEnumTypeSgd CurrencyEnumType = "SGD"
	CurrencyEnumTypeShp CurrencyEnumType = "SHP"
	CurrencyEnumTypeSll CurrencyEnumType = "SLL"
	CurrencyEnumTypeSos CurrencyEnumType = "SOS"
	CurrencyEnumTypeSrd CurrencyEnumType = "SRD"
	CurrencyEnumTypeSsp CurrencyEnumType = "SSP"
	CurrencyEnumTypeStd CurrencyEnumType = "STD"
	CurrencyEnumTypeSvc CurrencyEnumType = "SVC"
	CurrencyEnumTypeSyp CurrencyEnumType = "SYP"
	CurrencyEnumTypeSzl CurrencyEnumType = "SZL"
	CurrencyEnumTypeThb CurrencyEnumType = "THB"
	CurrencyEnumTypeTjs CurrencyEnumType = "TJS"
	CurrencyEnumTypeTmt CurrencyEnumType = "TMT"
	CurrencyEnumTypeTnd CurrencyEnumType = "TND"
	CurrencyEnumTypeTop CurrencyEnumType = "TOP"
	CurrencyEnumTypeTry CurrencyEnumType = "TRY"
	CurrencyEnumTypeTtd CurrencyEnumType = "TTD"
	CurrencyEnumTypeTwd CurrencyEnumType = "TWD"
	CurrencyEnumTypeTzs CurrencyEnumType = "TZS"
	CurrencyEnumTypeUah CurrencyEnumType = "UAH"
	CurrencyEnumTypeUgx CurrencyEnumType = "UGX"
	CurrencyEnumTypeUsd CurrencyEnumType = "USD"
	CurrencyEnumTypeUsn CurrencyEnumType = "USN"
	CurrencyEnumTypeUyi CurrencyEnumType = "UYI"
	CurrencyEnumTypeUyu CurrencyEnumType = "UYU"
	CurrencyEnumTypeUzs CurrencyEnumType = "UZS"
	CurrencyEnumTypeVef CurrencyEnumType = "VEF"
	CurrencyEnumTypeVnd CurrencyEnumType = "VND"
	CurrencyEnumTypeVuv CurrencyEnumType = "VUV"
	CurrencyEnumTypeWst CurrencyEnumType = "WST"
	CurrencyEnumTypeXaf CurrencyEnumType = "XAF"
	CurrencyEnumTypeXag CurrencyEnumType = "XAG"
	CurrencyEnumTypeXau CurrencyEnumType = "XAU"
	CurrencyEnumTypeXba CurrencyEnumType = "XBA"
	CurrencyEnumTypeXbb CurrencyEnumType = "XBB"
	CurrencyEnumTypeXbc CurrencyEnumType = "XBC"
	CurrencyEnumTypeXbd CurrencyEnumType = "XBD"
	CurrencyEnumTypeXcd CurrencyEnumType = "XCD"
	CurrencyEnumTypeXdr CurrencyEnumType = "XDR"
	CurrencyEnumTypeXof CurrencyEnumType = "XOF"
	CurrencyEnumTypeXpd CurrencyEnumType = "XPD"
	CurrencyEnumTypeXpf CurrencyEnumType = "XPF"
	CurrencyEnumTypeXpt CurrencyEnumType = "XPT"
	CurrencyEnumTypeXsu CurrencyEnumType = "XSU"
	CurrencyEnumTypeXts CurrencyEnumType = "XTS"
	CurrencyEnumTypeXua CurrencyEnumType = "XUA"
	CurrencyEnumTypeXxx CurrencyEnumType = "XXX"
	CurrencyEnumTypeYer CurrencyEnumType = "YER"
	CurrencyEnumTypeZar CurrencyEnumType = "ZAR"
	CurrencyEnumTypeZmw CurrencyEnumType = "ZMW"
	CurrencyEnumTypeZwl CurrencyEnumType = "ZWL"
)

type AddressDeviceType string

type AddressEntityType uint

type AddressFeatureType uint

type DeviceAddressType struct {
	Device *AddressDeviceType `json:"device,omitempty"`
}

type DeviceAddressElementsType struct {
	Device *ElementTagType `json:"device,omitempty"`
}

type EntityAddressType struct {
	Device *AddressDeviceType  `json:"device,omitempty"`
	Entity []AddressEntityType `json:"entity,omitempty"`
}

type EntityAddressElementsType struct {
	Device *ElementTagType `json:"device,omitempty"`
	Entity *ElementTagType `json:"entity,omitempty"`
}

type FeatureAddressType struct {
	Device  *AddressDeviceType  `json:"device,omitempty"`
	Entity  []AddressEntityType `json:"entity,omitempty"`
	Feature *AddressFeatureType `json:"feature,omitempty"`
}

type FeatureAddressElementsType struct {
	Device  *ElementTagType `json:"device,omitempty"`
	Entity  *ElementTagType `json:"entity,omitempty"`
	Feature *ElementTagType `json:"feature,omitempty"`
}

type ScopeTypeType ScopeTypeEnumType

type ScopeTypeEnumType string

const (
	ScopeTypeEnumTypeAC                    ScopeTypeEnumType = "ac"
	ScopeTypeEnumTypeACCosPhiGrid          ScopeTypeEnumType = "acCosPhiGrid"
	ScopeTypeEnumTypeACCurrentA            ScopeTypeEnumType = "acCurrentA"
	ScopeTypeEnumTypeACCurrentB            ScopeTypeEnumType = "acCurrentB"
	ScopeTypeEnumTypeACCurrentC            ScopeTypeEnumType = "acCurrentC"
	ScopeTypeEnumTypeACFrequencyGrid       ScopeTypeEnumType = "acFrequencyGrid"
	ScopeTypeEnumTypeACPowerA              ScopeTypeEnumType = "acPowerA"
	ScopeTypeEnumTypeACPowerB              ScopeTypeEnumType = "acPowerB"
	ScopeTypeEnumTypeACPowerC              ScopeTypeEnumType = "acPowerC"
	ScopeTypeEnumTypeACPowerLimitPct       ScopeTypeEnumType = "acPowerLimitPct"
	ScopeTypeEnumTypeACPowerTotal          ScopeTypeEnumType = "acPowerTotal"
	ScopeTypeEnumTypeACVoltageA            ScopeTypeEnumType = "acVoltageA"
	ScopeTypeEnumTypeACVoltageB            ScopeTypeEnumType = "acVoltageB"
	ScopeTypeEnumTypeACVoltageC            ScopeTypeEnumType = "acVoltageC"
	ScopeTypeEnumTypeACYieldDay            ScopeTypeEnumType = "acYieldDay"
	ScopeTypeEnumTypeACYieldTotal          ScopeTypeEnumType = "acYieldTotal"
	ScopeTypeEnumTypeDCCurrent             ScopeTypeEnumType = "dcCurrent"
	ScopeTypeEnumTypeDCPower               ScopeTypeEnumType = "dcPower"
	ScopeTypeEnumTypeDCString1             ScopeTypeEnumType = "dcString1"
	ScopeTypeEnumTypeDCString2             ScopeTypeEnumType = "dcString2"
	ScopeTypeEnumTypeDCString3             ScopeTypeEnumType = "dcString3"
	ScopeTypeEnumTypeDCString4             ScopeTypeEnumType = "dcString4"
	ScopeTypeEnumTypeDCString5             ScopeTypeEnumType = "dcString5"
	ScopeTypeEnumTypeDCString6             ScopeTypeEnumType = "dcString6"
	ScopeTypeEnumTypeDCTotal               ScopeTypeEnumType = "dcTotal"
	ScopeTypeEnumTypeDCVoltage             ScopeTypeEnumType = "dcVoltage"
	ScopeTypeEnumTypeDhwTemperature        ScopeTypeEnumType = "dhwTemperature"
	ScopeTypeEnumTypeFlowTemperature       ScopeTypeEnumType = "flowTemperature"
	ScopeTypeEnumTypeOutsideAirTemperature ScopeTypeEnumType = "outsideAirTemperature"
	ScopeTypeEnumTypeReturnTemperature     ScopeTypeEnumType = "returnTemperature"
	ScopeTypeEnumTypeRoomAirTemperature    ScopeTypeEnumType = "roomAirTemperature"
	ScopeTypeEnumTypeCharge                ScopeTypeEnumType = "charge"
	ScopeTypeEnumTypeStateOfCharge         ScopeTypeEnumType = "stateOfCharge"
	ScopeTypeEnumTypeDischarge             ScopeTypeEnumType = "discharge"
	ScopeTypeEnumTypeGridConsumption       ScopeTypeEnumType = "gridConsumption"
	ScopeTypeEnumTypeGridFeedIn            ScopeTypeEnumType = "gridFeedIn"
	ScopeTypeEnumTypeSelfConsumption       ScopeTypeEnumType = "selfConsumption"
	ScopeTypeEnumTypeOverloadProtection    ScopeTypeEnumType = "overloadProtection"
	ScopeTypeEnumTypeACPower               ScopeTypeEnumType = "acPower"
	ScopeTypeEnumTypeACEnergy              ScopeTypeEnumType = "acEnergy"
	ScopeTypeEnumTypeACCurrent             ScopeTypeEnumType = "acCurrent"
	ScopeTypeEnumTypeACVoltage             ScopeTypeEnumType = "acVoltage"
	ScopeTypeEnumTypeBatteryControl        ScopeTypeEnumType = "batteryControl"
	ScopeTypeEnumTypeSimpleIncentiveTable  ScopeTypeEnumType = "simpleIncentiveTable"
)

type RoleType string

const (
	RoleTypeClient  RoleType = "client"
	RoleTypeServer  RoleType = "server"
	RoleTypeSpecial RoleType = "special"
)

type FeatureGroupType string

type DeviceTypeType DeviceTypeEnumType

type DeviceTypeEnumType string

const (
	DeviceTypeEnumTypeDishwasher              DeviceTypeEnumType = "Dishwasher"
	DeviceTypeEnumTypeDryer                   DeviceTypeEnumType = "Dryer"
	DeviceTypeEnumTypeEnvironmentSensor       DeviceTypeEnumType = "EnvironmentSensor"
	DeviceTypeEnumTypeGeneric                 DeviceTypeEnumType = "Generic"
	DeviceTypeEnumTypeHeatgenerationSystem    DeviceTypeEnumType = "HeatGenerationSystem"
	DeviceTypeEnumTypeHeatsinkSystem          DeviceTypeEnumType = "HeatSinkSystem"
	DeviceTypeEnumTypeHeatstorageSystem       DeviceTypeEnumType = "HeatStorageSystem"
	DeviceTypeEnumTypeHVACController          DeviceTypeEnumType = "HVACController"
	DeviceTypeEnumTypeSubmeter                DeviceTypeEnumType = "SubMeter"
	DeviceTypeEnumTypeWasher                  DeviceTypeEnumType = "Washer"
	DeviceTypeEnumTypeElectricitySupplySystem DeviceTypeEnumType = "ElectricitySupplySystem"
	DeviceTypeEnumTypeEnergyManagementSystem  DeviceTypeEnumType = "EnergyManagementSystem"
	DeviceTypeEnumTypeInverter                DeviceTypeEnumType = "Inverter"
	DeviceTypeEnumTypeChargingStation         DeviceTypeEnumType = "ChargingStation"
)

type EntityTypeType EntityTypeEnumType

type EntityTypeEnumType string

const (
	EntityTypeEnumTypeBattery                       EntityTypeEnumType = "Battery"
	EntityTypeEnumTypeCompressor                    EntityTypeEnumType = "Compressor"
	EntityTypeEnumTypeDeviceInformation             EntityTypeEnumType = "DeviceInformation"
	EntityTypeEnumTypeDHWCircuit                    EntityTypeEnumType = "DHWCircuit"
	EntityTypeEnumTypeDHWStorage                    EntityTypeEnumType = "DHWStorage"
	EntityTypeEnumTypeDishwasher                    EntityTypeEnumType = "Dishwasher"
	EntityTypeEnumTypeDryer                         EntityTypeEnumType = "Dryer"
	EntityTypeEnumTypeElectricalImmersionheater     EntityTypeEnumType = "ElectricalImmersionHeater"
	EntityTypeEnumTypeFan                           EntityTypeEnumType = "Fan"
	EntityTypeEnumTypeGasHeatingAppliance           EntityTypeEnumType = "GasHeatingAppliance"
	EntityTypeEnumTypeGeneric                       EntityTypeEnumType = "Generic"
	EntityTypeEnumTypeHeatingBufferStorage          EntityTypeEnumType = "HeatingBufferStorage"
	EntityTypeEnumTypeHeatingCircuit                EntityTypeEnumType = "HeatingCircuit"
	EntityTypeEnumTypeHeatingObject                 EntityTypeEnumType = "HeatingObject"
	EntityTypeEnumTypeHeatingZone                   EntityTypeEnumType = "HeatingZone"
	EntityTypeEnumTypeHeatPumpAppliance             EntityTypeEnumType = "HeatPumpAppliance"
	EntityTypeEnumTypeHeatSinkCircuit               EntityTypeEnumType = "HeatSinkCircuit"
	EntityTypeEnumTypeHeatSourceCircuit             EntityTypeEnumType = "HeatSourceCircuit"
	EntityTypeEnumTypeHeatSourceUnit                EntityTypeEnumType = "HeatSourceUnit"
	EntityTypeEnumTypeHvacController                EntityTypeEnumType = "HVACController"
	EntityTypeEnumTypeHvacRoom                      EntityTypeEnumType = "HVACRoom"
	EntityTypeEnumTypeInstantDHWheater              EntityTypeEnumType = "InstantDHWHeater"
	EntityTypeEnumTypeInverter                      EntityTypeEnumType = "Inverter"
	EntityTypeEnumTypeOilHeatingAppliance           EntityTypeEnumType = "OilHeatingAppliance"
	EntityTypeEnumTypePump                          EntityTypeEnumType = "Pump"
	EntityTypeEnumTypeRefrigerantCircuit            EntityTypeEnumType = "RefrigerantCircuit"
	EntityTypeEnumTypeSmartEnergyAppliance          EntityTypeEnumType = "SmartEnergyAppliance"
	EntityTypeEnumTypeSolarDHWStorage               EntityTypeEnumType = "SolarDHWStorage"
	EntityTypeEnumTypeSolarThermalCircuit           EntityTypeEnumType = "SolarThermalCircuit"
	EntityTypeEnumTypeSubmeterElectricity           EntityTypeEnumType = "SubMeterElectricity"
	EntityTypeEnumTypeTemperatureSensor             EntityTypeEnumType = "TemperatureSensor"
	EntityTypeEnumTypeWasher                        EntityTypeEnumType = "Washer"
	EntityTypeEnumTypeBatterySystem                 EntityTypeEnumType = "BatterySystem"
	EntityTypeEnumTypeElectricityGenerationSystem   EntityTypeEnumType = "ElectricityGenerationSystem"
	EntityTypeEnumTypeElectricityStorageSystem      EntityTypeEnumType = "ElectricityStorageSystem"
	EntityTypeEnumTypeGridConnectionPointOfPremises EntityTypeEnumType = "GridConnectionPointOfPremises"
	EntityTypeEnumTypeHousehold                     EntityTypeEnumType = "Household"
	EntityTypeEnumTypePVSystem                      EntityTypeEnumType = "PVSystem"
	EntityTypeEnumTypeEV                            EntityTypeEnumType = "EV"
	EntityTypeEnumTypeEVSE                          EntityTypeEnumType = "EVSE"
	EntityTypeEnumTypeChargingOutlet                EntityTypeEnumType = "ChargingOutlet"
	EntityTypeEnumTypeCEM                           EntityTypeEnumType = "CEM"
)

type FeatureTypeType FeatureTypeEnumType

type FeatureTypeEnumType string

const (
	FeatureTypeEnumTypeActuatorLevel           FeatureTypeEnumType = "ActuatorLevel"
	FeatureTypeEnumTypeActuatorSwitch          FeatureTypeEnumType = "ActuatorSwitch"
	FeatureTypeEnumTypeAlarm                   FeatureTypeEnumType = "Alarm"
	FeatureTypeEnumTypeDataTunneling           FeatureTypeEnumType = "DataTunneling"
	FeatureTypeEnumTypeDeviceClassification    FeatureTypeEnumType = "DeviceClassification"
	FeatureTypeEnumTypeDeviceDiagnosis         FeatureTypeEnumType = "DeviceDiagnosis"
	FeatureTypeEnumTypeDirectControl           FeatureTypeEnumType = "DirectControl"
	FeatureTypeEnumTypeElectricalConnection    FeatureTypeEnumType = "ElectricalConnection"
	FeatureTypeEnumTypeGeneric                 FeatureTypeEnumType = "Generic"
	FeatureTypeEnumTypeHvac                    FeatureTypeEnumType = "HVAC"
	FeatureTypeEnumTypeLoadControl             FeatureTypeEnumType = "LoadControl"
	FeatureTypeEnumTypeMeasurement             FeatureTypeEnumType = "Measurement"
	FeatureTypeEnumTypeMessaging               FeatureTypeEnumType = "Messaging"
	FeatureTypeEnumTypeNetworkManagement       FeatureTypeEnumType = "NetworkManagement"
	FeatureTypeEnumTypeNodeManagement          FeatureTypeEnumType = "NodeManagement"
	FeatureTypeEnumTypeOperatingConstraints    FeatureTypeEnumType = "OperatingConstraints"
	FeatureTypeEnumTypePowerSequences          FeatureTypeEnumType = "PowerSequences"
	FeatureTypeEnumTypeSensing                 FeatureTypeEnumType = "Sensing"
	FeatureTypeEnumTypeSetpoint                FeatureTypeEnumType = "Setpoint"
	FeatureTypeEnumTypeSmartEnergyManagementPs FeatureTypeEnumType = "SmartEnergyManagementPs"
	FeatureTypeEnumTypeTaskManagement          FeatureTypeEnumType = "TaskManagement"
	FeatureTypeEnumTypeThreshold               FeatureTypeEnumType = "Threshold"
	FeatureTypeEnumTypeTimeInformation         FeatureTypeEnumType = "TimeInformation"
	FeatureTypeEnumTypeTimeTable               FeatureTypeEnumType = "TimeTable"
	FeatureTypeEnumTypeDeviceConfiguration     FeatureTypeEnumType = "DeviceConfiguration"
	FeatureTypeEnumTypeSupplyCondition         FeatureTypeEnumType = "SupplyCondition"
	FeatureTypeEnumTypeTimeSeries              FeatureTypeEnumType = "TimeSeries"
	FeatureTypeEnumTypeTariffInformation       FeatureTypeEnumType = "TariffInformation"
	FeatureTypeEnumTypeIncentiveTable          FeatureTypeEnumType = "IncentiveTable"
	FeatureTypeEnumTypeBill                    FeatureTypeEnumType = "Bill"
	FeatureTypeEnumTypeIdentification          FeatureTypeEnumType = "Identification"
)

type FeatureSpecificUsageType FeatureSpecificUsageEnumType

type FeatureSpecificUsageEnumType string

const (
// FeatureDirectControlSpecificUsageEnumType
// FeatureHvacSpecificUsageEnumType
// FeatureMeasurementSpecificUsageEnumType
// FeatureSetpointSpecificUsageEnumType
// FeatureSmartEnergyManagementPsSpecificUsageEnumType
)

type FeatureDirectControlSpecificUsageEnumType string

const (
	FeatureDirectControlSpecificUsageEnumTypeHistory  FeatureDirectControlSpecificUsageEnumType = "History"
	FeatureDirectControlSpecificUsageEnumTypeRealtime FeatureDirectControlSpecificUsageEnumType = "RealTime"
)

type FeatureHvacSpecificUsageEnumType string

const (
	FeatureHvacSpecificUsageEnumTypeOperationmode FeatureHvacSpecificUsageEnumType = "OperationMode"
	FeatureHvacSpecificUsageEnumTypeOverrun       FeatureHvacSpecificUsageEnumType = "Overrun"
)

type FeatureMeasurementSpecificUsageEnumType string

const (
	FeatureMeasurementSpecificUsageEnumTypeContact     FeatureMeasurementSpecificUsageEnumType = "Contact"
	FeatureMeasurementSpecificUsageEnumTypeElectrical  FeatureMeasurementSpecificUsageEnumType = "Electrical"
	FeatureMeasurementSpecificUsageEnumTypeHeat        FeatureMeasurementSpecificUsageEnumType = "Heat"
	FeatureMeasurementSpecificUsageEnumTypeLevel       FeatureMeasurementSpecificUsageEnumType = "Level"
	FeatureMeasurementSpecificUsageEnumTypePressure    FeatureMeasurementSpecificUsageEnumType = "Pressure"
	FeatureMeasurementSpecificUsageEnumTypeTemperature FeatureMeasurementSpecificUsageEnumType = "Temperature"
)

type FeatureSetpointSpecificUsageEnumType string

type FeatureSmartEnergyManagementPsSpecificUsageEnumType string

const (
	FeatureSmartEnergyManagementPsSpecificUsageEnumTypeFixedForecast                         FeatureSmartEnergyManagementPsSpecificUsageEnumType = "FixedForecast"
	FeatureSmartEnergyManagementPsSpecificUsageEnumTypeFlexibleChosenForecast                FeatureSmartEnergyManagementPsSpecificUsageEnumType = "FlexibleChosenForecast"
	FeatureSmartEnergyManagementPsSpecificUsageEnumTypeFlexibleOptionalForecast              FeatureSmartEnergyManagementPsSpecificUsageEnumType = "FlexibleOptionalForecast"
	FeatureSmartEnergyManagementPsSpecificUsageEnumTypeOptionalSequenceBasedImmediateControl FeatureSmartEnergyManagementPsSpecificUsageEnumType = "OptionalSequenceBasedImmediateControl"
)

type FunctionType FunctionEnumType

type FunctionEnumType string

const (
	FunctionEnumTypeActuatorLevelData                                  FunctionEnumType = "actuatorLevelData"
	FunctionEnumTypeActuatorLevelDescriptionData                       FunctionEnumType = "actuatorLevelDescriptionData"
	FunctionEnumTypeActuatorSwitchData                                 FunctionEnumType = "actuatorSwitchData"
	FunctionEnumTypeActuatorSwitchDescriptionData                      FunctionEnumType = "actuatorSwitchDescriptionData"
	FunctionEnumTypeAlarmListData                                      FunctionEnumType = "alarmListData"
	FunctionEnumTypeBindingManagementDeleteCall                        FunctionEnumType = "bindingManagementDeleteCall"
	FunctionEnumTypeBindingManagementEntryListData                     FunctionEnumType = "bindingManagementEntryListData"
	FunctionEnumTypeBindingManagementRequestCall                       FunctionEnumType = "bindingManagementRequestCall"
	FunctionEnumTypeDataTunnelingCall                                  FunctionEnumType = "dataTunnelingCall"
	FunctionEnumTypeDeviceClassificationManufacturerData               FunctionEnumType = "deviceClassificationManufacturerData"
	FunctionEnumTypeDeviceClassificationUserData                       FunctionEnumType = "deviceClassificationUserData"
	FunctionEnumTypeDeviceDiagnosisHeartbeatData                       FunctionEnumType = "deviceDiagnosisHeartbeatData"
	FunctionEnumTypeDeviceDiagnosisServiceData                         FunctionEnumType = "deviceDiagnosisServiceData"
	FunctionEnumTypeDeviceDiagnosisStateData                           FunctionEnumType = "deviceDiagnosisStateData"
	FunctionEnumTypeDirectControlActivityListData                      FunctionEnumType = "directControlActivityListData"
	FunctionEnumTypeDirectControlDescriptionData                       FunctionEnumType = "directControlDescriptionData"
	FunctionEnumTypeElectricalConnectionDescriptionListData            FunctionEnumType = "electricalConnectionDescriptionListData"
	FunctionEnumTypeElectricalConnectionParameterDescriptionListData   FunctionEnumType = "electricalConnectionParameterDescriptionListData"
	FunctionEnumTypeElectricalConnectionStateListData                  FunctionEnumType = "electricalConnectionStateListData"
	FunctionEnumTypeHvacOperationModeDescriptionListData               FunctionEnumType = "hvacOperationModeDescriptionListData"
	FunctionEnumTypeHvacOverrunDescriptionListData                     FunctionEnumType = "hvacOverrunDescriptionListData"
	FunctionEnumTypeHvacOverrunListData                                FunctionEnumType = "hvacOverrunListData"
	FunctionEnumTypeHvacSystemFunctionDescriptionListData              FunctionEnumType = "hvacSystemFunctionDescriptionListData"
	FunctionEnumTypeHvacSystemFunctionListData                         FunctionEnumType = "hvacSystemFunctionListData"
	FunctionEnumTypeHvacSystemFunctionOperationModeRelationListData    FunctionEnumType = "hvacSystemFunctionOperationModeRelationListData"
	FunctionEnumTypeHvacSystemFunctionPowerSequenceRelationListData    FunctionEnumType = "hvacSystemFunctionPowerSequenceRelationListData"
	FunctionEnumTypeHvacSystemFunctionSetPointRelationListData         FunctionEnumType = "hvacSystemFunctionSetpointRelationListData"
	FunctionEnumTypeLoadControlEventListData                           FunctionEnumType = "loadControlEventListData"
	FunctionEnumTypeLoadControlStateListData                           FunctionEnumType = "loadControlStateListData"
	FunctionEnumTypeMeasurementConstraintsListData                     FunctionEnumType = "measurementConstraintsListData"
	FunctionEnumTypeMeasurementDescriptionListData                     FunctionEnumType = "measurementDescriptionListData"
	FunctionEnumTypeMeasurementListData                                FunctionEnumType = "measurementListData"
	FunctionEnumTypeMeasurementThresholdRelationListData               FunctionEnumType = "measurementThresholdRelationListData"
	FunctionEnumTypeMessagingListData                                  FunctionEnumType = "messagingListData"
	FunctionEnumTypeNetworkManagementAbortCall                         FunctionEnumType = "networkManagementAbortCall"
	FunctionEnumTypeNetworkManagementAddnodeCall                       FunctionEnumType = "networkManagementAddNodeCall"
	FunctionEnumTypeNetworkManagementDeviceDescriptionListData         FunctionEnumType = "networkManagementDeviceDescriptionListData"
	FunctionEnumTypeNetworkManagementDiscoverCall                      FunctionEnumType = "networkManagementDiscoverCall"
	FunctionEnumTypeNetworkManagementEntityDescriptionListData         FunctionEnumType = "networkManagementEntityDescriptionListData"
	FunctionEnumTypeNetworkManagementFeatureDescriptionListData        FunctionEnumType = "networkManagementFeatureDescriptionListData"
	FunctionEnumTypeNetworkManagementJoiningModeData                   FunctionEnumType = "networkManagementJoiningModeData"
	FunctionEnumTypeNetworkManagementModifyNodeCall                    FunctionEnumType = "networkManagementModifyNodeCall"
	FunctionEnumTypeNetworkManagementProcessStateData                  FunctionEnumType = "networkManagementProcessStateData"
	FunctionEnumTypeNetworkManagementRemoveNodeCall                    FunctionEnumType = "networkManagementRemoveNodeCall"
	FunctionEnumTypeNetworkManagementReportCandidateData               FunctionEnumType = "networkManagementReportCandidateData"
	FunctionEnumTypeNetworkManagementScanNetworkCall                   FunctionEnumType = "networkManagementScanNetworkCall"
	FunctionEnumTypeNodeManagementBindingData                          FunctionEnumType = "nodeManagementBindingData"
	FunctionEnumTypeNodeManagementBindingDeleteCall                    FunctionEnumType = "nodeManagementBindingDeleteCall"
	FunctionEnumTypeNodeManagementBindingRequestCall                   FunctionEnumType = "nodeManagementBindingRequestCall"
	FunctionEnumTypeNodeManagementDestinationListData                  FunctionEnumType = "nodeManagementDestinationListData"
	FunctionEnumTypeNodeManagementDetailedDiscoveryData                FunctionEnumType = "nodeManagementDetailedDiscoveryData"
	FunctionEnumTypeNodeManagementSubscriptionData                     FunctionEnumType = "nodeManagementSubscriptionData"
	FunctionEnumTypeNodeManagementSubscriptionDeleteCall               FunctionEnumType = "nodeManagementSubscriptionDeleteCall"
	FunctionEnumTypeNodeManagementSubscriptionRequestCall              FunctionEnumType = "nodeManagementSubscriptionRequestCall"
	FunctionEnumTypeOperatingConstraintsDurationListData               FunctionEnumType = "operatingConstraintsDurationListData"
	FunctionEnumTypeOperatingConstraintsInterruptListData              FunctionEnumType = "operatingConstraintsInterruptListData"
	FunctionEnumTypeOperatingConstraintsPowerDescriptionListData       FunctionEnumType = "operatingConstraintsPowerDescriptionListData"
	FunctionEnumTypeOperatingConstraintsPowerLevelListData             FunctionEnumType = "operatingConstraintsPowerLevelListData"
	FunctionEnumTypeOperatingConstraintsPowerRangeListData             FunctionEnumType = "operatingConstraintsPowerRangeListData"
	FunctionEnumTypeOperatingConstraintsResumeImplicationListData      FunctionEnumType = "operatingConstraintsResumeImplicationListData"
	FunctionEnumTypePowerSequenceAlternativesRelationlistdata          FunctionEnumType = "powerSequenceAlternativesRelationListData"
	FunctionEnumTypePowerSequenceDescriptionListData                   FunctionEnumType = "powerSequenceDescriptionListData"
	FunctionEnumTypePowerSequenceNodeScheduleInformationData           FunctionEnumType = "powerSequenceNodeScheduleInformationData"
	FunctionEnumTypePowerSequencePriceCalculationRequestCall           FunctionEnumType = "powerSequencePriceCalculationRequestCall"
	FunctionEnumTypePowerSequencePriceListData                         FunctionEnumType = "powerSequencePriceListData"
	FunctionEnumTypePowerSequenceScheduleConfigurationRequestCall      FunctionEnumType = "powerSequenceScheduleConfigurationRequestCall"
	FunctionEnumTypePowerSequenceScheduleConstraintsListData           FunctionEnumType = "powerSequenceScheduleConstraintsListData"
	FunctionEnumTypePowerSequenceScheduleListData                      FunctionEnumType = "powerSequenceScheduleListData"
	FunctionEnumTypePowerSequenceSchedulePreferenceListData            FunctionEnumType = "powerSequenceSchedulePreferenceListData"
	FunctionEnumTypePowerSequenceStateListData                         FunctionEnumType = "powerSequenceStateListData"
	FunctionEnumTypePowerTimeslotScheduleConstraintsListData           FunctionEnumType = "powerTimeSlotScheduleConstraintsListData"
	FunctionEnumTypePowerTimeslotScheduleListData                      FunctionEnumType = "powerTimeSlotScheduleListData"
	FunctionEnumTypePowerTimeslotValueListData                         FunctionEnumType = "powerTimeSlotValueListData"
	FunctionEnumTypeResultData                                         FunctionEnumType = "resultData"
	FunctionEnumTypeSensingDescriptionData                             FunctionEnumType = "sensingDescriptionData"
	FunctionEnumTypeSensingListData                                    FunctionEnumType = "sensingListData"
	FunctionEnumTypeSetPointConstraintsListData                        FunctionEnumType = "setpointConstraintsListData"
	FunctionEnumTypeSetPointDescriptionListData                        FunctionEnumType = "setpointDescriptionListData"
	FunctionEnumTypeSetPointListData                                   FunctionEnumType = "setpointListData"
	FunctionEnumTypeSmartEnergyManagementPsConfigurationRequestCall    FunctionEnumType = "smartEnergyManagementPsConfigurationRequestCall"
	FunctionEnumTypeSmartEnergyManagementPsData                        FunctionEnumType = "smartEnergyManagementPsData"
	FunctionEnumTypeSmartEnergyManagementPsPriceCalculationRequestCall FunctionEnumType = "smartEnergyManagementPsPriceCalculationRequestCall"
	FunctionEnumTypeSmartEnergyManagementPsPriceData                   FunctionEnumType = "smartEnergyManagementPsPriceData"
	FunctionEnumTypeSpecificationVersionListData                       FunctionEnumType = "specificationVersionListData"
	FunctionEnumTypeSubscriptionManagementDeleteCall                   FunctionEnumType = "subscriptionManagementDeleteCall"
	FunctionEnumTypeSubscriptionManagementEntryListData                FunctionEnumType = "subscriptionManagementEntryListData"
	FunctionEnumTypeSubscriptionManagementRequestCall                  FunctionEnumType = "subscriptionManagementRequestCall"
	FunctionEnumTypeSupplyConditionDescriptionListData                 FunctionEnumType = "supplyConditionDescriptionListData"
	FunctionEnumTypeSupplyConditionListData                            FunctionEnumType = "supplyConditionListData"
	FunctionEnumTypeSupplyConditionThresholdRelationListData           FunctionEnumType = "supplyConditionThresholdRelationListData"
	FunctionEnumTypeTaskManagementJobDescriptionListData               FunctionEnumType = "taskManagementJobDescriptionListData"
	FunctionEnumTypeTaskManagementJobListData                          FunctionEnumType = "taskManagementJobListData"
	FunctionEnumTypeTaskManagementJobRelationListData                  FunctionEnumType = "taskManagementJobRelationListData"
	FunctionEnumTypeTaskManagementOverviewData                         FunctionEnumType = "taskManagementOverviewData"
	FunctionEnumTypeThresholdConstraintsListData                       FunctionEnumType = "thresholdConstraintsListData"
	FunctionEnumTypeThresholdDescriptionListData                       FunctionEnumType = "thresholdDescriptionListData"
	FunctionEnumTypeThresholdListData                                  FunctionEnumType = "thresholdListData"
	FunctionEnumTypeTimeDistributorData                                FunctionEnumType = "timeDistributorData"
	FunctionEnumTypeTimeDistributorEnquiryCall                         FunctionEnumType = "timeDistributorEnquiryCall"
	FunctionEnumTypeTimeInformationData                                FunctionEnumType = "timeInformationData"
	FunctionEnumTypeTimePrecisionData                                  FunctionEnumType = "timePrecisionData"
	FunctionEnumTypeTimeTableConstraintsListData                       FunctionEnumType = "timeTableConstraintsListData"
	FunctionEnumTypeTimeTableDescriptionListData                       FunctionEnumType = "timeTableDescriptionListData"
	FunctionEnumTypeTimeTableListData                                  FunctionEnumType = "timeTableListData"
	FunctionEnumTypeDeviceConfigurationKeyValueConstraintsListData     FunctionEnumType = "deviceConfigurationKeyValueConstraintsListData"
	FunctionEnumTypeDeviceConfigurationKeyValueListData                FunctionEnumType = "deviceConfigurationKeyValueListData"
	FunctionEnumTypeDeviceConfigurationKeyValueDescriptionListData     FunctionEnumType = "deviceConfigurationKeyValueDescriptionListData"
	FunctionEnumTypeLoadControlLimitConstraintsListData                FunctionEnumType = "loadControlLimitConstraintsListData"
	FunctionEnumTypeLoadControlLimitDescriptionListData                FunctionEnumType = "loadControlLimitDescriptionListData"
	FunctionEnumTypeLoadControlLimitListData                           FunctionEnumType = "loadControlLimitListData"
	FunctionEnumTypeLoadControlNodeData                                FunctionEnumType = "loadControlNodeData"
	FunctionEnumTypeTimeSeriesConstraintsListData                      FunctionEnumType = "timeSeriesConstraintsListData"
	FunctionEnumTypeTimeSeriesDescriptionListData                      FunctionEnumType = "timeSeriesDescriptionListData"
	FunctionEnumTypeTimeSeriesListData                                 FunctionEnumType = "timeSeriesListData"
	FunctionEnumTypeTariffOverallConstraintsData                       FunctionEnumType = "tariffOverallConstraintsData"
	FunctionEnumTypeTariffListData                                     FunctionEnumType = "tariffListData"
	FunctionEnumTypeTariffBoundaryrelationListData                     FunctionEnumType = "tariffBoundaryRelationListData"
	FunctionEnumTypeTariffTierRelationListData                         FunctionEnumType = "tariffTierRelationListData"
	FunctionEnumTypeTariffDescriptionListData                          FunctionEnumType = "tariffDescriptionListData"
	FunctionEnumTypeTierBoundaryListData                               FunctionEnumType = "tierBoundaryListData"
	FunctionEnumTypeTierBoundaryDescriptionListData                    FunctionEnumType = "tierBoundaryDescriptionListData"
	FunctionEnumTypeCommodityListData                                  FunctionEnumType = "commodityListData"
	FunctionEnumTypeTierListData                                       FunctionEnumType = "tierListData"
	FunctionEnumTypeTierIncentiveRelationListData                      FunctionEnumType = "tierIncentiveRelationListData"
	FunctionEnumTypeTierDescriptionListData                            FunctionEnumType = "tierDescriptionListData"
	FunctionEnumTypeIncentiveListData                                  FunctionEnumType = "incentiveListData"
	FunctionEnumTypeIncentiveDescriptionListData                       FunctionEnumType = "incentiveDescriptionListData"
	FunctionEnumTypeIncentiveTableData                                 FunctionEnumType = "incentiveTableData"
	FunctionEnumTypeIncentiveTableDescriptionData                      FunctionEnumType = "incentiveTableDescriptionData"
	FunctionEnumTypeIncentiveTableConstraintsData                      FunctionEnumType = "incentiveTableConstraintsData"
	FunctionEnumTypeElectricalConnectionPermittedValueSetListData      FunctionEnumType = "electricalConnectionPermittedValueSetListData"
	FunctionEnumTypeUseCaseInformationListData                         FunctionEnumType = "useCaseInformationListData"
	FunctionEnumTypeNodeManagementUseCaseData                          FunctionEnumType = "nodeManagementUseCaseData"
	FunctionEnumTypeBillConstraintsListData                            FunctionEnumType = "billConstraintsListData"
	FunctionEnumTypeBillDescriptionListData                            FunctionEnumType = "billDescriptionListData"
	FunctionEnumTypeBillListData                                       FunctionEnumType = "billListData"
	FunctionEnumTypeIdentificationListData                             FunctionEnumType = "identificationListData"
)

type PossibleOperationsClassifierType struct {
	Partial *ElementTagType `json:"partial,omitempty"`
}

type PossibleOperationsReadType struct {
	Partial *ElementTagType `json:"partial,omitempty"`
}

type PossibleOperationsWriteType struct {
	Partial *ElementTagType `json:"partial,omitempty"`
}

type PossibleOperationsType struct {
	Read  *PossibleOperationsReadType  `json:"read,omitempty"`
	Write *PossibleOperationsWriteType `json:"write,omitempty"`
}

type PossibleOperationsElementsType struct {
	Read  *ElementTagType `json:"read,omitempty"`
	Write *ElementTagType `json:"write,omitempty"`
}

type FunctionPropertyType struct {
	Function           *FunctionType           `json:"function,omitempty"`
	PossibleOperations *PossibleOperationsType `json:"possibleOperations,omitempty"`
}

type FunctionPropertyElementsType struct {
	Function           *ElementTagType `json:"function,omitempty"`
	PossibleOperations *ElementTagType `json:"possibleOperations,omitempty"`
}
