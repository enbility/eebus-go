package ship

import (
	"testing"
	"time"

	"github.com/DerAndereAndi/eebus-go/ship/model"
	"github.com/DerAndereAndi/eebus-go/spine"
	spineModel "github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestProClientSuite(t *testing.T) {
	suite.Run(t, new(ProClientSuite))
}

type ProClientSuite struct {
	suite.Suite

	sut *ShipConnection

	sentMessage []byte
}

var _ ConnectionHandler = (*ProClientSuite)(nil)

func (s *ProClientSuite) HandleClosedConnection(connection *ShipConnection) {}

var _ ShipServiceDataProvider = (*ProClientSuite)(nil)

func (s *ProClientSuite) IsRemoteServiceForSKIPaired(string) bool { return true }

var _ ShipDataConnection = (*ProClientSuite)(nil)

func (s *ProClientSuite) InitDataProcessing(dataProcessing ShipDataProcessing) {}

func (s *ProClientSuite) WriteMessageToDataConnection(message []byte) error {
	s.sentMessage = message
	return nil
}

func (s *ProClientSuite) CloseDataConnection() {}

func (s *ProClientSuite) SetupSuite()   {}
func (s *ProClientSuite) TearDownTest() {}

func (s *ProClientSuite) BeforeTest(suiteName, testName string) {
	s.sentMessage = nil

	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", spineModel.DeviceTypeTypeEnergyManagementSystem, spineModel.NetworkManagementFeatureSetTypeSmart)

	s.sut = NewConnectionHandler(s, s, localDevice, ShipRoleClient, "LocalShipID", "RemoveDevice", "RemoteShipID")

	s.sut.handshakeTimer = time.NewTimer(time.Hour * 1)
	s.sut.stopHandshakeTimer()
}

func (s *ProClientSuite) AfterTest(suiteName, testName string) {
	s.sut.stopHandshakeTimer()
}

func (s *ProClientSuite) Test_Init() {
	s.sut.setState(smeHelloStateOk)

	s.sut.handleState(false, nil)

	// the state goes from smeHelloStateOk to smeProtHStateClientInit to smeProtHStateClientListenChoice
	assert.Equal(s.T(), smeProtHStateClientListenChoice, s.sut.smeState)
	assert.NotNil(s.T(), s.sentMessage)
}

func (s *ProClientSuite) Test_ListenChoice() {
	s.sut.setState(smeProtHStateClientListenChoice)

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

	// state goes directly from smeProtHStateClientOk to smePinStateCheckInit to smePinStateCheckListen
	assert.Equal(s.T(), smePinStateCheckListen, s.sut.smeState)
	assert.NotNil(s.T(), s.sentMessage)
}
