package mpc

import (
	"github.com/enbility/spine-go/model"
	"strings"
)

type ConnectedPhases string

const ConnectedPhasesA ConnectedPhases = "a"
const ConnectedPhasesB ConnectedPhases = "b"
const ConnectedPhasesC ConnectedPhases = "c"
const ConnectedPhasesAB ConnectedPhases = "ab"
const ConnectedPhasesBC ConnectedPhases = "bc"
const ConnectedPhasesCA ConnectedPhases = "ac"
const ConnectedPhasesABC ConnectedPhases = "abc"

// MonitorPowerConfig is the configuration for the monitor use case
// This config is required by the mpc use case and must be used in mpc.NewMPC
type MonitorPowerConfig struct {
	ConnectedPhases ConnectedPhases // The phases that are measured

	ValueSourceTotal  *model.MeasurementValueSourceType // The source of the values from the acPowerTotal (not optional)
	ValueSourcePhaseA *model.MeasurementValueSourceType // The source of the values from the acPower for phase A (shall be set if the phase is supported)
	ValueSourcePhaseB *model.MeasurementValueSourceType // The source of the values from the acPower for phase B (shall be set if the phase is supported)
	ValueSourcePhaseC *model.MeasurementValueSourceType // The source of the values from the acPower for phase C (shall be set if the phase is supported)

	ValueConstraintsTotal  *model.MeasurementConstraintsDataType // The constraints for the acPowerTotal (optional can be nil)
	ValueConstraintsPhaseA *model.MeasurementConstraintsDataType // The constraints for the acPower for phase A (optional can be nil)
	ValueConstraintsPhaseB *model.MeasurementConstraintsDataType // The constraints for the acPower for phase B (optional can be nil)
	ValueConstraintsPhaseC *model.MeasurementConstraintsDataType // The constraints for the acPower for phase C (optional can be nil)
}

// MonitorEnergyConfig is the configuration for the monitor use case
// If this config is passed via NewMPC, the use case will support energy monitoring as specified
type MonitorEnergyConfig struct {
	ValueSourceProduction      *model.MeasurementValueSourceType     // The source of the production values (if this is set, the use case will support production) (optional can be nil)
	ValueConstraintsProduction *model.MeasurementConstraintsDataType // The constraints for the production values (optional can be nil) (needs ProductionValueSource to be set)

	ValueSourceConsumption      *model.MeasurementValueSourceType     // The source of the consumption values (if this is set, the use case will support consumption) (optional can be nil)
	ValueConstraintsConsumption *model.MeasurementConstraintsDataType // The constraints for the consumption values (optional can be nil) (needs ConsumptionValueSource to be set)
}

// MonitorCurrentConfig is the configuration for the monitor use case
// If this config is passed via NewMPC, the use case will support current monitoring
// The current phases will be the same as specified in MonitorPowerConfig
type MonitorCurrentConfig struct {
	ValueSourcePhaseA *model.MeasurementValueSourceType // The source of the values for phase A (shall be set if the phase is supported)
	ValueSourcePhaseB *model.MeasurementValueSourceType // The source of the values for phase B (shall be set if the phase is supported)
	ValueSourcePhaseC *model.MeasurementValueSourceType // The source of the values for phase C (shall be set if the phase is supported)

	ValueConstraintsPhaseA *model.MeasurementConstraintsDataType // The constraints for the current for phase A (optional can be nil) (needs ValueSourcePhaseA to be set)
	ValueConstraintsPhaseB *model.MeasurementConstraintsDataType // The constraints for the current for phase B (optional can be nil) (needs ValueSourcePhaseB to be set)
	ValueConstraintsPhaseC *model.MeasurementConstraintsDataType // The constraints for the current for phase C (optional can be nil) (needs ValueSourcePhaseC to be set)
}

// MonitorVoltageConfig is the configuration for the monitor use case
// If this config is passed via NewMPC, the use case will support voltage monitoring
// The voltage phases will be the same as specified in MonitorPowerConfig
type MonitorVoltageConfig struct {
	ValueSourcePhaseA *model.MeasurementValueSourceType // The source of the values for phase A (shall be set if the phase is supported)
	ValueSourcePhaseB *model.MeasurementValueSourceType // The source of the values for phase B (shall be set if the phase is supported)
	ValueSourcePhaseC *model.MeasurementValueSourceType // The source of the values for phase C (shall be set if the phase is supported)

	ValueConstraintsPhaseA *model.MeasurementConstraintsDataType // The constraints for the voltage for phase A (optional can be nil) (needs ValueSourcePhaseA to be set)
	ValueConstraintsPhaseB *model.MeasurementConstraintsDataType // The constraints for the voltage for phase B (optional can be nil) (needs ValueSourcePhaseB to be set)
	ValueConstraintsPhaseC *model.MeasurementConstraintsDataType // The constraints for the voltage for phase C (optional can be nil) (needs ValueSourcePhaseC to be set)

	SupportPhaseToPhase  bool                              // If the use case shall support phase to phase voltage monitoring
	ValueSourcePhaseAToB *model.MeasurementValueSourceType // The source of the values for phase A to B (shall be set if the phases are supported and SupportPhaseToPhase is true)
	ValueSourcePhaseBToC *model.MeasurementValueSourceType // The source of the values for phase B to C (shall be set if the phases are supported and SupportPhaseToPhase is true)
	ValueSourcePhaseCToA *model.MeasurementValueSourceType // The source of the values for phase C to A (shall be set if the phases are supported and SupportPhaseToPhase is true)

	ValueConstraintsPhaseAToB *model.MeasurementConstraintsDataType // The constraints for the voltage for phase A to B (optional can be nil) (needs ValueSourcePhaseAToB to be set)
	ValueConstraintsPhaseBToC *model.MeasurementConstraintsDataType // The constraints for the voltage for phase B to C (optional can be nil) (needs ValueSourcePhaseBToC to be set)
	ValueConstraintsPhaseCToA *model.MeasurementConstraintsDataType // The constraints for the voltage for phase C to A (optional can be nil) (needs ValueSourcePhaseCToA to be set)
}

// MonitorFrequencyConfig is the configuration for the monitor use case
type MonitorFrequencyConfig struct {
	ValueSource      *model.MeasurementValueSourceType     // The source of the values (not optional)
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
