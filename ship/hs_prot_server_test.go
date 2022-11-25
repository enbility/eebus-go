package ship

import (
	"sync"
	"testing"
	"time"

	"github.com/DerAndereAndi/eebus-go/ship/model"
	"github.com/DerAndereAndi/eebus-go/spine"
	spineModel "github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestProServerSuite(t *testing.T) {
	suite.Run(t, new(ProServerSuite))
}

type ProServerSuite struct {
	suite.Suite

	sut *ShipConnection

	sentMessage []byte

	mux sync.Mutex
}

func (s *ProServerSuite) lastMessage() []byte {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.sentMessage
}

func (s *ProServerSuite) setSentMessage(value []byte) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.sentMessage = value
}

var _ ConnectionHandler = (*ProServerSuite)(nil)

func (s *ProServerSuite) HandleClosedConnection(connection *ShipConnection) {}

var _ ShipServiceDataProvider = (*ProServerSuite)(nil)

func (s *ProServerSuite) IsRemoteServiceForSKIPaired(string) bool { return true }

var _ ShipDataConnection = (*ProServerSuite)(nil)

func (s *ProServerSuite) InitDataProcessing(dataProcessing ShipDataProcessing) {}

func (s *ProServerSuite) WriteMessageToDataConnection(message []byte) error {
	s.setSentMessage(message)

	return nil
}

func (s *ProServerSuite) CloseDataConnection() {}

func (s *ProServerSuite) SetupSuite()   {}
func (s *ProServerSuite) TearDownTest() {}

func (s *ProServerSuite) BeforeTest(suiteName, testName string) {
	s.setSentMessage(nil)

	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", spineModel.DeviceTypeTypeEnergyManagementSystem, spineModel.NetworkManagementFeatureSetTypeSmart)

	s.sut = NewConnectionHandler(s, s, localDevice, ShipRoleServer, "LocalShipID", "RemoveDevice", "RemoteShipID")

	s.sut.handshakeTimer = time.NewTimer(time.Hour * 1)
	s.sut.stopHandshakeTimer()
}

func (s *ProServerSuite) AfterTest(suiteName, testName string) {
	s.sut.stopHandshakeTimer()
}

func (s *ProServerSuite) Test_Init() {
	s.sut.setState(smeHelloStateOk)

	s.sut.handleState(false, nil)

	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)

	// the state goes from smeHelloStateOk to smeProtHStateServerInit to smeProtHStateServerListenProposal
	assert.Equal(s.T(), smeProtHStateServerListenProposal, s.sut.getState())
	assert.Nil(s.T(), s.lastMessage())
}

func (s *ProServerSuite) Test_ListenProposal() {
	s.sut.setState(smeProtHStateServerListenProposal)

	protMsg := model.MessageProtocolHandshake{
		MessageProtocolHandshake: model.MessageProtocolHandshakeType{
			HandshakeType: model.ProtocolHandshakeTypeTypeAnnounceMax,
			Version:       model.Version{Major: 1, Minor: 0},
			Formats: model.MessageProtocolFormatsType{
				Format: []model.MessageProtocolFormatType{model.MessageProtocolFormatTypeUTF8},
			},
		},
	}

	msg, err := s.sut.shipMessage(model.MsgTypeControl, protMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	s.sut.handleState(false, msg)

	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)

	assert.Equal(s.T(), smeProtHStateServerListenConfirm, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *ProServerSuite) Test_ListenConfirm() {
	s.sut.setState(smeProtHStateServerListenConfirm)

	protMsg := model.MessageProtocolHandshake{
		MessageProtocolHandshake: model.MessageProtocolHandshakeType{
			HandshakeType: model.ProtocolHandshakeTypeTypeSelect,
			Version:       model.Version{Major: 1, Minor: 0},
			Formats: model.MessageProtocolFormatsType{
				Format: []model.MessageProtocolFormatType{model.MessageProtocolFormatTypeUTF8},
			},
		},
	}

	msg, err := s.sut.shipMessage(model.MsgTypeControl, protMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	s.sut.handleState(false, msg)

	assert.Equal(s.T(), false, s.sut.handshakeTimerRunning)

	// state smeProtHStateServerOk directly goes to smePinStateCheckInit to smePinStateCheckListen
	assert.Equal(s.T(), smePinStateCheckListen, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}
