package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func TestTypesSuite(t *testing.T) {
	suite.Run(t, new(TypesSuite))
}

type TypesSuite struct {
	suite.Suite
}

func (s *TypesSuite) Test_ConnectionState() {
	conState := NewConnectionStateDetail(ConnectionStateNone, nil)
	assert.Equal(s.T(), ConnectionStateNone, conState.State())
	assert.Nil(s.T(), conState.Error())

	conState.SetState(ConnectionStateError)
	assert.Equal(s.T(), ConnectionStateError, conState.State())

	conState.SetError(errors.New("test"))
	assert.NotNil(s.T(), conState.Error())
}
