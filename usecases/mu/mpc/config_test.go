package mpc

import (
	"github.com/stretchr/testify/assert"
)

func (s *MuMPCSuite) Test_SupportsPhases() {
	allowedConstellations := map[ConnectedPhases][][]string{
		ConnectedPhasesA:   {{"a"}},
		ConnectedPhasesB:   {{"b"}},
		ConnectedPhasesC:   {{"C"}},
		ConnectedPhasesAB:  {{"a"}, {"b"}, {"a", "b"}},
		ConnectedPhasesBC:  {{"b"}, {"c"}, {"B", "c"}},
		ConnectedPhasesCA:  {{"a"}, {"c"}, {"A", "C"}},
		ConnectedPhasesABC: {{"a"}, {"b"}, {"c"}, {"a", "b"}, {"b", "c"}, {"a", "c"}, {"A", "b", "c"}},
	}

	for constellation, phases := range allowedConstellations {
		config := MonitorPowerConfig{
			ConnectedPhases: constellation,
		}

		for _, phase := range phases {
			assert.True(s.T(), config.SupportsPhases(phase))
		}
	}

	notAllowedConstellations := map[ConnectedPhases][]string{
		ConnectedPhasesA:  {"b", "c", "ab", "bc", "ac", "abc"},
		ConnectedPhasesB:  {"a", "c", "ab", "bc", "ac", "abc"},
		ConnectedPhasesC:  {"a", "b", "ab", "bc", "ac", "abc"},
		ConnectedPhasesAB: {"c", "ac", "abc"},
		ConnectedPhasesBC: {"a", "ab", "abc"},
		ConnectedPhasesCA: {"b", "bc", "abc"},
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
