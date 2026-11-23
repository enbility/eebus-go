package ohpcf

import (
	"github.com/enbility/eebus-go/usecases/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
	"time"
)

// Scenario 1

func (s *CemOhPCFSuite) Test_OptionalPowerConsumptionAvailable() {
	_, err := s.sut.OptionalPowerConsumptionAvailable(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.OptionalPowerConsumptionAvailable(s.monitoredEntity)
	assert.NotNil(s.T(), err)

	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeInactive),
				},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	available, err := s.sut.OptionalPowerConsumptionAvailable(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, available)
}

func (s *CemOhPCFSuite) Test_Power() {
	_, err := s.sut.RequestedPowerEstimate(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.RequestedPowerEstimate(s.monitoredEntity)
	assert.NotNil(s.T(), err)

	data := &model.SmartEnergyManagementPsDataType{
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

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	available, err := s.sut.RequestedPowerEstimate(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1004.0, available)
}

func (s *CemOhPCFSuite) Test_MaxPower() {
	_, err := s.sut.RequestedPowerMax(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.RequestedPowerMax(s.monitoredEntity)
	assert.NotNil(s.T(), err)

	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				PowerTimeSlot: []model.SmartEnergyManagementPsPowerTimeSlotType{{
					ValueList: &model.SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []model.PowerTimeSlotValueDataType{{
							Value:     model.NewScaledNumberType(1006),
							ValueType: util.Ptr(model.PowerTimeSlotValueTypeTypePowerMax),
						}},
					},
				}},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	available, err := s.sut.RequestedPowerMax(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), 1006.0, available)
}

func (s *CemOhPCFSuite) Test_ConsumptionIsStoppable() {
	_, err := s.sut.ConsumptionIsStoppable(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.ConsumptionIsStoppable(s.monitoredEntity)
	assert.NotNil(s.T(), err)

	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				OperatingConstraintsInterrupt: &model.OperatingConstraintsInterruptDataType{
					IsStoppable: util.Ptr(true),
				},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	stoppable, err := s.sut.ConsumptionIsStoppable(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, stoppable)
}

func (s *CemOhPCFSuite) Test_ConsumptionIsPausable() {
	_, err := s.sut.ConsumptionIsPausable(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.ConsumptionIsPausable(s.monitoredEntity)
	assert.NotNil(s.T(), err)

	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				OperatingConstraintsInterrupt: &model.OperatingConstraintsInterruptDataType{
					IsPausable: util.Ptr(true),
				},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	pausable, err := s.sut.ConsumptionIsPausable(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, pausable)
}

func (s *CemOhPCFSuite) Test_PowerConsumptionProcessStartTime() {
	_, err := s.sut.PowerConsumptionProcessStartTime(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.PowerConsumptionProcessStartTime(s.monitoredEntity)
	assert.NotNil(s.T(), err)

	utcNow := time.Now().UTC()
	utcNowTimeObj := model.NewAbsoluteOrRelativeTimeTypeFromTime(utcNow)

	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				Schedule: &model.PowerSequenceScheduleDataType{
					StartTime: utcNowTimeObj,
				},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	startTime, err := s.sut.PowerConsumptionProcessStartTime(s.monitoredEntity)
	assert.Nil(s.T(), err)
	expected, _ := utcNowTimeObj.GetTime()
	assert.Equal(s.T(), expected, startTime)
}

func (s *CemOhPCFSuite) Test_PowerConsumptionProcessState() {
	_, err := s.sut.PowerConsumptionProcessState(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.PowerConsumptionProcessState(s.monitoredEntity)
	assert.NotNil(s.T(), err)

	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeInactive),
				},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	state, err := s.sut.PowerConsumptionProcessState(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), api.CompressorPowerConsumptionStateAvailable, state)

	// A paused process is written as "paused" [OHPCF-012] and must read back as paused
	for _, pausedState := range []model.PowerSequenceStateType{
		model.PowerSequenceStateTypePaused,
		model.PowerSequenceStateTypeScheduledPaused,
	} {
		data.Alternatives[0].PowerSequence[0].State.State = util.Ptr(pausedState)
		_, fErr = rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
		assert.Nil(s.T(), fErr)

		state, err = s.sut.PowerConsumptionProcessState(s.monitoredEntity)
		assert.Nil(s.T(), err)
		assert.Equal(s.T(), api.CompressorPowerConsumptionStatePaused, state)
	}
}

func (s *CemOhPCFSuite) Test_PowerConsumptionMinimalRunDuration() {
	_, err := s.sut.PowerConsumptionMinimalRunDuration(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.PowerConsumptionMinimalRunDuration(s.monitoredEntity)
	assert.NotNil(s.T(), err)

	duration := time.Duration(120000000000000000)

	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				OperatingConstraintsDuration: &model.OperatingConstraintsDurationDataType{
					ActiveDurationMin: model.NewDurationType(duration),
				},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	minDuration, err := s.sut.PowerConsumptionMinimalRunDuration(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), duration, minDuration)
}

func (s *CemOhPCFSuite) Test_PowerConsumptionMinimalPauseDuration() {
	_, err := s.sut.PowerConsumptionMinimalPauseDuration(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.PowerConsumptionMinimalPauseDuration(s.monitoredEntity)
	assert.NotNil(s.T(), err)

	duration := time.Duration(120000000000000000)

	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				OperatingConstraintsDuration: &model.OperatingConstraintsDurationDataType{
					PauseDurationMin: model.NewDurationType(duration),
				},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	minDuration, err := s.sut.PowerConsumptionMinimalPauseDuration(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), duration, minDuration)
}

// Scenario 2

func (s *CemOhPCFSuite) Test_SchedulePowerConsumptionProcess() {
	_, err := s.sut.SchedulePowerConsumptionProcess(s.mockRemoteEntity, 0, nil)
	assert.NotNil(s.T(), err)

	// Without valid data, the call should fail
	_, err = s.sut.SchedulePowerConsumptionProcess(s.monitoredEntity, 0, nil)
	assert.NotNil(s.T(), err)

	// Set up valid SmartEnergyManagementPs data
	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				Description: &model.PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				},
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeInactive),
				},
				OperatingConstraintsInterrupt: &model.OperatingConstraintsInterruptDataType{
					IsPausable:  util.Ptr(true),
					IsStoppable: util.Ptr(true),
				},
				PowerTimeSlot: []model.SmartEnergyManagementPsPowerTimeSlotType{{
					ValueList: &model.SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []model.PowerTimeSlotValueDataType{{
							Value:     model.NewScaledNumberType(1000),
							ValueType: util.Ptr(model.PowerTimeSlotValueTypeTypePower),
						}},
					},
				}},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	msgCounter, err := s.sut.SchedulePowerConsumptionProcess(s.monitoredEntity, time.Hour, nil)
	assert.NotNil(s.T(), msgCounter)
	assert.Nil(s.T(), err)
}

func (s *CemOhPCFSuite) Test_StopAbortPowerConsumptionProcess() {
	_, err := s.sut.AbortPowerConsumptionProcess(s.mockRemoteEntity, nil)
	assert.NotNil(s.T(), err)

	// Without valid data, the call should fail
	_, err = s.sut.AbortPowerConsumptionProcess(s.monitoredEntity, nil)
	assert.NotNil(s.T(), err)

	// Set up valid SmartEnergyManagementPs data
	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				Description: &model.PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				},
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeRunning),
				},
				OperatingConstraintsInterrupt: &model.OperatingConstraintsInterruptDataType{
					IsPausable:  util.Ptr(true),
					IsStoppable: util.Ptr(true),
				},
				PowerTimeSlot: []model.SmartEnergyManagementPsPowerTimeSlotType{{
					ValueList: &model.SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []model.PowerTimeSlotValueDataType{{
							Value:     model.NewScaledNumberType(1000),
							ValueType: util.Ptr(model.PowerTimeSlotValueTypeTypePower),
						}},
					},
				}},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	msgCounter, err := s.sut.AbortPowerConsumptionProcess(s.monitoredEntity, nil)
	assert.NotNil(s.T(), msgCounter)
	assert.Nil(s.T(), err)
}

func (s *CemOhPCFSuite) Test_PausePowerConsumptionProcess() {
	_, err := s.sut.PausePowerConsumptionProcess(s.mockRemoteEntity, nil)
	assert.NotNil(s.T(), err)

	// Without valid data, the call should fail
	_, err = s.sut.PausePowerConsumptionProcess(s.monitoredEntity, nil)
	assert.NotNil(s.T(), err)

	// Set up valid SmartEnergyManagementPs data
	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				Description: &model.PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				},
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeRunning),
				},
				OperatingConstraintsInterrupt: &model.OperatingConstraintsInterruptDataType{
					IsPausable:  util.Ptr(true),
					IsStoppable: util.Ptr(true),
				},
				PowerTimeSlot: []model.SmartEnergyManagementPsPowerTimeSlotType{{
					ValueList: &model.SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []model.PowerTimeSlotValueDataType{{
							Value:     model.NewScaledNumberType(1000),
							ValueType: util.Ptr(model.PowerTimeSlotValueTypeTypePower),
						}},
					},
				}},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	msgCounter, err := s.sut.PausePowerConsumptionProcess(s.monitoredEntity, nil)
	assert.NotNil(s.T(), msgCounter)
	assert.Nil(s.T(), err)
}

func (s *CemOhPCFSuite) Test_ResumePowerConsumptionProcess() {
	_, err := s.sut.ResumePowerConsumptionProcess(s.mockRemoteEntity, nil)
	assert.NotNil(s.T(), err)

	// Without valid data, the call should fail
	_, err = s.sut.ResumePowerConsumptionProcess(s.monitoredEntity, nil)
	assert.NotNil(s.T(), err)

	// Set up valid SmartEnergyManagementPs data with paused state
	data := &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{
				Description: &model.PowerSequenceDescriptionDataType{
					SequenceId: util.Ptr(model.PowerSequenceIdType(1)),
				},
				State: &model.PowerSequenceStateDataType{
					State: util.Ptr(model.PowerSequenceStateTypeScheduledPaused),
				},
				OperatingConstraintsInterrupt: &model.OperatingConstraintsInterruptDataType{
					IsPausable:  util.Ptr(true),
					IsStoppable: util.Ptr(true),
				},
				PowerTimeSlot: []model.SmartEnergyManagementPsPowerTimeSlotType{{
					ValueList: &model.SmartEnergyManagementPsPowerTimeSlotValueListType{
						Value: []model.PowerTimeSlotValueDataType{{
							Value:     model.NewScaledNumberType(1000),
							ValueType: util.Ptr(model.PowerTimeSlotValueTypeTypePower),
						}},
					},
				}},
			}},
		}},
	}

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, data, nil, nil)
	assert.Nil(s.T(), fErr)

	msgCounter, err := s.sut.ResumePowerConsumptionProcess(s.monitoredEntity, nil)
	assert.NotNil(s.T(), msgCounter)
	assert.Nil(s.T(), err)
}

// Helper methods

func (s *CemOhPCFSuite) Test_checkEntityTypeAndGetData() {
	_, err := s.sut.checkEntityTypeAndGetData(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)

	_, err = s.sut.checkEntityTypeAndGetData(s.monitoredEntity)
	assert.NotNil(s.T(), err)

	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, &model.SmartEnergyManagementPsDataType{}, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err := s.sut.checkEntityTypeAndGetData(s.monitoredEntity)
	assert.NotNil(s.T(), data)
	assert.Nil(s.T(), err)
}

func (s *CemOhPCFSuite) Test_isDataAvailable() {
	rFeature := s.remoteDevice.FeatureByEntityTypeAndRole(s.monitoredEntity, model.FeatureTypeTypeSmartEnergyManagementPs, model.RoleTypeServer)
	_, fErr := rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, &model.SmartEnergyManagementPsDataType{}, nil, nil)
	assert.Nil(s.T(), fErr)

	fData, err := s.sut.checkEntityTypeAndGetData(s.monitoredEntity)
	assert.NotNil(s.T(), fData)
	assert.Nil(s.T(), err)

	available := s.sut.isDataAvailable(fData)
	assert.Equal(s.T(), false, available)

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{},
	}, nil, nil)
	assert.Nil(s.T(), fErr)

	fData, err = s.sut.checkEntityTypeAndGetData(s.monitoredEntity)
	assert.NotNil(s.T(), fData)
	assert.Nil(s.T(), err)

	available = s.sut.isDataAvailable(fData)
	assert.Equal(s.T(), false, available)

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{}},
	}, nil, nil)
	assert.Nil(s.T(), fErr)

	fData, err = s.sut.checkEntityTypeAndGetData(s.monitoredEntity)
	assert.NotNil(s.T(), fData)
	assert.Nil(s.T(), err)

	available = s.sut.isDataAvailable(fData)
	assert.Equal(s.T(), false, available)

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{},
		}},
	}, nil, nil)
	assert.Nil(s.T(), fErr)

	fData, err = s.sut.checkEntityTypeAndGetData(s.monitoredEntity)
	assert.NotNil(s.T(), fData)
	assert.Nil(s.T(), err)

	available = s.sut.isDataAvailable(fData)
	assert.Equal(s.T(), false, available)

	_, fErr = rFeature.UpdateData(true, model.FunctionTypeSmartEnergyManagementPsData, &model.SmartEnergyManagementPsDataType{
		Alternatives: []model.SmartEnergyManagementPsAlternativesType{{
			PowerSequence: []model.SmartEnergyManagementPsPowerSequenceType{{}},
		}},
	}, nil, nil)
	assert.Nil(s.T(), fErr)

	fData, err = s.sut.checkEntityTypeAndGetData(s.monitoredEntity)
	assert.NotNil(s.T(), fData)
	assert.Nil(s.T(), err)

	available = s.sut.isDataAvailable(fData)
	assert.Equal(s.T(), true, available)
}
