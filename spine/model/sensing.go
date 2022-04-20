package model

type SensingStateType string

const (
	SensingStateTypeOn                      SensingStateType = "on"
	SensingStateTypeOff                     SensingStateType = "off"
	SensingStateTypeToggle                  SensingStateType = "toggle"
	SensingStateTypeLevel                   SensingStateType = "level"
	SensingStateTypeLevelUp                 SensingStateType = "levelUp"
	SensingStateTypeLevelDown               SensingStateType = "levelDown"
	SensingStateTypeLevelStart              SensingStateType = "levelStart"
	SensingStateTypeLevelStop               SensingStateType = "levelStop"
	SensingStateTypeLevelAbsolute           SensingStateType = "levelAbsolute"
	SensingStateTypeLevelRelative           SensingStateType = "levelRelative"
	SensingStateTypeLevelPercentageAbsolute SensingStateType = "levelPercentageAbsolute"
	SensingStateTypeLevelPercentageRelative SensingStateType = "levelPercentageRelative"
	SensingStateTypePressed                 SensingStateType = "pressed"
	SensingStateTypeLongPressed             SensingStateType = "longPressed"
	SensingStateTypeReleased                SensingStateType = "released"
	SensingStateTypeChanged                 SensingStateType = "changed"
	SensingStateTypeStarted                 SensingStateType = "started"
	SensingStateTypeStopped                 SensingStateType = "stopped"
	SensingStateTypePaused                  SensingStateType = "paused"
	SensingStateTypeMiddle                  SensingStateType = "middle"
	SensingStateTypeUp                      SensingStateType = "up"
	SensingStateTypeDown                    SensingStateType = "down"
	SensingStateTypeForward                 SensingStateType = "forward"
	SensingStateTypeBackwards               SensingStateType = "backwards"
	SensingStateTypeOpen                    SensingStateType = "open"
	SensingStateTypeClosed                  SensingStateType = "closed"
	SensingStateTypeOpening                 SensingStateType = "opening"
	SensingStateTypeClosing                 SensingStateType = "closing"
	SensingStateTypeHigh                    SensingStateType = "high"
	SensingStateTypeLow                     SensingStateType = "low"
	SensingStateTypeDay                     SensingStateType = "day"
	SensingStateTypeNight                   SensingStateType = "night"
	SensingStateTypeDetected                SensingStateType = "detected"
	SensingStateTypeNotDetected             SensingStateType = "notDetected"
	SensingStateTypeAlarmed                 SensingStateType = "alarmed"
	SensingStateTypeNotAlarmed              SensingStateType = "notAlarmed"
)

type SensingTypeType string

const (
	SensingTypeTypeSwitch            SensingTypeType = "switch"
	SensingTypeTypeButton            SensingTypeType = "button"
	SensingTypeTypeLevel             SensingTypeType = "level"
	SensingTypeTypeLevelSwitch       SensingTypeType = "levelSwitch"
	SensingTypeTypeWindowHandle      SensingTypeType = "windowHandle"
	SensingTypeTypeContactSensor     SensingTypeType = "contactSensor"
	SensingTypeTypeOccupancySensor   SensingTypeType = "occupancySensor"
	SensingTypeTypeMotionDetector    SensingTypeType = "motionDetector"
	SensingTypeTypeFireDetector      SensingTypeType = "fireDetector"
	SensingTypeTypeSmokeDetector     SensingTypeType = "smokeDetector"
	SensingTypeTypeHeatDetector      SensingTypeType = "heatDetector"
	SensingTypeTypeWaterDetector     SensingTypeType = "waterDetector"
	SensingTypeTypeGasDetector       SensingTypeType = "gasDetector"
	SensingTypeTypeAlarmSensor       SensingTypeType = "alarmSensor"
	SensingTypeTypePowerAlarmSensor  SensingTypeType = "powerAlarmSensor"
	SensingTypeTypeDayNightIndicator SensingTypeType = "dayNightIndicator"
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
