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

func TestInitServerSuite(t *testing.T) {
	suite.Run(t, new(InitServerSuite))
}

type InitServerSuite struct {
	suite.Suite

	sut *ShipConnection

	sentMessage []byte

	mux sync.Mutex
}

func (s *InitServerSuite) lastMessage() []byte {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.sentMessage
}

func (s *InitServerSuite) setSentMessage(value []byte) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.sentMessage = value
}

var _ ConnectionHandler = (*InitServerSuite)(nil)

func (s *InitServerSuite) HandleClosedConnection(connection *ShipConnection) {}

var _ ShipServiceDataProvider = (*InitServerSuite)(nil)

func (s *InitServerSuite) IsRemoteServiceForSKIPaired(string) bool { return true }

var _ ShipDataConnection = (*InitServerSuite)(nil)

func (s *InitServerSuite) InitDataProcessing(dataProcessing ShipDataProcessing) {}

func (s *InitServerSuite) WriteMessageToDataConnection(message []byte) error {
	s.setSentMessage(message)

	return nil
}

func (s *InitServerSuite) CloseDataConnection() {}

func (s *InitServerSuite) SetupSuite()   {}
func (s *InitServerSuite) TearDownTest() {}

func (s *InitServerSuite) BeforeTest(suiteName, testName string) {
	s.setSentMessage(nil)

	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", spineModel.DeviceTypeTypeEnergyManagementSystem, spineModel.NetworkManagementFeatureSetTypeSmart)

	s.sut = NewConnectionHandler(s, s, localDevice, ShipRoleServer, "LocalShipID", "RemoveDevice", "RemoteShipID")

	s.sut.handshakeTimer = time.NewTimer(time.Hour * 1)
	s.sut.stopHandshakeTimer()
}

func (s *InitServerSuite) AfterTest(suiteName, testName string) {
	s.sut.stopHandshakeTimer()
}

func (s *InitServerSuite) Test_Init() {
	assert.Equal(s.T(), cmiStateInitStart, s.sut.getState())
}

func (s *InitServerSuite) Test_Start() {
	s.sut.setState(cmiStateInitStart)

	s.sut.handleState(false, nil)

	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), cmiStateServerWait, s.sut.getState())
}

func (s *InitServerSuite) Test_ServerWait() {
	s.sut.setState(cmiStateServerWait)

	s.sut.handleState(false, shipInit)

	// the state goes from smeHelloState directly to smeHelloStateReadyInit to smeHelloStateReadyListen
	assert.Equal(s.T(), smeHelloStateReadyListen, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *InitServerSuite) Test_ServerWait_InvalidMsgType() {
	s.sut.setState(cmiStateServerWait)

	s.sut.handleState(false, []byte{0x05, 0x00})

	assert.Equal(s.T(), smeError, s.sut.getState())
	assert.Nil(s.T(), s.lastMessage())
}

func (s *InitServerSuite) Test_ServerWait_InvalidData() {
	s.sut.setState(cmiStateServerWait)

	s.sut.handleState(false, []byte{model.MsgTypeInit, 0x05})

	assert.Equal(s.T(), smeError, s.sut.getState())
	assert.Nil(s.T(), s.lastMessage())
}
