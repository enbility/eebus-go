package spine

import (
	"testing"
	"time"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type PendingRequestsTestSuite struct {
	suite.Suite
	sut     PendingRequests
	counter model.MsgCounterType
}

func TestPendingRequestsSuite(t *testing.T) {
	suite.Run(t, new(PendingRequestsTestSuite))
}

func (suite *PendingRequestsTestSuite) SetupSuite() {
	suite.counter = model.MsgCounterType(1)
	suite.sut = NewPendingRequest()
}

func (suite *PendingRequestsTestSuite) SetupTest() {
	suite.sut.Add(suite.counter, defaultMaxResponseDelay)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_Timeout() {
	suite.sut.Remove(suite.counter)
	suite.sut.Add(suite.counter, time.Duration(time.Millisecond*10))

	time.Sleep(time.Duration(time.Millisecond * 20))

	// Act
	data, err := suite.sut.GetData(suite.counter)
	assert.Nil(suite.T(), data)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), model.ErrorNumberTypeTimeout, err.ErrorNumber)
	assert.Equal(suite.T(), "the request with the message counter '1' timed out", string(err.Description))
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_Remove() {
	// Act
	err := suite.sut.Remove(suite.counter)
	assert.Nil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_Remove_GetData() {
	suite.sut.Remove(suite.counter)

	// Act
	_, err := suite.sut.GetData(suite.counter)
	assert.NotNil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetData() {
	// Act
	err := suite.sut.SetData(suite.counter, 1)
	assert.Nil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetData_UnknownCounter() {
	// Act
	err := suite.sut.SetData(model.MsgCounterType(2), 1)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), "No pending request with message counter '2' found", string(err.Description))
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetData_SetData() {
	suite.sut.SetData(suite.counter, 1)
	// Act
	err := suite.sut.SetData(suite.counter, 2)
	assert.NotNil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetResult() {
	// Act
	err := suite.sut.SetResult(suite.counter, NewErrorTypeFromString("unknown error"))
	assert.Nil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetResult_SetResult() {
	suite.sut.SetResult(suite.counter, NewErrorTypeFromString("unknown error"))
	// Act
	err := suite.sut.SetResult(suite.counter, NewErrorTypeFromString("unknown error"))
	assert.NotNil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetData_SetResult() {
	suite.sut.SetData(suite.counter, 1)
	// Act
	err := suite.sut.SetResult(suite.counter, NewErrorTypeFromString("unknown error"))
	assert.NotNil(suite.T(), err)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetData_GetData() {
	data := 1
	suite.sut.SetData(suite.counter, data)

	// Act
	result, err := suite.sut.GetData(suite.counter)
	assert.Nil(suite.T(), err)
	assert.NotNil(suite.T(), result)
	assert.Equal(suite.T(), data, result)
}

func (suite *PendingRequestsTestSuite) TestPendingRequests_SetResult_GetData() {
	errNo := model.ErrorNumberTypeTimeout
	errDesc := "Timeout occured"
	suite.sut.SetResult(suite.counter, NewErrorType(errNo, errDesc))

	// Act
	result, err := suite.sut.GetData(suite.counter)
	assert.Nil(suite.T(), result)
	assert.NotNil(suite.T(), err)
	assert.Equal(suite.T(), errNo, err.ErrorNumber)
	assert.Equal(suite.T(), errDesc, string(err.Description))
}
