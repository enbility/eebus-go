package mpc

import (
	"strings"

	"github.com/enbility/spine-go/model"
)

type ConnectedPhases string

const (
	ConnectedPhasesA   ConnectedPhases = "a"
	ConnectedPhasesB   ConnectedPhases = "b"
	ConnectedPhasesC   ConnectedPhases = "c"
	ConnectedPhasesAB  ConnectedPhases = "ab"
	ConnectedPhasesBC  ConnectedPhases = "bc"
	ConnectedPhasesCA  ConnectedPhases = "ac"
	ConnectedPhasesABC ConnectedPhases = "abc"
)

// MonitorPowerConfig is the configuration for the monitor use case
// This config is required by the mpc use case and must be used in mpc.NewMPC
type MonitorPowerConfig struct {
	ConnectedPhases ConnectedPhases // The phases that are measured

	ValueSourceTotal  *model.MeasurementValueSourceType // The source of the values from the acPowerTotal (required)
	ValueSourcePhaseA *model.MeasurementValueSourceType // The source of the values from the acPower for phase A (required if the phase is supported)
	ValueSourcePhaseB *model.MeasurementValueSourceType // The source of the values from the acPower for phase B (required if the phase is supported)
	ValueSourcePhaseC *model.MeasurementValueSourceType // The source of the values from the acPower for phase C (required if the phase is supported)

	ValueConstraintsTotal  *model.MeasurementConstraintsDataType // The constraints for the acPowerTotal (optional can be nil)
	ValueConstraintsPhaseA *model.MeasurementConstraintsDataType // The constraints for the acPower for phase A (optional can be nil)
	ValueConstraintsPhaseB *model.MeasurementConstraintsDataType // The constraints for the acPower for phase B (optional can be nil)
	ValueConstraintsPhaseC *model.MeasurementConstraintsDataType // The constraints for the acPower for phase C (optional can be nil)
}

// MonitorEnergyConfig is the configuration for the monitor use case
// If this config is passed via NewMPC, the use case will support energy monitoring as specified
type MonitorEnergyConfig struct {
	ValueSourceProduction      *model.MeasurementValueSourceType     // The source of the production values (if this is set, the use case will support production) (optional can be nil)
	ValueConstraintsProduction *model.MeasurementConstraintsDataType // The constraints for the production values (optional can be nil) (requires ProductionValueSource to be set)

	ValueSourceConsumption      *model.MeasurementValueSourceType     // The source of the consumption values (if this is set, the use case will support consumption) (optional can be nil)
	ValueConstraintsConsumption *model.MeasurementConstraintsDataType // The constraints for the consumption values (optional can be nil) (requires ConsumptionValueSource to be set)
}

// MonitorCurrentConfig is the configuration for the monitor use case
// If this config is passed via NewMPC, the use case will support current monitoring
// The current phases will be the same as specified in MonitorPowerConfig
type MonitorCurrentConfig struct {
	ValueSourcePhaseA *model.MeasurementValueSourceType // The source of the values for phase A (required if the phase is supported)
	ValueSourcePhaseB *model.MeasurementValueSourceType // The source of the values for phase B (required if the phase is supported)
	ValueSourcePhaseC *model.MeasurementValueSourceType // The source of the values for phase C (required if the phase is supported)

	ValueConstraintsPhaseA *model.MeasurementConstraintsDataType // The constraints for the current for phase A (optional can be nil) (requires ValueSourcePhaseA to be set)
	ValueConstraintsPhaseB *model.MeasurementConstraintsDataType // The constraints for the current for phase B (optional can be nil) (requires ValueSourcePhaseB to be set)
	ValueConstraintsPhaseC *model.MeasurementConstraintsDataType // The constraints for the current for phase C (optional can be nil) (requires ValueSourcePhaseC to be set)
}

// MonitorVoltageConfig is the configuration for the monitor use case
// If this config is passed via NewMPC, the use case will support voltage monitoring
// The voltage phases will be the same as specified in MonitorPowerConfig
type MonitorVoltageConfig struct {
	ValueSourcePhaseA *model.MeasurementValueSourceType // The source of the values for phase A (required if the phase is supported)
	ValueSourcePhaseB *model.MeasurementValueSourceType // The source of the values for phase B (required if the phase is supported)
	ValueSourcePhaseC *model.MeasurementValueSourceType // The source of the values for phase C (required if the phase is supported)

	ValueConstraintsPhaseA *model.MeasurementConstraintsDataType // The constraints for the voltage for phase A (optional can be nil) (requires ValueSourcePhaseA to be set)
	ValueConstraintsPhaseB *model.MeasurementConstraintsDataType // The constraints for the voltage for phase B (optional can be nil) (requires ValueSourcePhaseB to be set)
	ValueConstraintsPhaseC *model.MeasurementConstraintsDataType // The constraints for the voltage for phase C (optional can be nil) (requires ValueSourcePhaseC to be set)

	SupportPhaseToPhase  bool                              // Set to true if the use case supports phase to phase voltage monitoring
	ValueSourcePhaseAToB *model.MeasurementValueSourceType // The source of the values for phase A to B (required if the phases are supported and SupportPhaseToPhase is true)
	ValueSourcePhaseBToC *model.MeasurementValueSourceType // The source of the values for phase B to C (required if the phases are supported and SupportPhaseToPhase is true)
	ValueSourcePhaseCToA *model.MeasurementValueSourceType // The source of the values for phase C to A (required if the phases are supported and SupportPhaseToPhase is true)

	ValueConstraintsPhaseAToB *model.MeasurementConstraintsDataType // The constraints for the voltage for phase A to B (optional can be nil) (requires ValueSourcePhaseAToB to be set)
	ValueConstraintsPhaseBToC *model.MeasurementConstraintsDataType // The constraints for the voltage for phase B to C (optional can be nil) (requires ValueSourcePhaseBToC to be set)
	ValueConstraintsPhaseCToA *model.MeasurementConstraintsDataType // The constraints for the voltage for phase C to A (optional can be nil) (requires ValueSourcePhaseCToA to be set)
}

// MonitorFrequencyConfig is the configuration for the monitor use case
type MonitorFrequencyConfig struct {
	ValueSource      *model.MeasurementValueSourceType     // The source of the values (required)
	ValueConstraints *model.MeasurementConstraintsDataType // The constraints for the frequency values (optional can be nil)
}

func (c *MonitorPowerConfig) SupportsPhases(phase []string) bool {
	phasesString := string(c.ConnectedPhases)
	supports := true
	for _, p := range phase {
		if !strings.Contains(phasesString, p) {
			supports = false
			break
		}
	}
	return supports
}
