package ship

import (
	"testing"

	"github.com/enbility/eebus-go/ship/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestProClientSuite(t *testing.T) {
	suite.Run(t, new(ProClientSuite))
}

type ProClientSuite struct {
	suite.Suite

	role shipRole
}

func (s *ProClientSuite) BeforeTest(suiteName, testName string) {
	s.role = ShipRoleClient
}

func (s *ProClientSuite) Test_Init() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloStateOk, nil)

	sut.handleState(false, nil)

	// the state goes from smeHelloStateOk to smeProtHStateClientInit to smeProtHStateClientListenChoice
	assert.Equal(s.T(), SmeProtHStateClientListenChoice, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *ProClientSuite) Test_ListenChoice() {
	sut, data := initTest(s.role)

	sut.setState(SmeProtHStateClientListenChoice, nil)

	protMsg := model.MessageProtocolHandshake{
		MessageProtocolHandshake: model.MessageProtocolHandshakeType{
			HandshakeType: model.ProtocolHandshakeTypeTypeSelect,
			Version:       model.Version{Major: 1, Minor: 0},
			Formats: model.MessageProtocolFormatsType{
				Format: []model.MessageProtocolFormatType{model.MessageProtocolFormatTypeUTF8},
			},
		},
	}

	msg, err := sut.shipMessage(model.MsgTypeControl, protMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)

	// state goes directly from smeProtHStateClientOk to smePinStateCheckInit to smePinStateCheckListen
	assert.Equal(s.T(), SmePinStateCheckListen, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}
