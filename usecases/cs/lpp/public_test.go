package lpp

import (
	"time"

	"github.com/enbility/eebus-go/features/client"
	ucapi "github.com/enbility/eebus-go/usecases/api"
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CsLPPSuite) Test_LoadControlLimit() {
	limit, err := s.sut.ProductionLimit()
	assert.Equal(s.T(), 0.0, limit.Value)
	assert.Nil(s.T(), err)

	newLimit := ucapi.LoadLimit{
		Duration:     time.Duration(time.Hour * 2),
		IsActive:     true,
		IsChangeable: true,
		Value:        16,
	}
	err = s.sut.SetProductionLimit(newLimit)
	assert.Nil(s.T(), err)

	limit, err = s.sut.ProductionLimit()
	assert.Equal(s.T(), 16.0, limit.Value)
	assert.Nil(s.T(), err)
}

func (s *CsLPPSuite) Test_PendingProductionLimits() {
	data := s.sut.PendingProductionLimits()
	assert.Equal(s.T(), 0, len(data))

	msgCounter := model.MsgCounterType(500)

	msg := &spineapi.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(msgCounter),
		},
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{
				LoadControlLimitData: []model.LoadControlLimitDataType{
					{
						LimitId:       util.Ptr(model.LoadControlLimitIdType(0)),
						IsLimitActive: util.Ptr(true),
						Value:         model.NewScaledNumberType(1000),
						TimePeriod:    model.NewTimePeriodTypeWithRelativeEndTime(time.Minute * 2),
					},
					{
						LimitId:       util.Ptr(model.LoadControlLimitIdType(1)),
						IsLimitActive: util.Ptr(true),
						Value:         model.NewScaledNumberType(1000),
						TimePeriod:    model.NewTimePeriodTypeWithRelativeEndTime(time.Minute * 2),
					}},
			},
		},
		DeviceRemote: s.remoteDevice,
		EntityRemote: s.monitoredEntity,
	}

	s.sut.loadControlWriteCB(msg)

	data = s.sut.PendingProductionLimits()
	assert.Equal(s.T(), 1, len(data))

	s.sut.ApproveOrDenyProductionLimit(model.MsgCounterType(499), true, "")

	data = s.sut.PendingProductionLimits()
	assert.Equal(s.T(), 1, len(data))

	s.sut.ApproveOrDenyProductionLimit(msgCounter, false, "leave me alone")

	data = s.sut.PendingProductionLimits()
	assert.Equal(s.T(), 0, len(data))
}

func (s *CsLPPSuite) Test_Failsafe() {
	limit, changeable, err := s.sut.FailsafeProductionActivePowerLimit()
	assert.Equal(s.T(), 0.0, limit)
	assert.Equal(s.T(), true, changeable)
	assert.Nil(s.T(), err)

	err = s.sut.SetFailsafeProductionActivePowerLimit(10, false)
	assert.Nil(s.T(), err)

	limit, changeable, err = s.sut.FailsafeProductionActivePowerLimit()
	assert.Equal(s.T(), 10.0, limit)
	assert.Equal(s.T(), false, changeable)
	assert.Nil(s.T(), err)

	// The actual tests of the functionality is located in the util package
	duration, changeable, err := s.sut.FailsafeDurationMinimum()
	assert.Equal(s.T(), time.Duration(0), duration)
	assert.Equal(s.T(), true, changeable)
	assert.Nil(s.T(), err)

	err = s.sut.SetFailsafeDurationMinimum(time.Duration(time.Hour*1), true)
	assert.NotNil(s.T(), err)

	err = s.sut.SetFailsafeDurationMinimum(time.Duration(time.Hour*2), true)
	assert.Nil(s.T(), err)

	err = s.sut.SetFailsafeProductionActivePowerLimit(10, true)
	assert.Nil(s.T(), err)

	limit, changeable, err = s.sut.FailsafeProductionActivePowerLimit()
	assert.Equal(s.T(), 10.0, limit)
	assert.Equal(s.T(), true, changeable)
	assert.Nil(s.T(), err)

	duration, changeable, err = s.sut.FailsafeDurationMinimum()
	assert.Equal(s.T(), time.Duration(time.Hour*2), duration)
	assert.Equal(s.T(), true, changeable)
	assert.Nil(s.T(), err)
}

func (s *CsLPPSuite) Test_Heartbeat() {
	assert.Nil(s.T(), s.sut.heartbeatDiag)

	value := s.sut.IsHeartbeatWithinDuration()
	assert.False(s.T(), value)

	remoteDiagServer := s.monitoredEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeDeviceDiagnosis, model.RoleTypeServer)
	assert.NotNil(s.T(), remoteDiagServer)

	var err error
	s.sut.heartbeatDiag, err = client.NewDeviceDiagnosis(s.sut.LocalEntity, s.monitoredEntity)
	assert.NotNil(s.T(), s.sut.heartbeatDiag)
	assert.Nil(s.T(), err)

	// add heartbeat data to the remoteDiagServer
	timestamp := time.Now().Add(-time.Second * 121)
	data := &model.DeviceDiagnosisHeartbeatDataType{
		Timestamp:        model.NewAbsoluteOrRelativeTimeTypeFromTime(timestamp),
		HeartbeatCounter: util.Ptr(uint64(1)),
		HeartbeatTimeout: model.NewDurationType(time.Second * 120),
	}
	_, err1 := remoteDiagServer.UpdateData(true, model.FunctionTypeDeviceDiagnosisHeartbeatData, data, nil, nil)
	assert.Nil(s.T(), err1)

	value = s.sut.IsHeartbeatWithinDuration()
	assert.False(s.T(), value)

	timestamp = time.Now()
	data.Timestamp = model.NewAbsoluteOrRelativeTimeTypeFromTime(timestamp)

	_, err1 = remoteDiagServer.UpdateData(true, model.FunctionTypeDeviceDiagnosisHeartbeatData, data, nil, nil)
	assert.Nil(s.T(), err1)

	value = s.sut.IsHeartbeatWithinDuration()
	assert.True(s.T(), value)

	s.sut.StopHeartbeat()
	s.sut.StartHeartbeat()
}

func (s *CsLPPSuite) Test_ContractualProductionNominalMax() {
	value, err := s.sut.ProductionNominalMax()
	assert.Equal(s.T(), 0.0, value)
	assert.NotNil(s.T(), err)

	err = s.sut.SetProductionNominalMax(10)
	assert.Nil(s.T(), err)

	value, err = s.sut.ProductionNominalMax()
	assert.Equal(s.T(), 10.0, value)
	assert.Nil(s.T(), err)
}
