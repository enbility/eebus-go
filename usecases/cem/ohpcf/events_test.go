package ohpcf

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"time"
)

func (s *CemOhPCFSuite) Test_Events() {
	payload := spineapi.EventPayload{
		Entity: s.mockRemoteEntity,
	}
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeEntityChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.ChangeType = spineapi.ElementChangeRemove
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeAdd
	s.sut.HandleEvent(payload)

	payload.EventType = spineapi.EventTypeDataChange
	payload.ChangeType = spineapi.ElementChangeUpdate
	payload.Data = util.Ptr(model.SmartEnergyManagementPsDataType{})
	s.sut.HandleEvent(payload)
}

func (s *CemOhPCFSuite) Test_loadSmartEnergyManagementPsDataType() {
	payload := spineapi.EventPayload{
		Ski:    remoteSki,
		Device: s.remoteDevice,
		Entity: s.mockRemoteEntity,
	}

	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeInactive),
				},
			}},
		}},
	}
	payload.Data = data
	s.sut.loadSmartEnergyManagementPsDataType(payload)
	assert.True(s.T(), s.eventCalled)

	s.eventCalled = false

	data = &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				PowerTimeSlot: []model.SmartEnergyManagementPsPowerTimeSlotType{{
					ValueList: &model.SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []model.PowerTimeSlotValueDataType{{
							Value:     model.NewScaledNumberType(1004),
							ValueType: util.Ptr(model.PowerTimeSlotValueTypeTypePower),
						}},
					},
				}},
			}},
		}},
	}
	payload.Data = data
	s.sut.loadSmartEnergyManagementPsDataType(payload)
	assert.True(s.T(), s.eventCalled)

	s.eventCalled = false

	data = &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				PowerTimeSlot: []model.SmartEnergyManagementPsPowerTimeSlotType{{
					ValueList: &model.SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []model.PowerTimeSlotValueDataType{{
							Value:     model.NewScaledNumberType(10432),
							ValueType: util.Ptr(model.PowerTimeSlotValueTypeTypePowerMax),
						}},
					},
				}},
			}},
		}},
	}
	payload.Data = data
	s.sut.loadSmartEnergyManagementPsDataType(payload)
	assert.True(s.T(), s.eventCalled)

	s.eventCalled = false

	data = &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				OperatingConstraintsInterrupt: &model.OperatingConstraintsInterruptDataType{
					IsStoppable: util.Ptr(true),
				},
			}},
		}},
	}
	payload.Data = data
	s.sut.loadSmartEnergyManagementPsDataType(payload)
	assert.True(s.T(), s.eventCalled)

	s.eventCalled = false

	data = &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				OperatingConstraintsInterrupt: &model.OperatingConstraintsInterruptDataType{
					IsPausable: util.Ptr(true),
				},
			}},
		}},
	}
	payload.Data = data
	s.sut.loadSmartEnergyManagementPsDataType(payload)
	assert.True(s.T(), s.eventCalled)

	s.eventCalled = false

	data = &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				Schedule: &model.PowerSequenceScheduleDataType{
					StartTime: model.NewAbsoluteOrRelativeTimeTypeFromTime(time.Time{}),
				},
			}},
		}},
	}
	payload.Data = data
	s.sut.loadSmartEnergyManagementPsDataType(payload)
	assert.True(s.T(), s.eventCalled)

	s.eventCalled = false

	data = &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeInvalid),
				},
			}},
		}},
	}
	payload.Data = data
	s.sut.loadSmartEnergyManagementPsDataType(payload)
	assert.True(s.T(), s.eventCalled)

	s.eventCalled = false

	data = &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				OperatingConstraintsDuration: &model.OperatingConstraintsDurationDataType{
					ActiveDurationMin: model.NewDurationType(1000),
				},
			}},
		}},
	}
	payload.Data = data
	s.sut.loadSmartEnergyManagementPsDataType(payload)
	assert.True(s.T(), s.eventCalled)

	s.eventCalled = false

	data = &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				OperatingConstraintsDuration: &model.OperatingConstraintsDurationDataType{
					PauseDurationMin: model.NewDurationType(1000),
				},
			}},
		}},
	}
	payload.Data = data
	s.sut.loadSmartEnergyManagementPsDataType(payload)
	assert.True(s.T(), s.eventCalled)
}
