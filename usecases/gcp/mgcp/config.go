package mgcp

import "github.com/enbility/spine-go/model"

// MonitorPvFeedInPowerLimitationFactorConfig is the configuration for the power limitation factor monitoring use case in the MGCP
// If this config is passed via NewMGCP, the MGCP use case will support power limitation factor monitoring
type MonitorPvFeedInPowerLimitationFactorConfig struct {
}

// MonitorPowerConfig is the configuration for the power monitoring use case in the MGCP
// This config is required by the MGCP use case and must be used in NewMGCP
type MonitorPowerConfig struct {
	ValueSource      *model.MeasurementValueSourceType     // The source of the values from the acPowerTotal (not optional)
	ValueConstraints *model.MeasurementConstraintsDataType // The constraints for the acPowerTotal (optional, can be nil)
}

// MonitorEnergyConfig is the configuration for the energy monitoring use case in the MGCP
// This config is required by the MGCP use case and must be used in NewMGCP
type MonitorEnergyConfig struct {
	ValueSourceProduction  *model.MeasurementValueSourceType // The source of the production values (not optional)
	ValueSourceConsumption *model.MeasurementValueSourceType // The source of the consumption values (not optional)

	ValueConstraintsProduction  *model.MeasurementConstraintsDataType // The constraints for the production values (optional, can be nil)
	ValueConstraintsConsumption *model.MeasurementConstraintsDataType // The constraints for the consumption values (optional, can be nil)
}

// MonitorCurrentConfig is the configuration for the current monitoring use case in the MGCP
// If this config is passed via NewMGCP, the MGCP use case will support current monitoring
type MonitorCurrentConfig struct {
	ValueSourcePhaseA *model.MeasurementValueSourceType // The source of the values for phase A (not optional)
	ValueSourcePhaseB *model.MeasurementValueSourceType // The source of the values for phase B (not optional)
	ValueSourcePhaseC *model.MeasurementValueSourceType // The source of the values for phase C (not optional)

	ValueConstraintsPhaseA *model.MeasurementConstraintsDataType // The constraints for the current for phase A (optional, can be nil)
	ValueConstraintsPhaseB *model.MeasurementConstraintsDataType // The constraints for the current for phase B (optional, can be nil)
	ValueConstraintsPhaseC *model.MeasurementConstraintsDataType // The constraints for the current for phase C (optional, can be nil)
}

// MonitorVoltageConfig is the configuration for the voltage monitoring use case in the MGCP
// If this config is passed via NewMGCP, the MGCP use case will support voltage monitoring
type MonitorVoltageConfig struct {
	// If the value source is not nil, the use case will support the voltage monitoring for the respective phase
	// If the value source is nil, the use case will not support the voltage monitoring for the respective phase
	ValueSourcePhaseA    *model.MeasurementValueSourceType // The source of the values for phase A (optional, can be nil)
	ValueSourcePhaseB    *model.MeasurementValueSourceType // The source of the values for phase B (optional, can be nil)
	ValueSourcePhaseC    *model.MeasurementValueSourceType // The source of the values for phase C (optional, can be nil)
	ValueSourcePhaseAToB *model.MeasurementValueSourceType // The source of the values for phase A to B (optional, can be nil)
	ValueSourcePhaseBToC *model.MeasurementValueSourceType // The source of the values for phase B to C (optional, can be nil)
	ValueSourcePhaseCToA *model.MeasurementValueSourceType // The source of the values for phase C to A (optional, can be nil)

	ValueConstraintsPhaseA    *model.MeasurementConstraintsDataType // The constraints for the voltage for phase A (optional, can be nil) (needs ValueSourcePhaseA to be set)
	ValueConstraintsPhaseB    *model.MeasurementConstraintsDataType // The constraints for the voltage for phase B (optional, can be nil) (needs ValueSourcePhaseB to be set)
	ValueConstraintsPhaseC    *model.MeasurementConstraintsDataType // The constraints for the voltage for phase C (optional, can be nil) (needs ValueSourcePhaseC to be set)
	ValueConstraintsPhaseAToB *model.MeasurementConstraintsDataType // The constraints for the voltage for phase A to B (optional, can be nil) (needs ValueSourcePhaseAToB to be set)
	ValueConstraintsPhaseBToC *model.MeasurementConstraintsDataType // The constraints for the voltage for phase B to C (optional, can be nil) (needs ValueSourcePhaseBToC to be set)
	ValueConstraintsPhaseCToA *model.MeasurementConstraintsDataType // The constraints for the voltage for phase C to A (optional, can be nil) (needs ValueSourcePhaseCToA to be set)
}

// MonitorFrequencyConfig is the configuration for the frequency monitoring use case in the MGCP
// If this config is passed via NewMGCP, the MGCP use case will support frequency monitoring
type MonitorFrequencyConfig struct {
	ValueSource      *model.MeasurementValueSourceType     // The source of the values (not optional)
	ValueConstraints *model.MeasurementConstraintsDataType // The constraints for the frequency values (optional can be nil)
}
