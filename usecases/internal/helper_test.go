package internal

import (
	spineapi "github.com/enbility/spine-go/api"
	"github.com/enbility/spine-go/mocks"
	"github.com/stretchr/testify/assert"
)

func (s *InternalSuite) Test_IsDeviceConnected() {
	payload := spineapi.EventPayload{}
	result := IsDeviceConnected(payload)
	assert.Equal(s.T(), false, result)

	device := mocks.NewDeviceRemoteInterface(s.T())
	payload = spineapi.EventPayload{
		Device:     device,
		EventType:  spineapi.EventTypeDeviceChange,
		ChangeType: spineapi.ElementChangeAdd,
	}
	result = IsDeviceConnected(payload)
	assert.Equal(s.T(), true, result)
}

func (s *InternalSuite) Test_IsDeviceDisconnected() {
	payload := spineapi.EventPayload{}
	result := IsDeviceDisconnected(payload)
	assert.Equal(s.T(), false, result)

	device := mocks.NewDeviceRemoteInterface(s.T())
	payload = spineapi.EventPayload{
		Ski:        "12345",
		Device:     device,
		EventType:  spineapi.EventTypeDeviceChange,
		ChangeType: spineapi.ElementChangeRemove,
	}
	result = IsDeviceDisconnected(payload)
	assert.Equal(s.T(), true, result)
}

func (s *InternalSuite) Test_IsEntityAdded() {
	payload := spineapi.EventPayload{}
	result := IsEntityAdded(payload)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity:     s.evseEntity,
		EventType:  spineapi.EventTypeEntityChange,
		ChangeType: spineapi.ElementChangeAdd,
	}
	result = IsEntityAdded(payload)
	assert.Equal(s.T(), true, result)
}

func (s *InternalSuite) Test_IsEntityRemoved() {
	payload := spineapi.EventPayload{}
	result := IsEntityRemoved(payload)
	assert.Equal(s.T(), false, result)

	payload = spineapi.EventPayload{
		Entity:     s.evseEntity,
		EventType:  spineapi.EventTypeEntityChange,
		ChangeType: spineapi.ElementChangeRemove,
	}
	result = IsEntityRemoved(payload)
	assert.Equal(s.T(), true, result)
}
