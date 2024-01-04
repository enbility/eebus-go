package ship

import (
	"errors"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestWebsocketSuite(t *testing.T) {
	suite.Run(t, new(WebsocketSuite))
}

type WebsocketSuite struct {
	suite.Suite

	sut *websocketConnection

	testServer *httptest.Server
	testWsConn *websocket.Conn

	shipDataProcessing *MockShipDataProcessing
}

func (s *WebsocketSuite) SetupSuite()   {}
func (s *WebsocketSuite) TearDownTest() {}

func (s *WebsocketSuite) BeforeTest(suiteName, testName string) {
	ctrl := gomock.NewController(s.T())

	s.shipDataProcessing = NewMockShipDataProcessing(ctrl)
	s.shipDataProcessing.EXPECT().ReportConnectionError(gomock.Any()).AnyTimes()
	s.shipDataProcessing.EXPECT().HandleIncomingShipMessage(gomock.Any()).AnyTimes()

	ts := &testServer{}
	s.testServer, s.testWsConn = newWSServer(s.T(), ts)

	s.sut = NewWebsocketConnection(s.testWsConn, "remoteSki")
	s.sut.InitDataProcessing(s.shipDataProcessing)
}

func (s *WebsocketSuite) AfterTest(suiteName, testName string) {
	s.testWsConn.Close()
	s.testServer.Close()
}

func (s *WebsocketSuite) TestConnection() {
	isClosed := s.sut.isConnClosed()
	assert.Equal(s.T(), false, isClosed)

	msg := []byte{0, 0}
	err := s.sut.WriteMessageToDataConnection(msg)
	assert.Nil(s.T(), err)

	// make sure we have enough time to read and write
	time.Sleep(time.Millisecond * 500)

	msg = []byte{1}
	msg = append(msg, []byte("message")...)
	err = s.sut.WriteMessageToDataConnection(msg)
	assert.Nil(s.T(), err)

	// make sure we have enough time to read and write
	time.Sleep(time.Millisecond * 500)

	isConnClosed, err := s.sut.IsDataConnectionClosed()
	assert.Equal(s.T(), false, isConnClosed)
	assert.Nil(s.T(), err)

	s.sut.CloseDataConnection(450, "User Close")

	isConnClosed, err = s.sut.IsDataConnectionClosed()
	assert.Equal(s.T(), true, isConnClosed)
	assert.NotNil(s.T(), err)

	err = s.sut.WriteMessageToDataConnection(msg)
	assert.NotNil(s.T(), err)
}

func (s *WebsocketSuite) TestConnectionInvalid() {
	msg := []byte{100}
	err := s.sut.WriteMessageToDataConnection(msg)
	assert.Nil(s.T(), err)

	// make sure we have enough time to read and write
	time.Sleep(time.Millisecond * 500)

	isConnClosed, err := s.sut.IsDataConnectionClosed()
	assert.Equal(s.T(), true, isConnClosed)
	assert.NotNil(s.T(), err)
}

func (s *WebsocketSuite) TestConnectionClose() {
	s.sut.close()

	isClosed, err := s.sut.IsDataConnectionClosed()
	assert.Equal(s.T(), true, isClosed)
	assert.NotNil(s.T(), err)
}

func (s *WebsocketSuite) TestPingPeriod() {
	isClosed, err := s.sut.IsDataConnectionClosed()
	assert.Equal(s.T(), false, isClosed)
	assert.Nil(s.T(), err)

	if !isRunningOnCI() {
		// test if the function is triggered correctly via the timer
		time.Sleep(time.Second * 51)
	} else {
		// speed up the test by running the method directly
		s.sut.handlePing()
	}

	isClosed, err = s.sut.IsDataConnectionClosed()
	assert.Equal(s.T(), false, isClosed)
	assert.Nil(s.T(), err)
}

func (s *WebsocketSuite) TestCloseWithError() {
	isClosed, err := s.sut.IsDataConnectionClosed()
	assert.Equal(s.T(), false, isClosed)
	assert.Nil(s.T(), err)

	err = errors.New("test error")
	s.sut.closeWithError(err, "test error")

	isClosed, err = s.sut.IsDataConnectionClosed()
	assert.Equal(s.T(), true, isClosed)
	assert.NotNil(s.T(), err)
}

var upgrader = websocket.Upgrader{}

func newWSServer(t *testing.T, h http.Handler) (*httptest.Server, *websocket.Conn) {
	t.Helper()

	s := httptest.NewServer(h)
	wsURL := strings.Replace(s.URL, "http://", "ws://", -1)
	wsURL = strings.Replace(wsURL, "https://", "wss://", -1)

	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatal(err)
	}

	return s, ws
}

type testServer struct {
}

func (s *testServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer ws.Close()

	for {
		_, msg, err := ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			return
		}

		err = ws.WriteMessage(websocket.BinaryMessage, msg)
		if err != nil {
			continue
		}
	}
}
