package mpc

import (
	"strings"

	"github.com/enbility/spine-go/model"
)

type PhaseMeasurementSourceMap map[model.ElectricalConnectionPhaseNameType]*model.MeasurementValueSourceType
type PhaseMeasurementConstraintsMap map[model.ElectricalConnectionPhaseNameType]*model.MeasurementConstraintsDataType

// MonitorPowerConfig is the configuration for the monitor use case
// This config is required by the mpc use case and must be used in mpc.NewMPC
type MonitorPowerConfig struct {
	ConnectedPhases model.ElectricalConnectionPhaseNameType // The phases that are measured

	ValueSourceTotal    *model.MeasurementValueSourceType // The source of the values from the acPowerTotal (required)
	ValueSourcePerPhase PhaseMeasurementSourceMap         // The source of the values from the acPower per phase (required if the phase is supported)

	ValueConstraintsTotal    *model.MeasurementConstraintsDataType // The constraints for the acPowerTotal (optional can be nil)
	ValueConstraintsPerPhase PhaseMeasurementConstraintsMap        // The constraints for the acPower per phase (optional can be nil)
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
	ValueSourcePerPhase      PhaseMeasurementSourceMap      // The source of the values per phase (required if the phase is supported)
	ValueConstraintsPerPhase PhaseMeasurementConstraintsMap // The constraints for the current per phase (optional can be nil) (requires ValueSourcePerPhase to be set)
}

// MonitorVoltageConfig is the configuration for the monitor use case
// If this config is passed via NewMPC, the use case will support voltage monitoring
// The voltage phases will be the same as specified in MonitorPowerConfig
type MonitorVoltageConfig struct {
	ValueSourcePerPhase      PhaseMeasurementSourceMap      // The source of the values per phase (required if the phase is supported)
	ValueConstraintsPerPhase PhaseMeasurementConstraintsMap // The constraints for the voltage per phase (optional can be nil) (requires ValueSourcePerPhase to be set)
}

// MonitorFrequencyConfig is the configuration for the monitor use case
type MonitorFrequencyConfig struct {
	ValueSource      *model.MeasurementValueSourceType     // The source of the values (required)
	ValueConstraints *model.MeasurementConstraintsDataType // The constraints for the frequency values (optional can be nil)
}

// SupportsPhases checks if the config supports the given phases
// e.g. SupportsPhases([]string{"a", "B"}) will return true if the config has ConnectedPhases set to "ab" or "abc"
func (c *MonitorPowerConfig) SupportsPhases(phase []string) bool {
	phasesString := string(c.ConnectedPhases)
	supports := true
	for _, p := range phase {
		if !strings.Contains(strings.ToLower(phasesString), strings.ToLower(p)) {
			supports = false
			break
		}
	}
	return supports
}
