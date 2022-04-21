package model

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

type RecurringIntervalType string

const (
	RecurringIntervalTypeYearly      RecurringIntervalType = "yearly"
	RecurringIntervalTypeMonthly     RecurringIntervalType = "monthly"
	RecurringIntervalTypeWeekly      RecurringIntervalType = "weekly"
	RecurringIntervalTypeDaily       RecurringIntervalType = "daily"
	RecurringIntervalTypeHourly      RecurringIntervalType = "hourly"
	RecurringIntervalTypeEveryminute RecurringIntervalType = "everyMinute"
	RecurringIntervalTypeEverysecond RecurringIntervalType = "everySecond"
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

const (
	OccurrenceTypeFirst  OccurrenceType = "first"
	OccurrenceTypeSecond OccurrenceType = "second"
	OccurrenceTypeThird  OccurrenceType = "third"
	OccurrenceTypeFourth OccurrenceType = "fourth"
	OccurrenceTypeLast   OccurrenceType = "last"
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

type CommodityTypeType string

const (
	CommodityTypeTypeElectricity      CommodityTypeType = "electricity"
	CommodityTypeTypeGas              CommodityTypeType = "gas"
	CommodityTypeTypeOil              CommodityTypeType = "oil"
	CommodityTypeTypeWater            CommodityTypeType = "water"
	CommodityTypeTypeWastewater       CommodityTypeType = "wasteWater"
	CommodityTypeTypeDomestichotwater CommodityTypeType = "domesticHotWater"
	CommodityTypeTypeHeatingwater     CommodityTypeType = "heatingWater"
	CommodityTypeTypeSteam            CommodityTypeType = "steam"
	CommodityTypeTypeHeat             CommodityTypeType = "heat"
	CommodityTypeTypeCoolingload      CommodityTypeType = "coolingLoad"
	CommodityTypeTypeAir              CommodityTypeType = "air"
)

type EnergyDirectionType string

const (
	EnergyDirectionTypeConsume EnergyDirectionType = "consume"
	EnergyDirectionTypeProduce EnergyDirectionType = "produce"
)

type EnergyModeType string

const (
	EnergyModeTypeConsume EnergyModeType = "consume"
	EnergyModeTypeProduce EnergyModeType = "produce"
	EnergyModeTypeIdle    EnergyModeType = "idle"
	EnergyModeTypeAuto    EnergyModeType = "auto"
)

type UnitOfMeasurementType string

const (
	UnitOfMeasurementTypeUnknown   UnitOfMeasurementType = "unknown"
	UnitOfMeasurementType1         UnitOfMeasurementType = "1"
	UnitOfMeasurementTypem         UnitOfMeasurementType = "m"
	UnitOfMeasurementTypekg        UnitOfMeasurementType = "kg"
	UnitOfMeasurementTypes         UnitOfMeasurementType = "s"
	UnitOfMeasurementTypeA         UnitOfMeasurementType = "A"
	UnitOfMeasurementTypeK         UnitOfMeasurementType = "K"
	UnitOfMeasurementTypemol       UnitOfMeasurementType = "mol"
	UnitOfMeasurementTypecd        UnitOfMeasurementType = "cd"
	UnitOfMeasurementTypeV         UnitOfMeasurementType = "V"
	UnitOfMeasurementTypeW         UnitOfMeasurementType = "W"
	UnitOfMeasurementTypeWh        UnitOfMeasurementType = "Wh"
	UnitOfMeasurementTypeVA        UnitOfMeasurementType = "VA"
	UnitOfMeasurementTypeVAh       UnitOfMeasurementType = "VAh"
	UnitOfMeasurementTypevar       UnitOfMeasurementType = "var"
	UnitOfMeasurementTypevarh      UnitOfMeasurementType = "varh"
	UnitOfMeasurementTypedegC      UnitOfMeasurementType = "degC"
	UnitOfMeasurementTypedegF      UnitOfMeasurementType = "degF"
	UnitOfMeasurementTypeLm        UnitOfMeasurementType = "Lm"
	UnitOfMeasurementTypelx        UnitOfMeasurementType = "lx"
	UnitOfMeasurementTypeOhm       UnitOfMeasurementType = "Ohm"
	UnitOfMeasurementTypeHz        UnitOfMeasurementType = "Hz"
	UnitOfMeasurementTypedB        UnitOfMeasurementType = "dB"
	UnitOfMeasurementTypedBm       UnitOfMeasurementType = "dBm"
	UnitOfMeasurementTypepct       UnitOfMeasurementType = "pct"
	UnitOfMeasurementTypeppm       UnitOfMeasurementType = "ppm"
	UnitOfMeasurementTypel         UnitOfMeasurementType = "l"
	UnitOfMeasurementTypels        UnitOfMeasurementType = "l/s"
	UnitOfMeasurementTypelh        UnitOfMeasurementType = "l/h"
	UnitOfMeasurementTypedeg       UnitOfMeasurementType = "deg"
	UnitOfMeasurementTyperad       UnitOfMeasurementType = "rad"
	UnitOfMeasurementTyperads      UnitOfMeasurementType = "rad/s"
	UnitOfMeasurementTypesr        UnitOfMeasurementType = "sr"
	UnitOfMeasurementTypeGy        UnitOfMeasurementType = "Gy"
	UnitOfMeasurementTypeBq        UnitOfMeasurementType = "Bq"
	UnitOfMeasurementTypeBqm3      UnitOfMeasurementType = "Bq/m^3"
	UnitOfMeasurementTypeSv        UnitOfMeasurementType = "Sv"
	UnitOfMeasurementTypeRd        UnitOfMeasurementType = "Rd"
	UnitOfMeasurementTypeC         UnitOfMeasurementType = "C"
	UnitOfMeasurementTypeF         UnitOfMeasurementType = "F"
	UnitOfMeasurementTypeH         UnitOfMeasurementType = "H"
	UnitOfMeasurementTypeJ         UnitOfMeasurementType = "J"
	UnitOfMeasurementTypeN         UnitOfMeasurementType = "N"
	UnitOfMeasurementTypeNm        UnitOfMeasurementType = "N_m"
	UnitOfMeasurementTypeNs        UnitOfMeasurementType = "N_s"
	UnitOfMeasurementTypeWb        UnitOfMeasurementType = "Wb"
	UnitOfMeasurementTypeT         UnitOfMeasurementType = "T"
	UnitOfMeasurementTypePa        UnitOfMeasurementType = "Pa"
	UnitOfMeasurementTypebar       UnitOfMeasurementType = "bar"
	UnitOfMeasurementTypeatm       UnitOfMeasurementType = "atm"
	UnitOfMeasurementTypepsi       UnitOfMeasurementType = "psi"
	UnitOfMeasurementTypemmHg      UnitOfMeasurementType = "mmHg"
	UnitOfMeasurementTypem2        UnitOfMeasurementType = "m^2"
	UnitOfMeasurementTypem3        UnitOfMeasurementType = "m^3"
	UnitOfMeasurementTypem3h       UnitOfMeasurementType = "m^3/h"
	UnitOfMeasurementTypems        UnitOfMeasurementType = "m/s"
	UnitOfMeasurementTypems2       UnitOfMeasurementType = "m/s^2"
	UnitOfMeasurementTypem3s       UnitOfMeasurementType = "m^3/s"
	UnitOfMeasurementTypemm3       UnitOfMeasurementType = "m/m^3"
	UnitOfMeasurementTypekgm3      UnitOfMeasurementType = "kg/m^3"
	UnitOfMeasurementTypekgm       UnitOfMeasurementType = "kg_m"
	UnitOfMeasurementTypem2s       UnitOfMeasurementType = "m^2/s"
	UnitOfMeasurementTypewmk       UnitOfMeasurementType = "W/m_K"
	UnitOfMeasurementTypeJK        UnitOfMeasurementType = "J/K"
	UnitOfMeasurementType1s        UnitOfMeasurementType = "1/s"
	UnitOfMeasurementTypeWm2       UnitOfMeasurementType = "W/m^2"
	UnitOfMeasurementTypeJm2       UnitOfMeasurementType = "J/m^2"
	UnitOfMeasurementTypeS         UnitOfMeasurementType = "S"
	UnitOfMeasurementTypeSm        UnitOfMeasurementType = "S/m"
	UnitOfMeasurementTypeKs        UnitOfMeasurementType = "K/s"
	UnitOfMeasurementTypePas       UnitOfMeasurementType = "Pa/s"
	UnitOfMeasurementTypeJkgK      UnitOfMeasurementType = "J/kg_K"
	UnitOfMeasurementTypeVs        UnitOfMeasurementType = "Vs"
	UnitOfMeasurementTypeVm        UnitOfMeasurementType = "V/m"
	UnitOfMeasurementTypeVHz       UnitOfMeasurementType = "V/Hz"
	UnitOfMeasurementTypeAs        UnitOfMeasurementType = "As"
	UnitOfMeasurementTypeAm        UnitOfMeasurementType = "A/m"
	UnitOfMeasurementTypeHzs       UnitOfMeasurementType = "Hz/s"
	UnitOfMeasurementTypekgs       UnitOfMeasurementType = "kg/s"
	UnitOfMeasurementTypekgm2      UnitOfMeasurementType = "kg_m^2"
	UnitOfMeasurementTypeJWh       UnitOfMeasurementType = "J/Wh"
	UnitOfMeasurementTypeWs        UnitOfMeasurementType = "W/s"
	UnitOfMeasurementTypeft3       UnitOfMeasurementType = "ft^3"
	UnitOfMeasurementTypeft3h      UnitOfMeasurementType = "ft^3/h"
	UnitOfMeasurementTypeccf       UnitOfMeasurementType = "ccf"
	UnitOfMeasurementTypeccfh      UnitOfMeasurementType = "ccf/h"
	UnitOfMeasurementTypeUSliqgal  UnitOfMeasurementType = "US.liq.gal"
	UnitOfMeasurementTypeUSliqgalh UnitOfMeasurementType = "US.liq.gal/h"
	UnitOfMeasurementTypeImpgal    UnitOfMeasurementType = "Imp.gal"
	UnitOfMeasurementTypeImpgalh   UnitOfMeasurementType = "Imp.gal/h"
	UnitOfMeasurementTypeBtu       UnitOfMeasurementType = "Btu"
	UnitOfMeasurementTypeBtuh      UnitOfMeasurementType = "Btu/h"
	UnitOfMeasurementTypeAh        UnitOfMeasurementType = "Ah"
	UnitOfMeasurementTypekgWh      UnitOfMeasurementType = "kg/Wh"
)

type CurrencyType string

const (
	CurrencyTypeAed CurrencyType = "AED"
	CurrencyTypeAfn CurrencyType = "AFN"
	CurrencyTypeAll CurrencyType = "ALL"
	CurrencyTypeAmd CurrencyType = "AMD"
	CurrencyTypeAng CurrencyType = "ANG"
	CurrencyTypeAoa CurrencyType = "AOA"
	CurrencyTypeArs CurrencyType = "ARS"
	CurrencyTypeAud CurrencyType = "AUD"
	CurrencyTypeAwg CurrencyType = "AWG"
	CurrencyTypeAzn CurrencyType = "AZN"
	CurrencyTypeBam CurrencyType = "BAM"
	CurrencyTypeBbd CurrencyType = "BBD"
	CurrencyTypeBdt CurrencyType = "BDT"
	CurrencyTypeBgn CurrencyType = "BGN"
	CurrencyTypeBhd CurrencyType = "BHD"
	CurrencyTypeBif CurrencyType = "BIF"
	CurrencyTypeBmd CurrencyType = "BMD"
	CurrencyTypeBnd CurrencyType = "BND"
	CurrencyTypeBob CurrencyType = "BOB"
	CurrencyTypeBov CurrencyType = "BOV"
	CurrencyTypeBrl CurrencyType = "BRL"
	CurrencyTypeBsd CurrencyType = "BSD"
	CurrencyTypeBtn CurrencyType = "BTN"
	CurrencyTypeBwp CurrencyType = "BWP"
	CurrencyTypeByr CurrencyType = "BYR"
	CurrencyTypeBzd CurrencyType = "BZD"
	CurrencyTypeCad CurrencyType = "CAD"
	CurrencyTypeCdf CurrencyType = "CDF"
	CurrencyTypeChe CurrencyType = "CHE"
	CurrencyTypeChf CurrencyType = "CHF"
	CurrencyTypeChw CurrencyType = "CHW"
	CurrencyTypeClf CurrencyType = "CLF"
	CurrencyTypeClp CurrencyType = "CLP"
	CurrencyTypeCny CurrencyType = "CNY"
	CurrencyTypeCop CurrencyType = "COP"
	CurrencyTypeCou CurrencyType = "COU"
	CurrencyTypeCrc CurrencyType = "CRC"
	CurrencyTypeCuc CurrencyType = "CUC"
	CurrencyTypeCup CurrencyType = "CUP"
	CurrencyTypeCve CurrencyType = "CVE"
	CurrencyTypeCzk CurrencyType = "CZK"
	CurrencyTypeDjf CurrencyType = "DJF"
	CurrencyTypeDkk CurrencyType = "DKK"
	CurrencyTypeDop CurrencyType = "DOP"
	CurrencyTypeDzd CurrencyType = "DZD"
	CurrencyTypeEgp CurrencyType = "EGP"
	CurrencyTypeErn CurrencyType = "ERN"
	CurrencyTypeEtb CurrencyType = "ETB"
	CurrencyTypeEur CurrencyType = "EUR"
	CurrencyTypeFjd CurrencyType = "FJD"
	CurrencyTypeFkp CurrencyType = "FKP"
	CurrencyTypeGbp CurrencyType = "GBP"
	CurrencyTypeGel CurrencyType = "GEL"
	CurrencyTypeGhs CurrencyType = "GHS"
	CurrencyTypeGip CurrencyType = "GIP"
	CurrencyTypeGmd CurrencyType = "GMD"
	CurrencyTypeGnf CurrencyType = "GNF"
	CurrencyTypeGtq CurrencyType = "GTQ"
	CurrencyTypeGyd CurrencyType = "GYD"
	CurrencyTypeHkd CurrencyType = "HKD"
	CurrencyTypeHnl CurrencyType = "HNL"
	CurrencyTypeHrk CurrencyType = "HRK"
	CurrencyTypeHtg CurrencyType = "HTG"
	CurrencyTypeHuf CurrencyType = "HUF"
	CurrencyTypeIdr CurrencyType = "IDR"
	CurrencyTypeIls CurrencyType = "ILS"
	CurrencyTypeInr CurrencyType = "INR"
	CurrencyTypeIqd CurrencyType = "IQD"
	CurrencyTypeIrr CurrencyType = "IRR"
	CurrencyTypeIsk CurrencyType = "ISK"
	CurrencyTypeJmd CurrencyType = "JMD"
	CurrencyTypeJod CurrencyType = "JOD"
	CurrencyTypeJpy CurrencyType = "JPY"
	CurrencyTypeKes CurrencyType = "KES"
	CurrencyTypeKgs CurrencyType = "KGS"
	CurrencyTypeKhr CurrencyType = "KHR"
	CurrencyTypeKmf CurrencyType = "KMF"
	CurrencyTypeKpw CurrencyType = "KPW"
	CurrencyTypeKrw CurrencyType = "KRW"
	CurrencyTypeKwd CurrencyType = "KWD"
	CurrencyTypeKyd CurrencyType = "KYD"
	CurrencyTypeKzt CurrencyType = "KZT"
	CurrencyTypeLak CurrencyType = "LAK"
	CurrencyTypeLbp CurrencyType = "LBP"
	CurrencyTypeLkr CurrencyType = "LKR"
	CurrencyTypeLrd CurrencyType = "LRD"
	CurrencyTypeLsl CurrencyType = "LSL"
	CurrencyTypeLyd CurrencyType = "LYD"
	CurrencyTypeMad CurrencyType = "MAD"
	CurrencyTypeMdl CurrencyType = "MDL"
	CurrencyTypeMga CurrencyType = "MGA"
	CurrencyTypeMkd CurrencyType = "MKD"
	CurrencyTypeMmk CurrencyType = "MMK"
	CurrencyTypeMnt CurrencyType = "MNT"
	CurrencyTypeMop CurrencyType = "MOP"
	CurrencyTypeMro CurrencyType = "MRO"
	CurrencyTypeMur CurrencyType = "MUR"
	CurrencyTypeMvr CurrencyType = "MVR"
	CurrencyTypeMwk CurrencyType = "MWK"
	CurrencyTypeMxn CurrencyType = "MXN"
	CurrencyTypeMxv CurrencyType = "MXV"
	CurrencyTypeMyr CurrencyType = "MYR"
	CurrencyTypeMzn CurrencyType = "MZN"
	CurrencyTypeNad CurrencyType = "NAD"
	CurrencyTypeNgn CurrencyType = "NGN"
	CurrencyTypeNio CurrencyType = "NIO"
	CurrencyTypeNok CurrencyType = "NOK"
	CurrencyTypeNpr CurrencyType = "NPR"
	CurrencyTypeNzd CurrencyType = "NZD"
	CurrencyTypeOmr CurrencyType = "OMR"
	CurrencyTypePab CurrencyType = "PAB"
	CurrencyTypePen CurrencyType = "PEN"
	CurrencyTypePgk CurrencyType = "PGK"
	CurrencyTypePhp CurrencyType = "PHP"
	CurrencyTypePkr CurrencyType = "PKR"
	CurrencyTypePln CurrencyType = "PLN"
	CurrencyTypePyg CurrencyType = "PYG"
	CurrencyTypeQar CurrencyType = "QAR"
	CurrencyTypeRon CurrencyType = "RON"
	CurrencyTypeRsd CurrencyType = "RSD"
	CurrencyTypeRub CurrencyType = "RUB"
	CurrencyTypeRwf CurrencyType = "RWF"
	CurrencyTypeSar CurrencyType = "SAR"
	CurrencyTypeSbd CurrencyType = "SBD"
	CurrencyTypeScr CurrencyType = "SCR"
	CurrencyTypeSdg CurrencyType = "SDG"
	CurrencyTypeSek CurrencyType = "SEK"
	CurrencyTypeSgd CurrencyType = "SGD"
	CurrencyTypeShp CurrencyType = "SHP"
	CurrencyTypeSll CurrencyType = "SLL"
	CurrencyTypeSos CurrencyType = "SOS"
	CurrencyTypeSrd CurrencyType = "SRD"
	CurrencyTypeSsp CurrencyType = "SSP"
	CurrencyTypeStd CurrencyType = "STD"
	CurrencyTypeSvc CurrencyType = "SVC"
	CurrencyTypeSyp CurrencyType = "SYP"
	CurrencyTypeSzl CurrencyType = "SZL"
	CurrencyTypeThb CurrencyType = "THB"
	CurrencyTypeTjs CurrencyType = "TJS"
	CurrencyTypeTmt CurrencyType = "TMT"
	CurrencyTypeTnd CurrencyType = "TND"
	CurrencyTypeTop CurrencyType = "TOP"
	CurrencyTypeTry CurrencyType = "TRY"
	CurrencyTypeTtd CurrencyType = "TTD"
	CurrencyTypeTwd CurrencyType = "TWD"
	CurrencyTypeTzs CurrencyType = "TZS"
	CurrencyTypeUah CurrencyType = "UAH"
	CurrencyTypeUgx CurrencyType = "UGX"
	CurrencyTypeUsd CurrencyType = "USD"
	CurrencyTypeUsn CurrencyType = "USN"
	CurrencyTypeUyi CurrencyType = "UYI"
	CurrencyTypeUyu CurrencyType = "UYU"
	CurrencyTypeUzs CurrencyType = "UZS"
	CurrencyTypeVef CurrencyType = "VEF"
	CurrencyTypeVnd CurrencyType = "VND"
	CurrencyTypeVuv CurrencyType = "VUV"
	CurrencyTypeWst CurrencyType = "WST"
	CurrencyTypeXaf CurrencyType = "XAF"
	CurrencyTypeXag CurrencyType = "XAG"
	CurrencyTypeXau CurrencyType = "XAU"
	CurrencyTypeXba CurrencyType = "XBA"
	CurrencyTypeXbb CurrencyType = "XBB"
	CurrencyTypeXbc CurrencyType = "XBC"
	CurrencyTypeXbd CurrencyType = "XBD"
	CurrencyTypeXcd CurrencyType = "XCD"
	CurrencyTypeXdr CurrencyType = "XDR"
	CurrencyTypeXof CurrencyType = "XOF"
	CurrencyTypeXpd CurrencyType = "XPD"
	CurrencyTypeXpf CurrencyType = "XPF"
	CurrencyTypeXpt CurrencyType = "XPT"
	CurrencyTypeXsu CurrencyType = "XSU"
	CurrencyTypeXts CurrencyType = "XTS"
	CurrencyTypeXua CurrencyType = "XUA"
	CurrencyTypeXxx CurrencyType = "XXX"
	CurrencyTypeYer CurrencyType = "YER"
	CurrencyTypeZar CurrencyType = "ZAR"
	CurrencyTypeZmw CurrencyType = "ZMW"
	CurrencyTypeZwl CurrencyType = "ZWL"
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

type ScopeTypeType string

const (
	ScopeTypeTypeAC                    ScopeTypeType = "ac"
	ScopeTypeTypeACCosPhiGrid          ScopeTypeType = "acCosPhiGrid"
	ScopeTypeTypeACCurrentA            ScopeTypeType = "acCurrentA"
	ScopeTypeTypeACCurrentB            ScopeTypeType = "acCurrentB"
	ScopeTypeTypeACCurrentC            ScopeTypeType = "acCurrentC"
	ScopeTypeTypeACFrequencyGrid       ScopeTypeType = "acFrequencyGrid"
	ScopeTypeTypeACPowerA              ScopeTypeType = "acPowerA"
	ScopeTypeTypeACPowerB              ScopeTypeType = "acPowerB"
	ScopeTypeTypeACPowerC              ScopeTypeType = "acPowerC"
	ScopeTypeTypeACPowerLimitPct       ScopeTypeType = "acPowerLimitPct"
	ScopeTypeTypeACPowerTotal          ScopeTypeType = "acPowerTotal"
	ScopeTypeTypeACVoltageA            ScopeTypeType = "acVoltageA"
	ScopeTypeTypeACVoltageB            ScopeTypeType = "acVoltageB"
	ScopeTypeTypeACVoltageC            ScopeTypeType = "acVoltageC"
	ScopeTypeTypeACYieldDay            ScopeTypeType = "acYieldDay"
	ScopeTypeTypeACYieldTotal          ScopeTypeType = "acYieldTotal"
	ScopeTypeTypeDCCurrent             ScopeTypeType = "dcCurrent"
	ScopeTypeTypeDCPower               ScopeTypeType = "dcPower"
	ScopeTypeTypeDCString1             ScopeTypeType = "dcString1"
	ScopeTypeTypeDCString2             ScopeTypeType = "dcString2"
	ScopeTypeTypeDCString3             ScopeTypeType = "dcString3"
	ScopeTypeTypeDCString4             ScopeTypeType = "dcString4"
	ScopeTypeTypeDCString5             ScopeTypeType = "dcString5"
	ScopeTypeTypeDCString6             ScopeTypeType = "dcString6"
	ScopeTypeTypeDCTotal               ScopeTypeType = "dcTotal"
	ScopeTypeTypeDCVoltage             ScopeTypeType = "dcVoltage"
	ScopeTypeTypeDhwTemperature        ScopeTypeType = "dhwTemperature"
	ScopeTypeTypeFlowTemperature       ScopeTypeType = "flowTemperature"
	ScopeTypeTypeOutsideAirTemperature ScopeTypeType = "outsideAirTemperature"
	ScopeTypeTypeReturnTemperature     ScopeTypeType = "returnTemperature"
	ScopeTypeTypeRoomAirTemperature    ScopeTypeType = "roomAirTemperature"
	ScopeTypeTypeCharge                ScopeTypeType = "charge"
	ScopeTypeTypeStateOfCharge         ScopeTypeType = "stateOfCharge"
	ScopeTypeTypeDischarge             ScopeTypeType = "discharge"
	ScopeTypeTypeGridConsumption       ScopeTypeType = "gridConsumption"
	ScopeTypeTypeGridFeedIn            ScopeTypeType = "gridFeedIn"
	ScopeTypeTypeSelfConsumption       ScopeTypeType = "selfConsumption"
	ScopeTypeTypeOverloadProtection    ScopeTypeType = "overloadProtection"
	ScopeTypeTypeACPower               ScopeTypeType = "acPower"
	ScopeTypeTypeACEnergy              ScopeTypeType = "acEnergy"
	ScopeTypeTypeACCurrent             ScopeTypeType = "acCurrent"
	ScopeTypeTypeACVoltage             ScopeTypeType = "acVoltage"
	ScopeTypeTypeBatteryControl        ScopeTypeType = "batteryControl"
	ScopeTypeTypeSimpleIncentiveTable  ScopeTypeType = "simpleIncentiveTable"
)

type RoleType string

const (
	RoleTypeClient  RoleType = "client"
	RoleTypeServer  RoleType = "server"
	RoleTypeSpecial RoleType = "special"
)

type FeatureGroupType string

type DeviceTypeType string

const (
	DeviceTypeTypeDishwasher              DeviceTypeType = "Dishwasher"
	DeviceTypeTypeDryer                   DeviceTypeType = "Dryer"
	DeviceTypeTypeEnvironmentSensor       DeviceTypeType = "EnvironmentSensor"
	DeviceTypeTypeGeneric                 DeviceTypeType = "Generic"
	DeviceTypeTypeHeatgenerationSystem    DeviceTypeType = "HeatGenerationSystem"
	DeviceTypeTypeHeatsinkSystem          DeviceTypeType = "HeatSinkSystem"
	DeviceTypeTypeHeatstorageSystem       DeviceTypeType = "HeatStorageSystem"
	DeviceTypeTypeHVACController          DeviceTypeType = "HVACController"
	DeviceTypeTypeSubmeter                DeviceTypeType = "SubMeter"
	DeviceTypeTypeWasher                  DeviceTypeType = "Washer"
	DeviceTypeTypeElectricitySupplySystem DeviceTypeType = "ElectricitySupplySystem"
	DeviceTypeTypeEnergyManagementSystem  DeviceTypeType = "EnergyManagementSystem"
	DeviceTypeTypeInverter                DeviceTypeType = "Inverter"
	DeviceTypeTypeChargingStation         DeviceTypeType = "ChargingStation"
)

type EntityTypeType string

const (
	EntityTypeTypeBattery                       EntityTypeType = "Battery"
	EntityTypeTypeCompressor                    EntityTypeType = "Compressor"
	EntityTypeTypeDeviceInformation             EntityTypeType = "DeviceInformation"
	EntityTypeTypeDHWCircuit                    EntityTypeType = "DHWCircuit"
	EntityTypeTypeDHWStorage                    EntityTypeType = "DHWStorage"
	EntityTypeTypeDishwasher                    EntityTypeType = "Dishwasher"
	EntityTypeTypeDryer                         EntityTypeType = "Dryer"
	EntityTypeTypeElectricalImmersionheater     EntityTypeType = "ElectricalImmersionHeater"
	EntityTypeTypeFan                           EntityTypeType = "Fan"
	EntityTypeTypeGasHeatingAppliance           EntityTypeType = "GasHeatingAppliance"
	EntityTypeTypeGeneric                       EntityTypeType = "Generic"
	EntityTypeTypeHeatingBufferStorage          EntityTypeType = "HeatingBufferStorage"
	EntityTypeTypeHeatingCircuit                EntityTypeType = "HeatingCircuit"
	EntityTypeTypeHeatingObject                 EntityTypeType = "HeatingObject"
	EntityTypeTypeHeatingZone                   EntityTypeType = "HeatingZone"
	EntityTypeTypeHeatPumpAppliance             EntityTypeType = "HeatPumpAppliance"
	EntityTypeTypeHeatSinkCircuit               EntityTypeType = "HeatSinkCircuit"
	EntityTypeTypeHeatSourceCircuit             EntityTypeType = "HeatSourceCircuit"
	EntityTypeTypeHeatSourceUnit                EntityTypeType = "HeatSourceUnit"
	EntityTypeTypeHvacController                EntityTypeType = "HVACController"
	EntityTypeTypeHvacRoom                      EntityTypeType = "HVACRoom"
	EntityTypeTypeInstantDHWheater              EntityTypeType = "InstantDHWHeater"
	EntityTypeTypeInverter                      EntityTypeType = "Inverter"
	EntityTypeTypeOilHeatingAppliance           EntityTypeType = "OilHeatingAppliance"
	EntityTypeTypePump                          EntityTypeType = "Pump"
	EntityTypeTypeRefrigerantCircuit            EntityTypeType = "RefrigerantCircuit"
	EntityTypeTypeSmartEnergyAppliance          EntityTypeType = "SmartEnergyAppliance"
	EntityTypeTypeSolarDHWStorage               EntityTypeType = "SolarDHWStorage"
	EntityTypeTypeSolarThermalCircuit           EntityTypeType = "SolarThermalCircuit"
	EntityTypeTypeSubmeterElectricity           EntityTypeType = "SubMeterElectricity"
	EntityTypeTypeTemperatureSensor             EntityTypeType = "TemperatureSensor"
	EntityTypeTypeWasher                        EntityTypeType = "Washer"
	EntityTypeTypeBatterySystem                 EntityTypeType = "BatterySystem"
	EntityTypeTypeElectricityGenerationSystem   EntityTypeType = "ElectricityGenerationSystem"
	EntityTypeTypeElectricityStorageSystem      EntityTypeType = "ElectricityStorageSystem"
	EntityTypeTypeGridConnectionPointOfPremises EntityTypeType = "GridConnectionPointOfPremises"
	EntityTypeTypeHousehold                     EntityTypeType = "Household"
	EntityTypeTypePVSystem                      EntityTypeType = "PVSystem"
	EntityTypeTypeEV                            EntityTypeType = "EV"
	EntityTypeTypeEVSE                          EntityTypeType = "EVSE"
	EntityTypeTypeChargingOutlet                EntityTypeType = "ChargingOutlet"
	EntityTypeTypeCEM                           EntityTypeType = "CEM"
)

type FeatureTypeType string

const (
	FeatureTypeTypeActuatorLevel           FeatureTypeType = "ActuatorLevel"
	FeatureTypeTypeActuatorSwitch          FeatureTypeType = "ActuatorSwitch"
	FeatureTypeTypeAlarm                   FeatureTypeType = "Alarm"
	FeatureTypeTypeDataTunneling           FeatureTypeType = "DataTunneling"
	FeatureTypeTypeDeviceClassification    FeatureTypeType = "DeviceClassification"
	FeatureTypeTypeDeviceDiagnosis         FeatureTypeType = "DeviceDiagnosis"
	FeatureTypeTypeDirectControl           FeatureTypeType = "DirectControl"
	FeatureTypeTypeElectricalConnection    FeatureTypeType = "ElectricalConnection"
	FeatureTypeTypeGeneric                 FeatureTypeType = "Generic"
	FeatureTypeTypeHvac                    FeatureTypeType = "HVAC"
	FeatureTypeTypeLoadControl             FeatureTypeType = "LoadControl"
	FeatureTypeTypeMeasurement             FeatureTypeType = "Measurement"
	FeatureTypeTypeMessaging               FeatureTypeType = "Messaging"
	FeatureTypeTypeNetworkManagement       FeatureTypeType = "NetworkManagement"
	FeatureTypeTypeNodeManagement          FeatureTypeType = "NodeManagement"
	FeatureTypeTypeOperatingConstraints    FeatureTypeType = "OperatingConstraints"
	FeatureTypeTypePowerSequences          FeatureTypeType = "PowerSequences"
	FeatureTypeTypeSensing                 FeatureTypeType = "Sensing"
	FeatureTypeTypeSetpoint                FeatureTypeType = "Setpoint"
	FeatureTypeTypeSmartEnergyManagementPs FeatureTypeType = "SmartEnergyManagementPs"
	FeatureTypeTypeTaskManagement          FeatureTypeType = "TaskManagement"
	FeatureTypeTypeThreshold               FeatureTypeType = "Threshold"
	FeatureTypeTypeTimeInformation         FeatureTypeType = "TimeInformation"
	FeatureTypeTypeTimeTable               FeatureTypeType = "TimeTable"
	FeatureTypeTypeDeviceConfiguration     FeatureTypeType = "DeviceConfiguration"
	FeatureTypeTypeSupplyCondition         FeatureTypeType = "SupplyCondition"
	FeatureTypeTypeTimeSeries              FeatureTypeType = "TimeSeries"
	FeatureTypeTypeTariffInformation       FeatureTypeType = "TariffInformation"
	FeatureTypeTypeIncentiveTable          FeatureTypeType = "IncentiveTable"
	FeatureTypeTypeBill                    FeatureTypeType = "Bill"
	FeatureTypeTypeIdentification          FeatureTypeType = "Identification"
)

type FeatureSpecificUsageType string

const (
	// FeatureDirectControlSpecificUsageEnumType
	FeatureSpecificUsageTypeHistory  FeatureSpecificUsageType = "History"
	FeatureSpecificUsageTypeRealtime FeatureSpecificUsageType = "RealTime"

	// FeatureHvacSpecificUsageEnumType
	FeatureSpecificUsageTypeOperationmode FeatureSpecificUsageType = "OperationMode"
	FeatureSpecificUsageTypeOverrun       FeatureSpecificUsageType = "Overrun"

	// FeatureMeasurementSpecificUsageEnumType
	FeatureSpecificUsageTypeContact     FeatureSpecificUsageType = "Contact"
	FeatureSpecificUsageTypeElectrical  FeatureSpecificUsageType = "Electrical"
	FeatureSpecificUsageTypeHeat        FeatureSpecificUsageType = "Heat"
	FeatureSpecificUsageTypeLevel       FeatureSpecificUsageType = "Level"
	FeatureSpecificUsageTypePressure    FeatureSpecificUsageType = "Pressure"
	FeatureSpecificUsageTypeTemperature FeatureSpecificUsageType = "Temperature"

	// FeatureSetpointSpecificUsageEnumType

	// FeatureSmartEnergyManagementPsSpecificUsageEnumType
	FeatureSpecificUsageTypeFixedForecast                         FeatureSpecificUsageType = "FixedForecast"
	FeatureSpecificUsageTypeFlexibleChosenForecast                FeatureSpecificUsageType = "FlexibleChosenForecast"
	FeatureSpecificUsageTypeFlexibleOptionalForecast              FeatureSpecificUsageType = "FlexibleOptionalForecast"
	FeatureSpecificUsageTypeOptionalSequenceBasedImmediateControl FeatureSpecificUsageType = "OptionalSequenceBasedImmediateControl"
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

const (
	// FeatureMeasurementSpecificUsageEnumType
	FeatureSetpointSpecificUsageEnumTypeContact     FeatureSetpointSpecificUsageEnumType = "Contact"
	FeatureSetpointSpecificUsageEnumTypeElectrical  FeatureSetpointSpecificUsageEnumType = "Electrical"
	FeatureSetpointSpecificUsageEnumTypeHeat        FeatureSetpointSpecificUsageEnumType = "Heat"
	FeatureSetpointSpecificUsageEnumTypeLevel       FeatureSetpointSpecificUsageEnumType = "Level"
	FeatureSetpointSpecificUsageEnumTypePressure    FeatureSetpointSpecificUsageEnumType = "Pressure"
	FeatureSetpointSpecificUsageEnumTypeTemperature FeatureSetpointSpecificUsageEnumType = "Temperature"
)

type FeatureSmartEnergyManagementPsSpecificUsageEnumType string

const (
	FeatureSmartEnergyManagementPsSpecificUsageEnumTypeFixedForecast                         FeatureSmartEnergyManagementPsSpecificUsageEnumType = "FixedForecast"
	FeatureSmartEnergyManagementPsSpecificUsageEnumTypeFlexibleChosenForecast                FeatureSmartEnergyManagementPsSpecificUsageEnumType = "FlexibleChosenForecast"
	FeatureSmartEnergyManagementPsSpecificUsageEnumTypeFlexibleOptionalForecast              FeatureSmartEnergyManagementPsSpecificUsageEnumType = "FlexibleOptionalForecast"
	FeatureSmartEnergyManagementPsSpecificUsageEnumTypeOptionalSequenceBasedImmediateControl FeatureSmartEnergyManagementPsSpecificUsageEnumType = "OptionalSequenceBasedImmediateControl"
)

type FunctionType string

const (
	FunctionTypeActuatorLevelData                                  FunctionType = "actuatorLevelData"
	FunctionTypeActuatorLevelDescriptionData                       FunctionType = "actuatorLevelDescriptionData"
	FunctionTypeActuatorSwitchData                                 FunctionType = "actuatorSwitchData"
	FunctionTypeActuatorSwitchDescriptionData                      FunctionType = "actuatorSwitchDescriptionData"
	FunctionTypeAlarmListData                                      FunctionType = "alarmListData"
	FunctionTypeBindingManagementDeleteCall                        FunctionType = "bindingManagementDeleteCall"
	FunctionTypeBindingManagementEntryListData                     FunctionType = "bindingManagementEntryListData"
	FunctionTypeBindingManagementRequestCall                       FunctionType = "bindingManagementRequestCall"
	FunctionTypeDataTunnelingCall                                  FunctionType = "dataTunnelingCall"
	FunctionTypeDeviceClassificationManufacturerData               FunctionType = "deviceClassificationManufacturerData"
	FunctionTypeDeviceClassificationUserData                       FunctionType = "deviceClassificationUserData"
	FunctionTypeDeviceDiagnosisHeartbeatData                       FunctionType = "deviceDiagnosisHeartbeatData"
	FunctionTypeDeviceDiagnosisServiceData                         FunctionType = "deviceDiagnosisServiceData"
	FunctionTypeDeviceDiagnosisStateData                           FunctionType = "deviceDiagnosisStateData"
	FunctionTypeDirectControlActivityListData                      FunctionType = "directControlActivityListData"
	FunctionTypeDirectControlDescriptionData                       FunctionType = "directControlDescriptionData"
	FunctionTypeElectricalConnectionDescriptionListData            FunctionType = "electricalConnectionDescriptionListData"
	FunctionTypeElectricalConnectionParameterDescriptionListData   FunctionType = "electricalConnectionParameterDescriptionListData"
	FunctionTypeElectricalConnectionStateListData                  FunctionType = "electricalConnectionStateListData"
	FunctionTypeHvacOperationModeDescriptionListData               FunctionType = "hvacOperationModeDescriptionListData"
	FunctionTypeHvacOverrunDescriptionListData                     FunctionType = "hvacOverrunDescriptionListData"
	FunctionTypeHvacOverrunListData                                FunctionType = "hvacOverrunListData"
	FunctionTypeHvacSystemFunctionDescriptionListData              FunctionType = "hvacSystemFunctionDescriptionListData"
	FunctionTypeHvacSystemFunctionListData                         FunctionType = "hvacSystemFunctionListData"
	FunctionTypeHvacSystemFunctionOperationModeRelationListData    FunctionType = "hvacSystemFunctionOperationModeRelationListData"
	FunctionTypeHvacSystemFunctionPowerSequenceRelationListData    FunctionType = "hvacSystemFunctionPowerSequenceRelationListData"
	FunctionTypeHvacSystemFunctionSetPointRelationListData         FunctionType = "hvacSystemFunctionSetpointRelationListData"
	FunctionTypeLoadControlEventListData                           FunctionType = "loadControlEventListData"
	FunctionTypeLoadControlStateListData                           FunctionType = "loadControlStateListData"
	FunctionTypeMeasurementConstraintsListData                     FunctionType = "measurementConstraintsListData"
	FunctionTypeMeasurementDescriptionListData                     FunctionType = "measurementDescriptionListData"
	FunctionTypeMeasurementListData                                FunctionType = "measurementListData"
	FunctionTypeMeasurementThresholdRelationListData               FunctionType = "measurementThresholdRelationListData"
	FunctionTypeMessagingListData                                  FunctionType = "messagingListData"
	FunctionTypeNetworkManagementAbortCall                         FunctionType = "networkManagementAbortCall"
	FunctionTypeNetworkManagementAddnodeCall                       FunctionType = "networkManagementAddNodeCall"
	FunctionTypeNetworkManagementDeviceDescriptionListData         FunctionType = "networkManagementDeviceDescriptionListData"
	FunctionTypeNetworkManagementDiscoverCall                      FunctionType = "networkManagementDiscoverCall"
	FunctionTypeNetworkManagementEntityDescriptionListData         FunctionType = "networkManagementEntityDescriptionListData"
	FunctionTypeNetworkManagementFeatureDescriptionListData        FunctionType = "networkManagementFeatureDescriptionListData"
	FunctionTypeNetworkManagementJoiningModeData                   FunctionType = "networkManagementJoiningModeData"
	FunctionTypeNetworkManagementModifyNodeCall                    FunctionType = "networkManagementModifyNodeCall"
	FunctionTypeNetworkManagementProcessStateData                  FunctionType = "networkManagementProcessStateData"
	FunctionTypeNetworkManagementRemoveNodeCall                    FunctionType = "networkManagementRemoveNodeCall"
	FunctionTypeNetworkManagementReportCandidateData               FunctionType = "networkManagementReportCandidateData"
	FunctionTypeNetworkManagementScanNetworkCall                   FunctionType = "networkManagementScanNetworkCall"
	FunctionTypeNodeManagementBindingData                          FunctionType = "nodeManagementBindingData"
	FunctionTypeNodeManagementBindingDeleteCall                    FunctionType = "nodeManagementBindingDeleteCall"
	FunctionTypeNodeManagementBindingRequestCall                   FunctionType = "nodeManagementBindingRequestCall"
	FunctionTypeNodeManagementDestinationListData                  FunctionType = "nodeManagementDestinationListData"
	FunctionTypeNodeManagementDetailedDiscoveryData                FunctionType = "nodeManagementDetailedDiscoveryData"
	FunctionTypeNodeManagementSubscriptionData                     FunctionType = "nodeManagementSubscriptionData"
	FunctionTypeNodeManagementSubscriptionDeleteCall               FunctionType = "nodeManagementSubscriptionDeleteCall"
	FunctionTypeNodeManagementSubscriptionRequestCall              FunctionType = "nodeManagementSubscriptionRequestCall"
	FunctionTypeOperatingConstraintsDurationListData               FunctionType = "operatingConstraintsDurationListData"
	FunctionTypeOperatingConstraintsInterruptListData              FunctionType = "operatingConstraintsInterruptListData"
	FunctionTypeOperatingConstraintsPowerDescriptionListData       FunctionType = "operatingConstraintsPowerDescriptionListData"
	FunctionTypeOperatingConstraintsPowerLevelListData             FunctionType = "operatingConstraintsPowerLevelListData"
	FunctionTypeOperatingConstraintsPowerRangeListData             FunctionType = "operatingConstraintsPowerRangeListData"
	FunctionTypeOperatingConstraintsResumeImplicationListData      FunctionType = "operatingConstraintsResumeImplicationListData"
	FunctionTypePowerSequenceAlternativesRelationlistdata          FunctionType = "powerSequenceAlternativesRelationListData"
	FunctionTypePowerSequenceDescriptionListData                   FunctionType = "powerSequenceDescriptionListData"
	FunctionTypePowerSequenceNodeScheduleInformationData           FunctionType = "powerSequenceNodeScheduleInformationData"
	FunctionTypePowerSequencePriceCalculationRequestCall           FunctionType = "powerSequencePriceCalculationRequestCall"
	FunctionTypePowerSequencePriceListData                         FunctionType = "powerSequencePriceListData"
	FunctionTypePowerSequenceScheduleConfigurationRequestCall      FunctionType = "powerSequenceScheduleConfigurationRequestCall"
	FunctionTypePowerSequenceScheduleConstraintsListData           FunctionType = "powerSequenceScheduleConstraintsListData"
	FunctionTypePowerSequenceScheduleListData                      FunctionType = "powerSequenceScheduleListData"
	FunctionTypePowerSequenceSchedulePreferenceListData            FunctionType = "powerSequenceSchedulePreferenceListData"
	FunctionTypePowerSequenceStateListData                         FunctionType = "powerSequenceStateListData"
	FunctionTypePowerTimeslotScheduleConstraintsListData           FunctionType = "powerTimeSlotScheduleConstraintsListData"
	FunctionTypePowerTimeslotScheduleListData                      FunctionType = "powerTimeSlotScheduleListData"
	FunctionTypePowerTimeslotValueListData                         FunctionType = "powerTimeSlotValueListData"
	FunctionTypeResultData                                         FunctionType = "resultData"
	FunctionTypeSensingDescriptionData                             FunctionType = "sensingDescriptionData"
	FunctionTypeSensingListData                                    FunctionType = "sensingListData"
	FunctionTypeSetPointConstraintsListData                        FunctionType = "setpointConstraintsListData"
	FunctionTypeSetPointDescriptionListData                        FunctionType = "setpointDescriptionListData"
	FunctionTypeSetPointListData                                   FunctionType = "setpointListData"
	FunctionTypeSmartEnergyManagementPsConfigurationRequestCall    FunctionType = "smartEnergyManagementPsConfigurationRequestCall"
	FunctionTypeSmartEnergyManagementPsData                        FunctionType = "smartEnergyManagementPsData"
	FunctionTypeSmartEnergyManagementPsPriceCalculationRequestCall FunctionType = "smartEnergyManagementPsPriceCalculationRequestCall"
	FunctionTypeSmartEnergyManagementPsPriceData                   FunctionType = "smartEnergyManagementPsPriceData"
	FunctionTypeSpecificationVersionListData                       FunctionType = "specificationVersionListData"
	FunctionTypeSubscriptionManagementDeleteCall                   FunctionType = "subscriptionManagementDeleteCall"
	FunctionTypeSubscriptionManagementEntryListData                FunctionType = "subscriptionManagementEntryListData"
	FunctionTypeSubscriptionManagementRequestCall                  FunctionType = "subscriptionManagementRequestCall"
	FunctionTypeSupplyConditionDescriptionListData                 FunctionType = "supplyConditionDescriptionListData"
	FunctionTypeSupplyConditionListData                            FunctionType = "supplyConditionListData"
	FunctionTypeSupplyConditionThresholdRelationListData           FunctionType = "supplyConditionThresholdRelationListData"
	FunctionTypeTaskManagementJobDescriptionListData               FunctionType = "taskManagementJobDescriptionListData"
	FunctionTypeTaskManagementJobListData                          FunctionType = "taskManagementJobListData"
	FunctionTypeTaskManagementJobRelationListData                  FunctionType = "taskManagementJobRelationListData"
	FunctionTypeTaskManagementOverviewData                         FunctionType = "taskManagementOverviewData"
	FunctionTypeThresholdConstraintsListData                       FunctionType = "thresholdConstraintsListData"
	FunctionTypeThresholdDescriptionListData                       FunctionType = "thresholdDescriptionListData"
	FunctionTypeThresholdListData                                  FunctionType = "thresholdListData"
	FunctionTypeTimeDistributorData                                FunctionType = "timeDistributorData"
	FunctionTypeTimeDistributorEnquiryCall                         FunctionType = "timeDistributorEnquiryCall"
	FunctionTypeTimeInformationData                                FunctionType = "timeInformationData"
	FunctionTypeTimePrecisionData                                  FunctionType = "timePrecisionData"
	FunctionTypeTimeTableConstraintsListData                       FunctionType = "timeTableConstraintsListData"
	FunctionTypeTimeTableDescriptionListData                       FunctionType = "timeTableDescriptionListData"
	FunctionTypeTimeTableListData                                  FunctionType = "timeTableListData"
	FunctionTypeDeviceConfigurationKeyValueConstraintsListData     FunctionType = "deviceConfigurationKeyValueConstraintsListData"
	FunctionTypeDeviceConfigurationKeyValueListData                FunctionType = "deviceConfigurationKeyValueListData"
	FunctionTypeDeviceConfigurationKeyValueDescriptionListData     FunctionType = "deviceConfigurationKeyValueDescriptionListData"
	FunctionTypeLoadControlLimitConstraintsListData                FunctionType = "loadControlLimitConstraintsListData"
	FunctionTypeLoadControlLimitDescriptionListData                FunctionType = "loadControlLimitDescriptionListData"
	FunctionTypeLoadControlLimitListData                           FunctionType = "loadControlLimitListData"
	FunctionTypeLoadControlNodeData                                FunctionType = "loadControlNodeData"
	FunctionTypeTimeSeriesConstraintsListData                      FunctionType = "timeSeriesConstraintsListData"
	FunctionTypeTimeSeriesDescriptionListData                      FunctionType = "timeSeriesDescriptionListData"
	FunctionTypeTimeSeriesListData                                 FunctionType = "timeSeriesListData"
	FunctionTypeTariffOverallConstraintsData                       FunctionType = "tariffOverallConstraintsData"
	FunctionTypeTariffListData                                     FunctionType = "tariffListData"
	FunctionTypeTariffBoundaryrelationListData                     FunctionType = "tariffBoundaryRelationListData"
	FunctionTypeTariffTierRelationListData                         FunctionType = "tariffTierRelationListData"
	FunctionTypeTariffDescriptionListData                          FunctionType = "tariffDescriptionListData"
	FunctionTypeTierBoundaryListData                               FunctionType = "tierBoundaryListData"
	FunctionTypeTierBoundaryDescriptionListData                    FunctionType = "tierBoundaryDescriptionListData"
	FunctionTypeCommodityListData                                  FunctionType = "commodityListData"
	FunctionTypeTierListData                                       FunctionType = "tierListData"
	FunctionTypeTierIncentiveRelationListData                      FunctionType = "tierIncentiveRelationListData"
	FunctionTypeTierDescriptionListData                            FunctionType = "tierDescriptionListData"
	FunctionTypeIncentiveListData                                  FunctionType = "incentiveListData"
	FunctionTypeIncentiveDescriptionListData                       FunctionType = "incentiveDescriptionListData"
	FunctionTypeIncentiveTableData                                 FunctionType = "incentiveTableData"
	FunctionTypeIncentiveTableDescriptionData                      FunctionType = "incentiveTableDescriptionData"
	FunctionTypeIncentiveTableConstraintsData                      FunctionType = "incentiveTableConstraintsData"
	FunctionTypeElectricalConnectionPermittedValueSetListData      FunctionType = "electricalConnectionPermittedValueSetListData"
	FunctionTypeUseCaseInformationListData                         FunctionType = "useCaseInformationListData"
	FunctionTypeNodeManagementUseCaseData                          FunctionType = "nodeManagementUseCaseData"
	FunctionTypeBillConstraintsListData                            FunctionType = "billConstraintsListData"
	FunctionTypeBillDescriptionListData                            FunctionType = "billDescriptionListData"
	FunctionTypeBillListData                                       FunctionType = "billListData"
	FunctionTypeIdentificationListData                             FunctionType = "identificationListData"
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
