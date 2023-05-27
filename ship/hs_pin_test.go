package ship

import (
	"testing"

	"github.com/enbility/eebus-go/ship/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestPinSuite(t *testing.T) {
	suite.Run(t, new(PinSuite))
}

type PinSuite struct {
	suite.Suite
}

func (s *PinSuite) Test_Init() {
	sut, data := initTest(ShipRoleClient)

	sut.setState(SmePinStateCheckInit, nil)
	sut.handleState(false, nil)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmePinStateCheckListen, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *PinSuite) Test_CheckListen_None() {
	sut, data := initTest(ShipRoleClient)

	sut.setState(SmePinStateCheckListen, nil)

	pinState := model.ConnectionPinState{
		ConnectionPinState: model.ConnectionPinStateType{
			PinState: model.PinStateTypeNone,
		},
	}
	msg, err := sut.shipMessage(model.MsgTypeControl, pinState)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	assert.Equal(s.T(), true, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeAccessMethodsRequest, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *PinSuite) Test_CheckListen_Required() {
	sut, data := initTest(ShipRoleClient)

	sut.setState(SmePinStateCheckListen, nil)

	pinState := model.ConnectionPinState{
		ConnectionPinState: model.ConnectionPinStateType{
			PinState: model.PinStateTypeRequired,
		},
	}
	msg, err := sut.shipMessage(model.MsgTypeControl, pinState)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeError, sut.getState())
	assert.Nil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *PinSuite) Test_CheckListen_Optional() {
	sut, data := initTest(ShipRoleClient)

	sut.setState(SmePinStateCheckListen, nil)

	pinState := model.ConnectionPinState{
		ConnectionPinState: model.ConnectionPinStateType{
			PinState: model.PinStateTypeOptional,
		},
	}
	msg, err := sut.shipMessage(model.MsgTypeControl, pinState)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeError, sut.getState())
	assert.Nil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *PinSuite) Test_CheckListen_Ok() {
	sut, data := initTest(ShipRoleClient)

	sut.setState(SmePinStateCheckListen, nil)

	pinState := model.ConnectionPinState{
		ConnectionPinState: model.ConnectionPinStateType{
			PinState: model.PinStateTypePinOk,
		},
	}
	msg, err := sut.shipMessage(model.MsgTypeControl, pinState)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeError, sut.getState())
	assert.Nil(s.T(), data.lastMessage())

	shutdownTest(sut)
}
