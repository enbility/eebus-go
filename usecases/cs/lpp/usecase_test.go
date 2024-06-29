package lpp

import (
	"time"

	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *CsLPPSuite) Test_loadControlServerAndLimitId() {
	lc, _, err := s.sut.loadControlServerAndLimitId()
	assert.NotNil(s.T(), lc)
	assert.Nil(s.T(), err)

	f := s.sut.LocalEntity.FeatureOfTypeAndRole(model.FeatureTypeTypeLoadControl, model.RoleTypeServer)
	f.UpdateData(model.FunctionTypeLoadControlLimitDescriptionListData, &model.LoadControlLimitDescriptionListDataType{}, nil, nil)
	lc, _, err = s.sut.loadControlServerAndLimitId()
	assert.NotNil(s.T(), lc)
	assert.NotNil(s.T(), err)

	s.sut.LocalEntity = nil
	lc, _, err = s.sut.loadControlServerAndLimitId()
	assert.Nil(s.T(), lc)
	assert.NotNil(s.T(), err)
}

func (s *CsLPPSuite) Test_loadControlWriteCB() {
	msg := &spineapi.Message{}

	s.sut.loadControlWriteCB(msg)

	msg = &spineapi.Message{
		RequestHeader: &model.HeaderType{
			MsgCounter: util.Ptr(model.MsgCounterType(500)),
		},
		Cmd: model.CmdType{
			LoadControlLimitListData: &model.LoadControlLimitListDataType{},
		},
		DeviceRemote: s.remoteDevice,
		EntityRemote: s.monitoredEntity,
	}

	s.sut.loadControlWriteCB(msg)

	msg.Cmd = model.CmdType{
		LoadControlLimitListData: &model.LoadControlLimitListDataType{
			LoadControlLimitData: []model.LoadControlLimitDataType{},
		},
	}

	s.sut.loadControlWriteCB(msg)

	msg.Cmd = model.CmdType{
		LoadControlLimitListData: &model.LoadControlLimitListDataType{
			LoadControlLimitData: []model.LoadControlLimitDataType{
				{},
			},
		},
	}

	s.sut.loadControlWriteCB(msg)

	msg.Cmd = model.CmdType{
		LoadControlLimitListData: &model.LoadControlLimitListDataType{
			LoadControlLimitData: []model.LoadControlLimitDataType{
				{
					LimitId:       util.Ptr(model.LoadControlLimitIdType(0)),
					IsLimitActive: util.Ptr(true),
					Value:         model.NewScaledNumberType(1000),
					TimePeriod:    model.NewTimePeriodTypeWithRelativeEndTime(time.Minute * 2),
				},
			},
		},
	}

	s.sut.loadControlWriteCB(msg)
}

func (s *CsLPPSuite) Test_UpdateUseCaseAvailability() {
	s.sut.UpdateUseCaseAvailability(true)
}
