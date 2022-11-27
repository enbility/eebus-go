package ship

import (
	"testing"

	"github.com/DerAndereAndi/eebus-go/ship/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestInitClientSuite(t *testing.T) {
	suite.Run(t, new(InitClientSuite))
}

type InitClientSuite struct {
	suite.Suite
	role shipRole
}

func (s *InitClientSuite) BeforeTest(suiteName, testName string) {
	s.role = ShipRoleClient
}

func (s *InitClientSuite) Test_Init() {
	sut, _ := initTest(s.role)

	assert.Equal(s.T(), cmiStateInitStart, sut.getState())

	shutdownTest(sut)
}

func (s *InitClientSuite) Test_Start() {
	sut, data := initTest(s.role)

	sut.setState(cmiStateInitStart)

	sut.handleState(false, nil)

	assert.Equal(s.T(), true, sut.handshakeTimerRunning)
	assert.Equal(s.T(), cmiStateClientWait, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())
	assert.Equal(s.T(), shipInit, data.lastMessage())

	shutdownTest(sut)
}

func (s *InitClientSuite) Test_ClientWait() {
	sut, data := initTest(s.role)

	sut.setState(cmiStateClientWait)

	sut.handleState(false, shipInit)

	// the state goes from smeHelloState directly to smeHelloStateReadyInit to smeHelloStateReadyListen
	assert.Equal(s.T(), smeHelloStateReadyListen, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *InitClientSuite) Test_ClientWait_InvalidMsgType() {
	sut, data := initTest(s.role)

	sut.setState(cmiStateClientWait)

	sut.handleState(false, []byte{0x05, 0x00})

	assert.Equal(s.T(), smeError, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *InitClientSuite) Test_ClientWait_InvalidData() {
	sut, data := initTest(s.role)

	sut.setState(cmiStateClientWait)

	sut.handleState(false, []byte{model.MsgTypeInit, 0x05})

	assert.Equal(s.T(), smeError, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}
