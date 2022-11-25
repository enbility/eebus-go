package ship

import (
	"os"
	"sync"
	"testing"
	"time"

	"github.com/DerAndereAndi/eebus-go/ship/model"
	"github.com/DerAndereAndi/eebus-go/spine"
	spineModel "github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func skipCI(t *testing.T) {
	if os.Getenv("ACTION_ENVIRONMENT") == "CI" {
		t.Skip("Skipping testing in CI environment")
	}
}

func TestHelloSuite(t *testing.T) {
	suite.Run(t, new(HelloSuite))
}

type HelloSuite struct {
	suite.Suite

	sut *ShipConnection

	sentMessage []byte

	mux sync.Mutex
}

func (s *HelloSuite) lastMessage() []byte {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.sentMessage
}

var _ ConnectionHandler = (*HelloSuite)(nil)

func (s *HelloSuite) HandleClosedConnection(connection *ShipConnection) {}

var _ ShipServiceDataProvider = (*HelloSuite)(nil)

func (s *HelloSuite) IsRemoteServiceForSKIPaired(string) bool { return true }

var _ ShipDataConnection = (*HelloSuite)(nil)

func (s *HelloSuite) InitDataProcessing(dataProcessing ShipDataProcessing) {}

func (s *HelloSuite) WriteMessageToDataConnection(message []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.sentMessage = message
	return nil
}

func (s *HelloSuite) CloseDataConnection() {}

func (s *HelloSuite) SetupSuite()   {}
func (s *HelloSuite) TearDownTest() {}

func (s *HelloSuite) BeforeTest(suiteName, testName string) {
	s.sentMessage = nil

	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", spineModel.DeviceTypeTypeEnergyManagementSystem, spineModel.NetworkManagementFeatureSetTypeSmart)

	s.sut = NewConnectionHandler(s, s, localDevice, ShipRoleServer, "LocalShipID", "RemoveDevice", "RemoteShipID")

	s.sut.handshakeTimer = time.NewTimer(time.Hour * 1)
	s.sut.stopHandshakeTimer()
}

func (s *HelloSuite) AfterTest(suiteName, testName string) {
	s.sut.stopHandshakeTimer()
}

func (s *HelloSuite) Test_InitialState() {
	s.sut.setState(smeHelloState)
	s.sut.handleState(false, nil)

	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeHelloStateReadyListen, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *HelloSuite) Test_ReadyListen_Init() {
	s.sut.setState(smeHelloStateReadyInit)
	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)
}

func (s *HelloSuite) Test_ReadyListen_Ok() {
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

	// the state goes from smeHelloStateOk directly to smeProtHStateServerInit to smeProtHStateClientListenProposal
	assert.Equal(s.T(), smeProtHStateServerListenProposal, s.sut.getState())
}

func (s *HelloSuite) Test_ReadyListen_Timeout() {
	skipCI(s.T())

	s.sut.setState(smeHelloStateReadyInit) // inits the timer
	s.sut.setState(smeHelloStateReadyListen)

	time.Sleep(tHelloInit + time.Second)

	assert.Equal(s.T(), smeHelloStateAbort, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *HelloSuite) Test_ReadyListen_Ignore() {
	s.sut.setState(smeHelloStateReadyInit) // inits the timer
	s.sut.setState(smeHelloStateReadyListen)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypePending,
		},
	}

	msg, err := s.sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	s.sut.handleState(false, msg)

	assert.Equal(s.T(), smeHelloStateReadyListen, s.sut.getState())
}

func (s *HelloSuite) Test_ReadyListen_Abort() {
	s.sut.setState(smeHelloStateReadyInit) // inits the timer
	s.sut.setState(smeHelloStateReadyListen)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypeAborted,
		},
	}

	msg, err := s.sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	s.sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), false, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeHelloStateAbort, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *HelloSuite) Test_PendingInit() {
	s.sut.setState(smeHelloStatePendingInit)
	s.sut.handleState(false, nil)

	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeHelloStatePendingListen, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *HelloSuite) Test_PendingListen() {
	s.sut.setState(smeHelloStatePendingInit) // inits the timer
	s.sut.setState(smeHelloStatePendingListen)
	s.sut.handleState(false, nil)
}

func (s *HelloSuite) Test_PendingListen_Timeout() {
	skipCI(s.T())

	s.sut.setState(smeHelloStatePendingInit) // inits the timer
	s.sut.setState(smeHelloStatePendingListen)

	time.Sleep(tHelloInit + time.Second)

	assert.Equal(s.T(), smeHelloStateAbort, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *HelloSuite) Test_PendingListen_ReadyAbort() {
	s.sut.setState(smeHelloStatePendingInit) // inits the timer
	s.sut.setState(smeHelloStatePendingListen)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypeReady,
		},
	}

	msg, err := s.sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	s.sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), false, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeHelloStateAbort, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *HelloSuite) Test_PendingListen_ReadyWaiting() {
	s.sut.setState(smeHelloStatePendingInit) // inits the timer
	s.sut.setState(smeHelloStatePendingListen)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase:   model.ConnectionHelloPhaseTypeReady,
			Waiting: util.Ptr(uint(tHelloInit.Milliseconds())),
		},
	}

	msg, err := s.sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	s.sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeHelloStatePendingListen, s.sut.getState())
}

func (s *HelloSuite) Test_PendingListen_Abort() {
	s.sut.setState(smeHelloStatePendingInit) // inits the timer
	s.sut.setState(smeHelloStatePendingListen)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase: model.ConnectionHelloPhaseTypeAborted,
		},
	}

	msg, err := s.sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	s.sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), false, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeHelloStateAbort, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}

func (s *HelloSuite) Test_PendingListen_PendingWaiting() {
	s.sut.setState(smeHelloStatePendingInit) // inits the timer
	s.sut.setState(smeHelloStatePendingListen)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase:   model.ConnectionHelloPhaseTypePending,
			Waiting: util.Ptr(uint(tHelloInit.Milliseconds())),
		},
	}

	msg, err := s.sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	s.sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeHelloStatePendingListen, s.sut.getState())
}

func (s *HelloSuite) Test_PendingListen_PendingProlongation() {
	s.sut.setState(smeHelloStatePendingInit) // inits the timer
	s.sut.setState(smeHelloStatePendingListen)

	helloMsg := model.ConnectionHello{
		ConnectionHello: model.ConnectionHelloType{
			Phase:               model.ConnectionHelloPhaseTypePending,
			ProlongationRequest: util.Ptr(true),
		},
	}

	msg, err := s.sut.shipMessage(model.MsgTypeControl, helloMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	s.sut.handleShipMessage(false, msg)

	assert.Equal(s.T(), true, s.sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeHelloStatePendingListen, s.sut.getState())
	assert.NotNil(s.T(), s.lastMessage())
}
