package ship

import (
	"os"
	"sync"
	"testing"

	"github.com/DerAndereAndi/eebus-go/spine"
	spineModel "github.com/DerAndereAndi/eebus-go/spine/model"
)

type dataHandlerTest struct {
	sentMessage []byte

	mux sync.Mutex
}

func (s *dataHandlerTest) lastMessage() []byte {
	s.mux.Lock()
	defer s.mux.Unlock()

	return s.sentMessage
}

var _ ShipDataConnection = (*dataHandlerTest)(nil)

func (s *dataHandlerTest) InitDataProcessing(dataProcessing ShipDataProcessing) {}

func (s *dataHandlerTest) WriteMessageToDataConnection(message []byte) error {
	s.mux.Lock()
	defer s.mux.Unlock()

	s.sentMessage = message

	return nil
}

func (s *dataHandlerTest) CloseDataConnection()         {}
func (w *dataHandlerTest) IsDataConnectionClosed() bool { return false }

var _ ConnectionHandler = (*dataHandlerTest)(nil)

func (s *dataHandlerTest) HandleClosedConnection(connection *ShipConnection) {}

var _ ShipServiceDataProvider = (*dataHandlerTest)(nil)

func (s *dataHandlerTest) IsRemoteServiceForSKIPaired(string) bool      { return true }
func (s *dataHandlerTest) HandleConnectionClosed(*ShipConnection, bool) {}
func (s *dataHandlerTest) ReportServiceShipID(string, string)           {}

func initTest(role shipRole) (*ShipConnection, *dataHandlerTest) {
	localDevice := spine.NewDeviceLocalImpl("TestBrandName", "TestDeviceModel", "TestSerialNumber", "TestDeviceCode",
		"TestDeviceAddress", spineModel.DeviceTypeTypeEnergyManagementSystem, spineModel.NetworkManagementFeatureSetTypeSmart)

	dataHandler := &dataHandlerTest{}
	conhandler := NewConnectionHandler(dataHandler, dataHandler, localDevice, role, "LocalShipID", "RemoveDevice", "RemoteShipID")

	return conhandler, dataHandler
}

func shutdownTest(conhandler *ShipConnection) {
	conhandler.stopHandshakeTimer()
}

func skipCI(t *testing.T) {
	if os.Getenv("ACTION_ENVIRONMENT") == "CI" {
		t.Skip("Skipping testing in CI environment")
	}
}
