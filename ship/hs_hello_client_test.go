package ship

import (
	"testing"

	"github.com/enbility/eebus-go/ship/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestHelloClientSuite(t *testing.T) {
	suite.Run(t, new(HelloClientSuite))
}

// Hello Client role specific tests
type HelloClientSuite struct {
	suite.Suite
	role shipRole
}

func (s *HelloClientSuite) BeforeTest(suiteName, testName string) {
	s.role = ShipRoleClient
}

func (s *HelloClientSuite) Test_InitialState() {
	sut, data := initTest(s.role)

	sut.setState(smeHelloState)
	sut.handleState(false, nil)

	assert.Equal(s.T(), true, sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeHelloStateReadyListen, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *HelloClientSuite) Test_ReadyListen_Ok() {
	sut, _ := initTest(s.role)

	sut.setState(smeHelloStateReadyInit) // inits the timer
	sut.setState(smeHelloStateReadyListen)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypeReady,
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	// the state goes from smeHelloStateOk directly to smeProtHStateClientInit to smeProtHStateClientListenChoice
	assert.Equal(s.T(), smeProtHStateClientListenChoice, sut.getState())

	shutdownTest(sut)
}
