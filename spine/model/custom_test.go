package model

import (
	"testing"

	"github.com/enbility/eebus-go/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestErrorTypeSuite(t *testing.T) {
	suite.Run(t, new(ErrorTypeSuite))
}

type ErrorTypeSuite struct {
	suite.Suite
}

func (s *ErrorTypeSuite) SetupSuite()   {}
func (s *ErrorTypeSuite) TearDownTest() {}

func (s *ErrorTypeSuite) BeforeTest(suiteName, testName string) {}

func (s *ErrorTypeSuite) Test_NewErrorType() {
	result := NewErrorType(ErrorNumberTypeNoError, "")
	assert.NotNil(s.T(), result)
}

func (s *ErrorTypeSuite) Test_NewErrorTypeFromNumber() {
	result := NewErrorTypeFromNumber(ErrorNumberTypeCommandRejected)
	assert.NotNil(s.T(), result)
}

func (s *ErrorTypeSuite) Test_NewErrorTypeFromString() {
	result := NewErrorTypeFromString("error")
	assert.NotNil(s.T(), result)

	assert.NotEqual(s.T(), 0, len(result.String()))
}

func (s *ErrorTypeSuite) Test_NewErrorTypeFromResult() {
	input := &ResultDataType{}

	result := NewErrorTypeFromResult(input)
	assert.Nil(s.T(), result)

	input = &ResultDataType{
		ErrorNumber: util.Ptr(ErrorNumberTypeNoError),
	}

	result = NewErrorTypeFromResult(input)
	assert.Nil(s.T(), result)

	input = &ResultDataType{
		ErrorNumber: util.Ptr(ErrorNumberTypeCommandNotSupported),
	}

	result = NewErrorTypeFromResult(input)
	assert.NotNil(s.T(), result)

	assert.NotEqual(s.T(), 0, len(result.String()))
}
