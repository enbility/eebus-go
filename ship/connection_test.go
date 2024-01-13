package ship

import (
	"encoding/json"
	"testing"

	"github.com/enbility/eebus-go/ship/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

func TestConnectionSuite(t *testing.T) {
	suite.Run(t, new(ConnectionSuite))
}

type ConnectionSuite struct {
	suite.Suite

	sut *ShipConnectionImpl

	shipDataProvider *MockShipServiceDataProvider
	shipDataConn     *MockWebsocketDataConnection

	spineDataProcessing *MockSpineDataProcessing

	sentMessage []byte
}

func (s *ConnectionSuite) BeforeTest(suiteName, testName string) {
	s.sentMessage = nil

	ctrl := gomock.NewController(s.T())

	s.shipDataProvider = NewMockShipServiceDataProvider(ctrl)
	s.shipDataProvider.EXPECT().HandleShipHandshakeStateUpdate(gomock.Any(), gomock.Any()).AnyTimes()
	s.shipDataProvider.EXPECT().HandleConnectionClosed(gomock.Any(), gomock.Any()).AnyTimes()
	s.shipDataProvider.EXPECT().IsRemoteServiceForSKIPaired(gomock.Any()).AnyTimes()
	s.shipDataProvider.EXPECT().AllowWaitingForTrust(gomock.Any()).AnyTimes()

	s.shipDataConn = NewMockWebsocketDataConnection(ctrl)
	s.shipDataConn.EXPECT().InitDataProcessing(gomock.Any()).AnyTimes()
	s.shipDataConn.EXPECT().WriteMessageToDataConnection(gomock.Any()).DoAndReturn(func(message []byte) error { s.sentMessage = message; return nil }).AnyTimes()
	s.shipDataConn.EXPECT().IsDataConnectionClosed().DoAndReturn(func() (bool, error) { return false, nil }).AnyTimes()
	s.shipDataConn.EXPECT().CloseDataConnection(gomock.Any(), gomock.Any()).AnyTimes()

	s.spineDataProcessing = NewMockSpineDataProcessing(ctrl)
	s.spineDataProcessing.EXPECT().HandleIncomingSpineMesssage(gomock.Any()).AnyTimes()

	s.sut = NewConnectionHandler(s.shipDataProvider, s.shipDataConn, ShipRoleServer, "LocalShipID", "RemoveDevice", "RemoteShipID")
}

func (s *ConnectionSuite) Test_RemoteSKI() {
	remoteSki := s.sut.RemoteSKI()
	assert.NotEqual(s.T(), "", remoteSki)
}

func (s *ConnectionSuite) Test_DataHandler() {
	handler := s.sut.DataHandler()
	assert.NotNil(s.T(), handler)
}

func (s *ConnectionSuite) TestRun() {
	s.sut.Run()
	state, err := s.sut.ShipHandshakeState()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), CmiStateServerWait, state)
}

func (s *ConnectionSuite) TestShipHandshakeState() {
	state, err := s.sut.ShipHandshakeState()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), CmiStateInitStart, state)
}

func (s *ConnectionSuite) TestApprovePendingHandshake() {
	s.sut.smeState = CmiStateInitStart
	s.sut.ApprovePendingHandshake()
	assert.Equal(s.T(), CmiStateInitStart, s.sut.smeState)

	s.sut.smeState = SmeHelloStatePendingListen
	s.sut.ApprovePendingHandshake()
	assert.Equal(s.T(), SmeProtHStateServerListenProposal, s.sut.smeState)
}

func (s *ConnectionSuite) TestAbortPendingHandshake() {
	s.sut.smeState = CmiStateInitStart
	s.sut.AbortPendingHandshake()
	assert.Equal(s.T(), CmiStateInitStart, s.sut.smeState)

	s.sut.smeState = SmeHelloStatePendingListen
	s.sut.AbortPendingHandshake()
	assert.Equal(s.T(), SmeHelloStateAbortDone, s.sut.smeState)
}

func (s *ConnectionSuite) TestCloseConnection_StateComplete() {
	s.sut.smeState = SmeStateComplete
	s.sut.CloseConnection(true, 450, "User Close")
	state, err := s.sut.ShipHandshakeState()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), SmeStateComplete, state)
}

func (s *ConnectionSuite) TestCloseConnection_StateComplete_2() {
	s.sut.smeState = SmeStateError
	s.sut.CloseConnection(false, 0, "User Close")
	state, err := s.sut.ShipHandshakeState()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), SmeStateError, state)
}

func (s *ConnectionSuite) TestCloseConnection_StateComplete_3() {
	s.sut.smeState = SmeStateError
	s.sut.CloseConnection(false, 450, "User Close")
	state, err := s.sut.ShipHandshakeState()
	assert.Nil(s.T(), err)
	assert.Equal(s.T(), SmeStateError, state)
}

func (s *ConnectionSuite) TestShipModelFromMessage() {
	msg := []byte{}
	data, err := s.sut.shipModelFromMessage(msg)
	assert.NotNil(s.T(), err)
	assert.Nil(s.T(), data)

	modelData := model.ShipData{}
	jsonData, err := json.Marshal(modelData)
	assert.Nil(s.T(), err)

	msg = []byte{0}
	msg = append(msg, jsonData...)
	data, err = s.sut.shipModelFromMessage(msg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), data)
}

func (s *ConnectionSuite) TestHandleIncomingShipMessage() {
	modelData := model.ShipData{}
	jsonData, err := json.Marshal(modelData)
	assert.Nil(s.T(), err)

	msg := []byte{0}
	msg = append(msg, jsonData...)

	s.sut.HandleIncomingShipMessage(msg)

	spineData := `{"datagram":{}}`
	jsonData = []byte(spineData)

	modelData = model.ShipData{
		Data: model.DataType{
			Payload: jsonData,
		},
	}
	jsonData, err = json.Marshal(modelData)
	assert.Nil(s.T(), err)

	msg = []byte{0}
	msg = append(msg, jsonData...)

	s.sut.HandleIncomingShipMessage(msg)

	s.sut.spineDataProcessing = s.spineDataProcessing

	s.sut.HandleIncomingShipMessage(msg)
}

func (s *ConnectionSuite) TestReportConnectionError() {
	s.sut.ReportConnectionError(nil)
	assert.Equal(s.T(), SmeStateError, s.sut.smeState)

	s.sut.smeState = SmeHelloStateReadyListen
	s.sut.ReportConnectionError(nil)
	assert.Equal(s.T(), SmeHelloStateRejected, s.sut.smeState)

	s.sut.smeState = SmeHelloStateRemoteAbortDone
	s.sut.ReportConnectionError(nil)
	assert.Equal(s.T(), SmeHelloStateRemoteAbortDone, s.sut.smeState)

	s.sut.smeState = SmeHelloStateAbort
	s.sut.ReportConnectionError(nil)
	assert.Equal(s.T(), SmeHelloStateAbort, s.sut.smeState)
}

func (s *ConnectionSuite) TestSendShipModel() {
	err := s.sut.sendShipModel(model.MsgTypeInit, nil)
	assert.NotNil(s.T(), err)

	closeMessage := model.ConnectionClose{
		ConnectionClose: model.ConnectionCloseType{
			Phase: model.ConnectionClosePhaseTypeAnnounce,
		},
	}

	err = s.sut.sendShipModel(model.MsgTypeControl, closeMessage)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), s.sentMessage)
}

func (s *ConnectionSuite) TestProcessShipJsonMessage() {
	closeMessage := model.ConnectionClose{
		ConnectionClose: model.ConnectionCloseType{
			Phase: model.ConnectionClosePhaseTypeAnnounce,
		},
	}
	msg, err := json.Marshal(closeMessage)
	assert.Nil(s.T(), err)

	newMsg := []byte{model.MsgTypeControl}
	newMsg = append(newMsg, msg...)

	var data any
	err = s.sut.processShipJsonMessage(newMsg, &data)
	assert.Nil(s.T(), err)
}

func (s *ConnectionSuite) TestSendSpineMessage() {
	data := `{"datagram":{"header":{},"payload":{"cmd":[]}}}`

	err := s.sut.sendSpineData([]byte(data))
	assert.Nil(s.T(), err)
}
