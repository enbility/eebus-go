package spine

type SensingStateType SensingStateEnumType

type SensingStateEnumType string

const (
	SensingStateEnumTypeOn                      SensingStateEnumType = "on"
	SensingStateEnumTypeOff                     SensingStateEnumType = "off"
	SensingStateEnumTypeToggle                  SensingStateEnumType = "toggle"
	SensingStateEnumTypeLevel                   SensingStateEnumType = "level"
	SensingStateEnumTypeLevelUp                 SensingStateEnumType = "levelUp"
	SensingStateEnumTypeLevelDown               SensingStateEnumType = "levelDown"
	SensingStateEnumTypeLevelStart              SensingStateEnumType = "levelStart"
	SensingStateEnumTypeLevelStop               SensingStateEnumType = "levelStop"
	SensingStateEnumTypeLevelAbsolute           SensingStateEnumType = "levelAbsolute"
	SensingStateEnumTypeLevelRelative           SensingStateEnumType = "levelRelative"
	SensingStateEnumTypeLevelPercentageAbsolute SensingStateEnumType = "levelPercentageAbsolute"
	SensingStateEnumTypeLevelPercentageRelative SensingStateEnumType = "levelPercentageRelative"
	SensingStateEnumTypePressed                 SensingStateEnumType = "pressed"
	SensingStateEnumTypeLongPressed             SensingStateEnumType = "longPressed"
	SensingStateEnumTypeReleased                SensingStateEnumType = "released"
	SensingStateEnumTypeChanged                 SensingStateEnumType = "changed"
	SensingStateEnumTypeStarted                 SensingStateEnumType = "started"
	SensingStateEnumTypeStopped                 SensingStateEnumType = "stopped"
	SensingStateEnumTypePaused                  SensingStateEnumType = "paused"
	SensingStateEnumTypeMiddle                  SensingStateEnumType = "middle"
	SensingStateEnumTypeUp                      SensingStateEnumType = "up"
	SensingStateEnumTypeDown                    SensingStateEnumType = "down"
	SensingStateEnumTypeForward                 SensingStateEnumType = "forward"
	SensingStateEnumTypeBackwards               SensingStateEnumType = "backwards"
	SensingStateEnumTypeOpen                    SensingStateEnumType = "open"
	SensingStateEnumTypeClosed                  SensingStateEnumType = "closed"
	SensingStateEnumTypeOpening                 SensingStateEnumType = "opening"
	SensingStateEnumTypeClosing                 SensingStateEnumType = "closing"
	SensingStateEnumTypeHigh                    SensingStateEnumType = "high"
	SensingStateEnumTypeLow                     SensingStateEnumType = "low"
	SensingStateEnumTypeDay                     SensingStateEnumType = "day"
	SensingStateEnumTypeNight                   SensingStateEnumType = "night"
	SensingStateEnumTypeDetected                SensingStateEnumType = "detected"
	SensingStateEnumTypeNotDetected             SensingStateEnumType = "notDetected"
	SensingStateEnumTypeAlarmed                 SensingStateEnumType = "alarmed"
	SensingStateEnumTypeNotAlarmed              SensingStateEnumType = "notAlarmed"
)

type SensingTypeType SensingTypeEnumType

type SensingTypeEnumType string

const (
	SensingTypeEnumTypeSwitch            SensingTypeEnumType = "switch"
	SensingTypeEnumTypeButton            SensingTypeEnumType = "button"
	SensingTypeEnumTypeLevel             SensingTypeEnumType = "level"
	SensingTypeEnumTypeLevelSwitch       SensingTypeEnumType = "levelSwitch"
	SensingTypeEnumTypeWindowHandle      SensingTypeEnumType = "windowHandle"
	SensingTypeEnumTypeContactSensor     SensingTypeEnumType = "contactSensor"
	SensingTypeEnumTypeOccupancySensor   SensingTypeEnumType = "occupancySensor"
	SensingTypeEnumTypeMotionDetector    SensingTypeEnumType = "motionDetector"
	SensingTypeEnumTypeFireDetector      SensingTypeEnumType = "fireDetector"
	SensingTypeEnumTypeSmokeDetector     SensingTypeEnumType = "smokeDetector"
	SensingTypeEnumTypeHeatDetector      SensingTypeEnumType = "heatDetector"
	SensingTypeEnumTypeWaterDetector     SensingTypeEnumType = "waterDetector"
	SensingTypeEnumTypeGasDetector       SensingTypeEnumType = "gasDetector"
	SensingTypeEnumTypeAlarmSensor       SensingTypeEnumType = "alarmSensor"
	SensingTypeEnumTypePowerAlarmSensor  SensingTypeEnumType = "powerAlarmSensor"
	SensingTypeEnumTypeDayNightIndicator SensingTypeEnumType = "dayNightIndicator"
)

type SensingDataType struct {
	Timestamp *AbsoluteOrRelativeTimeType `json:"timestamp,omitempty"`
	State     *SensingStateType           `json:"state,omitempty"`
	Value     *ScaledNumberType           `json:"value,omitempty"`
}

type SensingDataElementsType struct {
	Timestamp *ElementTagType           `json:"timestamp,omitempty"`
	State     *ElementTagType           `json:"state,omitempty"`
	Value     *ScaledNumberElementsType `json:"value,omitempty"`
}

type SensingListDataType struct {
	SensingData []SensingDataType `json:"sensingData,omitempty"`
}

type SensingListDataSelectorsType struct {
	TimestampInterval *TimestampIntervalType `json:"timestampInterval,omitempty"`
}

type SensingDescriptionDataType struct {
	SensingType *SensingTypeType       `json:"sensingType,omitempty"`
	Unit        *UnitOfMeasurementType `json:"unit,omitempty"`
	ScopeType   *ScopeTypeType         `json:"scopeType,omitempty"`
	Label       *LabelType             `json:"label,omitempty"`
	Description *DescriptionType       `json:"description,omitempty"`
}

type SensingDescriptionDataElementsType struct {
	SensingType *ElementTagType `json:"sensingType,omitempty"`
	Unit        *ElementTagType `json:"unit,omitempty"`
	ScopeType   *ElementTagType `json:"scopeType,omitempty"`
	Label       *ElementTagType `json:"label,omitempty"`
	Description *ElementTagType `json:"description,omitempty"`
}
