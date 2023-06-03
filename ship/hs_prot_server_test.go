package ship

import (
	"testing"

	"github.com/enbility/eebus-go/ship/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestProServerSuite(t *testing.T) {
	suite.Run(t, new(ProServerSuite))
}

type ProServerSuite struct {
	suite.Suite
	role shipRole
}

func (s *ProServerSuite) BeforeTest(suiteName, testName string) {
	s.role = ShipRoleServer
}

func (s *ProServerSuite) Test_Init() {
	sut, data := initTest(s.role)

	sut.setState(SmeHelloStateOk, nil)

	sut.handleState(false, nil)

	assert.Equal(s.T(), true, sut.handshakeTimerRunning)

	// the state goes from smeHelloStateOk to smeProtHStateServerInit to smeProtHStateServerListenProposal
	assert.Equal(s.T(), SmeProtHStateServerListenProposal, sut.getState())
	assert.Nil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *ProServerSuite) Test_ListenProposal() {
	sut, data := initTest(s.role)

	sut.setState(SmeProtHStateServerListenProposal, nil)

	protMsg := model.MessageProtocolHandshake{
		MessageProtocolHandshake: model.MessageProtocolHandshakeType{
			HandshakeType: model.ProtocolHandshakeTypeTypeAnnounceMax,
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

	assert.Equal(s.T(), true, sut.handshakeTimerRunning)

	assert.Equal(s.T(), SmeProtHStateServerListenConfirm, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *ProServerSuite) Test_ListenConfirm() {
	sut, data := initTest(s.role)

	sut.setState(SmeProtHStateServerListenConfirm, nil)

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

	// state smeProtHStateServerOk directly goes to smePinStateCheckInit to smePinStateCheckListen
	assert.Equal(s.T(), SmePinStateCheckListen, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}
