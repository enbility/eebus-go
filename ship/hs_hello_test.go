package ship

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/ship/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHelloSuite(t *testing.T) {
	suite.Run(t, new(HelloSuite))
}

type HelloSuite struct {
	suite.Suite
	role shipRole
}

func (s *HelloSuite) BeforeTest(suiteName, testName string) {
	s.role = ShipRoleServer
}

func (s *HelloSuite) Test_InitialState() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloState, nil)
	sut.handleState(false, nil)

	assert.Equal(s.T(), true, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeHelloStateReadyListen, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_ReadyListen_Init() {
	sut, _ := initTest(s.role)

	sut.setState(SmeHelloStateReadyInit, nil)
	assert.Equal(s.T(), true, sut.handshakeTimerRunning)

	shutdownTest(sut)
}

func (s *HelloSuite) Test_ReadyListen_Ok() {
	sut, _ := initTest(s.role)

	sut.setState(SmeHelloStateReadyInit, nil) // inits the timer
	sut.setState(SmeHelloStateReadyListen, nil)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypeReady,
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	// the state goes from smeHelloStateOk directly to smeProtHStateServerInit to smeProtHStateClientListenProposal
	assert.Equal(s.T(), SmeProtHStateServerListenProposal, sut.getState())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_ReadyListen_Timeout() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloStateReadyInit, nil) // inits the timer
	sut.setState(SmeHelloStateReadyListen, nil)

	if !isRunningOnCI() {
		// test if the function is triggered correctly via the timer
		time.Sleep(tHelloInit + time.Second)
	} else {
		// speed up the test by running the method directly
		sut.handshakeHello_ReadyListen(true, nil)
	}

	assert.Equal(s.T(), SmeHelloStateAbortDone, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_ReadyListen_Ignore() {
	sut, _ := initTest(s.role)

	sut.setState(SmeHelloStateReadyInit, nil) // inits the timer
	sut.setState(SmeHelloStateReadyListen, nil)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypePending,
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	assert.Equal(s.T(), SmeHelloStateReadyListen, sut.getState())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_ReadyListen_Prolongation() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloStateReadyInit, nil) // inits the timer
	sut.setState(SmeHelloStateReadyListen, nil)

	data.allowWaitingForTrust = true

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase:               model.ConnectionHelloPhaseTypePending,
			ProlongationRequest: util.Ptr(true),
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	assert.Equal(s.T(), SmeHelloStateReadyListen, sut.getState())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_ReadyListen_Abort() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloStateReadyInit, nil) // inits the timer
	sut.setState(SmeHelloStateReadyListen, nil)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypeAborted,
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeHelloStateRemoteAbortDone, sut.getState())
	assert.Nil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_PendingInit() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloStatePendingInit, nil)
	sut.handleState(false, nil)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeHelloStateAbortDone, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_PendingListen() {
	sut, _ := initTest(s.role)

	sut.setState(SmeHelloStatePendingInit, nil) // inits the timer
	sut.setState(SmeHelloStatePendingListen, nil)
	sut.handleState(false, nil)

	shutdownTest(sut)
}

func (s *HelloSuite) Test_PendingListen_Timeout() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloStatePendingInit, nil) // inits the timer
	sut.setState(SmeHelloStatePendingListen, nil)

	if !isRunningOnCI() {
		// test if the function is triggered correctly via the timer
		time.Sleep(tHelloInit + time.Second)
	} else {
		// speed up the test by running the method directly
		sut.handshakeHello_PendingListen(true, nil)
	}

	assert.Equal(s.T(), SmeHelloStateAbortDone, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_PendingListen_Timeout_Prolongation() {
	sut, data := initTest(s.role)

	data.allowWaitingForTrust = true

	sut.setState(SmeHelloStatePendingInit, nil) // inits the timer
	sut.setState(SmeHelloStatePendingListen, nil)

	// speed up the test by running the method directly, the timer is already checked
	sut.handshakeHello_PendingListen(true, nil)

	assert.Equal(s.T(), SmeHelloStatePendingListen, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_PendingListen_ReadyAbort() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloStatePendingInit, nil) // inits the timer
	sut.setState(SmeHelloStatePendingListen, nil)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypeReady,
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeHelloStateAbortDone, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_PendingListen_ReadyWaiting() {
	sut, _ := initTest(s.role)

	sut.setState(SmeHelloStatePendingInit, nil) // inits the timer
	sut.setState(SmeHelloStatePendingListen, nil)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase:   model.ConnectionHelloPhaseTypeReady,
			Waiting: util.Ptr(uint(tHelloInit.Milliseconds())),
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), true, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeHelloStatePendingListen, sut.getState())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_PendingListen_Abort() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloStatePendingInit, nil) // inits the timer
	sut.setState(SmeHelloStatePendingListen, nil)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypeAborted,
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeHelloStateRemoteAbortDone, sut.getState())
	assert.Nil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_PendingListen_PendingWaiting() {
	sut, _ := initTest(s.role)

	sut.setState(SmeHelloStatePendingInit, nil) // inits the timer
	sut.setState(SmeHelloStatePendingListen, nil)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase:   model.ConnectionHelloPhaseTypePending,
			Waiting: util.Ptr(uint(tHelloInit.Milliseconds())),
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), true, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeHelloStatePendingListen, sut.getState())

	shutdownTest(sut)
}

func (s *HelloSuite) Test_PendingListen_PendingProlongation() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloStatePendingInit, nil) // inits the timer
	sut.setState(SmeHelloStatePendingListen, nil)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase:               model.ConnectionHelloPhaseTypePending,
			ProlongationRequest: util.Ptr(true),
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), true, sut.handshakeTimerRunning)
	assert.Equal(s.T(), SmeHelloStatePendingListen, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}
