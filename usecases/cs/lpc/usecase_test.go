package lpc

import (
	"time"

	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/model"
	"github.com/enbility/spine-go/util"
	"github.com/stretchr/testify/assert"
)

func (s *LPCSuite) Test_loadControlServerAndLimitId() {
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

func (s *LPCSuite) Test_loadControlWriteCB() {
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

func (s *LPCSuite) Test_UpdateUseCaseAvailability() {
	s.sut.UpdateUseCaseAvailability(true)
}

func (s *LPCSuite) Test_IsUseCaseSupported() {
	data, err := s.sut.IsUseCaseSupported(s.mockRemoteEntity)
	assert.NotNil(s.T(), err)
	assert.Equal(s.T(), false, data)

	data, err = s.sut.IsUseCaseSupported(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), false, data)

	ucData := &model.NodeManagementUseCaseDataType{
		UseCaseInformation: []model.UseCaseInformationDataType{
			{
				Actor: util.Ptr(model.UseCaseActorTypeEnergyGuard),
				UseCaseSupport: []model.UseCaseSupportType{
					{
						UseCaseName:      util.Ptr(model.UseCaseNameTypeLimitationOfPowerConsumption),
						UseCaseAvailable: util.Ptr(true),
						ScenarioSupport:  []model.UseCaseScenarioSupportType{1, 2, 3, 4},
					},
				},
			},
		},
	}

	nodemgmtEntity := s.remoteDevice.Entity([]model.AddressEntityType{0})
	nodeFeature := s.remoteDevice.FeatureByEntityTypeAndRole(nodemgmtEntity, model.FeatureTypeTypeNodeManagement, model.RoleTypeSpecial)
	fErr := nodeFeature.UpdateData(model.FunctionTypeNodeManagementUseCaseData, ucData, nil, nil)
	assert.Nil(s.T(), fErr)

	data, err = s.sut.IsUseCaseSupported(s.monitoredEntity)
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), true, data)
}
