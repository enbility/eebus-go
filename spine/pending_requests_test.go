package spine

import (
	"testing"
	"time"

	"github.com/enbility/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PendingRequestsTestSuite struct {
	suite.Suite
	sut     PendingRequests
	ski     string
	counter model.MsgCounterType
}

func TestPendingRequestsSuite(t *testing.T) {
	suite.Run(t, new(PendingRequestsTestSuite))
}

func (suite *PendingRequestsTestSuite) SetupSuite() {
	suite.counter = model.MsgCounterType(1)
	suite.sut = NewPendingRequest()
	suite.ski = "test"
}

func (suite *PendingRequestsTestSuite) SetupTest() {
	suite.sut.Add(suite.ski, suite.counter, defaultMaxResponseDelay)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_Timeout() {
	_ = suite.sut.Remove(suite.ski, suite.counter)
	suite.sut.Add(suite.ski, suite.counter, time.Duration(time.Millisecond*10))

	time.Sleep(time.Duration(time.Millisecond * 20))

	// Act
	data, err := suite.sut.GetData(suite.ski, suite.counter)
	assert.Nil(suite.T(), data)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), model.ErrorNumberTypeTimeout, err.ErrorNumber)
	assert.Equal(suite.T(), "the request with the message counter '1' timed out", string(*err.Description))
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_Remove() {
	// Act
	err := suite.sut.Remove(suite.ski, suite.counter)
	assert.Nil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_Remove_GetData() {
	_ = suite.sut.Remove(suite.ski, suite.counter)

	// Act
	_, err := suite.sut.GetData(suite.ski, suite.counter)
	assert.NotNil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetData() {
	// Act
	err := suite.sut.SetData(suite.ski, suite.counter, 1)
	assert.Nil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetData_UnknownCounter() {
	// Act
	err := suite.sut.SetData(suite.ski, model.MsgCounterType(2), 1)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "No pending request with message counter '2' found", string(*err.Description))
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetData_SetData() {
	_ = suite.sut.SetData(suite.ski, suite.counter, 1)
	// Act
	err := suite.sut.SetData(suite.ski, suite.counter, 2)
	assert.NotNil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetResult() {
	// Act
	err := suite.sut.SetResult(suite.ski, suite.counter, NewErrorTypeFromString("unknown error"))
	assert.Nil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetResult_SetResult() {
	_ = suite.sut.SetResult(suite.ski, suite.counter, NewErrorTypeFromString("unknown error"))
	// Act
	err := suite.sut.SetResult(suite.ski, suite.counter, NewErrorTypeFromString("unknown error"))
	assert.NotNil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetData_SetResult() {
	_ = suite.sut.SetData(suite.ski, suite.counter, 1)
	// Act
	err := suite.sut.SetResult(suite.ski, suite.counter, NewErrorTypeFromString("unknown error"))
	assert.NotNil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetData_GetData() {
	data := 1
	_ = suite.sut.SetData(suite.ski, suite.counter, data)

	// Act
	result, err := suite.sut.GetData(suite.ski, suite.counter)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), data, result)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetResult_GetData() {
	errNo := model.ErrorNumberTypeTimeout
	errDesc := "Timeout occurred"
	_ = suite.sut.SetResult(suite.ski, suite.counter, NewErrorType(errNo, errDesc))

	// Act
	result, err := suite.sut.GetData(suite.ski, suite.counter)
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errNo, err.ErrorNumber)
	assert.Equal(suite.T(), errDesc, string(*err.Description))
}
