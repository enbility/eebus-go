package ship

import (
	"encoding/json"
	"testing"

	"github.com/enbility/eebus-go/ship/model"
	"github.com/enbility/eebus-go/spine"
	spineModel "github.com/enbility/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestConnectionSuite(t *testing.T) {
	suite.Run(t, new(ConnectionSuite))
}

type ConnectionSuite struct {
	suite.Suite

	sut *ShipConnection

	sentMessage []byte
}

var _ ConnectionHandler = (*ConnectionSuite)(nil)

func (s *ConnectionSuite) HandleClosedConnection(connection *ShipConnection) {}

var _ ShipServiceDataProvider = (*ConnectionSuite)(nil)

func (s *ConnectionSuite) IsRemoteServiceForSKIPaired(string) bool      { return true }
func (s *ConnectionSuite) HandleConnectionClosed(*ShipConnection, bool) {}
func (s *ConnectionSuite) ReportServiceShipID(string, string)           {}

var _ ShipDataConnection = (*ConnectionSuite)(nil)

func (s *ConnectionSuite) InitDataProcessing(dataProcessing ShipDataProcessing) {}

func (s *ConnectionSuite) WriteMessageToDataConnection(message []byte) error {
	s.sentMessage = message
	return nil
}

func (s *ConnectionSuite) CloseDataConnection()         {}
func (w *ConnectionSuite) IsDataConnectionClosed() bool { return false }

func (s *ConnectionSuite) SetupSuite()   {}
func (s *ConnectionSuite) TearDownTest() {}

func (s *ConnectionSuite) BeforeTest(suiteName, testName string) {
	s.sentMessage = nil
	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", spineModel.DeviceTypeTypeEnergyManagementSystem, spineModel.NetworkManagementFeatureSetTypeSmart)

	s.sut = NewConnectionHandler(s, s, localDevice, ShipRoleServer, "LocalShipID", "RemoveDevice", "RemoteShipID")
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
	data := spineModel.Datagram{
		Datagram: spineModel.DatagramType{
			Header: spineModel.HeaderType{},
			Payload: spineModel.PayloadType{
				Cmd: []spineModel.CmdType{},
			},
		},
	}

	msg, err := json.Marshal(data)
	assert.Nil(s.T(), err)

	err = s.sut.sendSpineData(msg)
	assert.Nil(s.T(), err)
}
