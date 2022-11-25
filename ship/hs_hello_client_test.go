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

func TestHelloClientSuite(t *testing.T) {
	suite.Run(t, new(HelloClientSuite))
}

// Hello Client role specific tests
type HelloClientSuite struct {
	suite.Suite

	sut *ShipConnection

	sentMessage []byte

	mux sync.Mutex
}

func (s *HelloClientSuite) lastMessage() []byte {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.sentMessage
}

func (s *HelloClientSuite) setSentMessage(value []byte) {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.sentMessage = value
}

var _ ConnectionHandler = (*HelloClientSuite)(nil)

func (s *HelloClientSuite) HandleClosedConnection(connection *ShipConnection) {}

var _ ShipServiceDataProvider = (*HelloClientSuite)(nil)

func (s *HelloClientSuite) IsRemoteServiceForSKIPaired(string) bool { return true }

var _ ShipDataConnection = (*HelloClientSuite)(nil)

func (s *HelloClientSuite) InitDataProcessing(dataProcessing ShipDataProcessing) {}

func (s *HelloClientSuite) WriteMessageToDataConnection(message []byte) error {
	s.setSentMessage(message)

	return nil
}

func (s *HelloClientSuite) CloseDataConnection() {}

func (s *HelloClientSuite) SetupSuite()   {}
func (s *HelloClientSuite) TearDownTest() {}

func (s *HelloClientSuite) BeforeTest(suiteName, testName string) {
	s.setSentMessage(nil)

	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", spineModel.DeviceTypeTypeEnergyManagementSystem, spineModel.NetworkManagementFeatureSetTypeSmart)

	s.sut = NewConnectionHandler(s, s, localDevice, ShipRoleClient, "LocalShipID", "RemoveDevice", "RemoteShipID")

	s.sut.handshakeTimer = time.NewTimer(time.Hour * 1)
	s.sut.stopHandshakeTimer()
}

func (s *HelloClientSuite) AfterTest(suiteName, testName string) {
	s.sut.stopHandshakeTimer()
}

func (s *HelloClientSuite) Test_InitialState() {
	s.sut.setState(smeHelloState)
	s.sut.handleState(false, nil)

	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeHelloStateReadyListen, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *HelloClientSuite) Test_ReadyListen_Ok() {
	s.sut.setState(smeHelloStateReadyInit) // inits the timer
	s.sut.setState(smeHelloStateReadyListen)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypeReady,
		},
	}

	msg, err := s.sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	s.sut.handleState(false, msg)

	// the state goes from smeHelloStateOk directly to smeProtHStateClientInit to smeProtHStateClientListenChoice
	assert.Equal(s.T(), smeProtHStateClientListenChoice, s.sut.getState())
}
