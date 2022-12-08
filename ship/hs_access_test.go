package ship

import (
	"testing"

	"github.com/enbility/eebus-go/ship/model"
	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestAccessSuite(t *testing.T) {
	suite.Run(t, new(AccessSuite))
}

type AccessSuite struct {
	suite.Suite
}

func (s *AccessSuite) Test_Init() {
	sut, data := initTest(ShipRoleClient)

	sut.setState(smePinStateCheckOk)
	sut.handleState(false, nil)

	assert.Equal(s.T(), true, sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeAccessMethodsRequest, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *AccessSuite) Test_Request() {
	sut, data := initTest(ShipRoleClient)

	sut.setState(smeAccessMethodsRequest)

	accessMsg := model.AccessMethodsRequest{
		AccessMethodsRequest: model.AccessMethodsRequestType{},
	}
	msg, err := sut.shipMessage(model.MsgTypeControl, accessMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeAccessMethodsRequest, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *AccessSuite) Test_Methods_Ok() {
	sut, data := initTest(ShipRoleClient)

	sut.setState(smeAccessMethodsRequest)

	accessMsg := model.AccessMethods{
		AccessMethods: model.AccessMethodsType{
			Id: util.Ptr("RemoteShipID"),
		},
	}
	msg, err := sut.shipMessage(model.MsgTypeControl, accessMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeComplete, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}

func (s *AccessSuite) Test_Methods_WrongShipID() {
	sut, data := initTest(ShipRoleClient)

	sut.setState(smeAccessMethodsRequest)

	accessMsg := model.AccessMethods{
		AccessMethods: model.AccessMethodsType{
			Id: util.Ptr("WrongRemoteShipID"),
		},
	}
	msg, err := sut.shipMessage(model.MsgTypeControl, accessMsg)
	assert.Nil(s.T(), err)
	assert.NotNil(s.T(), msg)

	sut.handleState(false, msg)

	assert.Equal(s.T(), false, sut.handshakeTimerRunning)
	assert.Equal(s.T(), smeError, sut.getState())
	assert.NotNil(s.T(), data.lastMessage())

	shutdownTest(sut)
}
