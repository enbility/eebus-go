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

func TestInitClientSuite(t *testing.T) {
	suite.Run(t, new(InitClientSuite))
}

type InitClientSuite struct {
	suite.Suite

	sut *ShipConnection

	sentMessage []byte

	mux sync.Mutex
}

func (s *InitClientSuite) lastMessage() []byte {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.sentMessage
}

func (s *InitClientSuite) setSentMessage(value []byte) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.sentMessage = value
}

var _ ConnectionHandler = (*InitClientSuite)(nil)

func (s *InitClientSuite) HandleClosedConnection(connection *ShipConnection) {}

var _ ShipServiceDataProvider = (*InitClientSuite)(nil)

func (s *InitClientSuite) IsRemoteServiceForSKIPaired(string) bool { return true }

var _ ShipDataConnection = (*InitClientSuite)(nil)

func (s *InitClientSuite) InitDataProcessing(dataProcessing ShipDataProcessing) {}

func (s *InitClientSuite) WriteMessageToDataConnection(message []byte) error {
	s.setSentMessage(message)

	return nil
}

func (s *InitClientSuite) CloseDataConnection() {}

func (s *InitClientSuite) SetupSuite()   {}
func (s *InitClientSuite) TearDownTest() {}

func (s *InitClientSuite) BeforeTest(suiteName, testName string) {
	s.setSentMessage(nil)

	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", spineModel.DeviceTypeTypeEnergyManagementSystem, spineModel.NetworkManagementFeatureSetTypeSmart)

	s.sut = NewConnectionHandler(s, s, localDevice, ShipRoleClient, "LocalShipID", "RemoveDevice", "RemoteShipID")

	s.sut.handshakeTimer = time.NewTimer(time.Hour * 1)
	s.sut.stopHandshakeTimer()
}

func (s *InitClientSuite) AfterTest(suiteName, testName string) {
	s.sut.stopHandshakeTimer()
}

func (s *InitClientSuite) Test_Init() {
	assert.Equal(s.T(), cmiStateInitStart, s.sut.getState())
}

func (s *InitClientSuite) Test_Start() {
	s.sut.setState(cmiStateInitStart)

	s.sut.handleState(false, nil)

	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), cmiStateClientWait, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
	assert.Equal(s.T(), shipInit, s.lastMessage())
}

func (s *InitClientSuite) Test_ClientWait() {
	s.sut.setState(cmiStateClientWait)

	s.sut.handleState(false, shipInit)

	// the state goes from smeHelloState directly to smeHelloStateReadyInit to smeHelloStateReadyListen
	assert.Equal(s.T(), smeHelloStateReadyListen, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *InitClientSuite) Test_ClientWait_InvalidMsgType() {
	s.sut.setState(cmiStateClientWait)

	s.sut.handleState(false, []byte{0x05, 0x00})

	assert.Equal(s.T(), smeError, s.sut.getState())
	assert.Nil(s.T(), s.lastMessage())
}

func (s *InitClientSuite) Test_ClientWait_InvalidData() {
	s.sut.setState(cmiStateClientWait)

	s.sut.handleState(false, []byte{model.MsgTypeInit, 0x05})

	assert.Equal(s.T(), smeError, s.sut.getState())
	assert.Nil(s.T(), s.lastMessage())
}
