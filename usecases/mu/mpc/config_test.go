package mpc

import (
	"github.com/enbility/spine-go/model"
	"github.com/stretchr/testify/assert"
)

func (s *MuMPCSuite) Test_SupportsPhases() {
	allowedConstellations := map[model.ElectricalConnectionPhaseNameType][][]string{
		model.ElectricalConnectionPhaseNameTypeA:   {{"a"}},
		model.ElectricalConnectionPhaseNameTypeB:   {{"b"}},
		model.ElectricalConnectionPhaseNameTypeC:   {{"C"}},
		model.ElectricalConnectionPhaseNameTypeAb:  {{"a"}, {"b"}, {"a", "b"}},
		model.ElectricalConnectionPhaseNameTypeBc:  {{"b"}, {"c"}, {"B", "c"}},
		model.ElectricalConnectionPhaseNameTypeAc:  {{"a"}, {"c"}, {"A", "C"}},
		model.ElectricalConnectionPhaseNameTypeAbc: {{"a"}, {"b"}, {"c"}, {"a", "b"}, {"b", "c"}, {"a", "c"}, {"A", "b", "c"}},
	}

	for constellation, phases := range allowedConstellations {
		config := MonitorPowerConfig{
			ConnectedPhases: constellation,
		}

		for _, phase := range phases {
			assert.True(s.T(), config.SupportsPhases(phase))
		}
	}

	notAllowedConstellations := map[model.ElectricalConnectionPhaseNameType][]string{
		model.ElectricalConnectionPhaseNameTypeA:  {"b", "c", "ab", "bc", "ac", "abc"},
		model.ElectricalConnectionPhaseNameTypeB:  {"a", "c", "ab", "bc", "ac", "abc"},
		model.ElectricalConnectionPhaseNameTypeC:  {"a", "b", "ab", "bc", "ac", "abc"},
		model.ElectricalConnectionPhaseNameTypeAb: {"c", "ac", "abc"},
		model.ElectricalConnectionPhaseNameTypeBc: {"a", "ab", "abc"},
		model.ElectricalConnectionPhaseNameTypeAc: {"b", "bc", "abc"},
	}

	for constellation, notSupportedPhases := range notAllowedConstellations {
		config := MonitorPowerConfig{
			ConnectedPhases: constellation,
		}

		for _, phase := range notSupportedPhases {
			assert.False(s.T(), config.SupportsPhases([]string{phase}))
		}
	}
}
