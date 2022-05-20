package spine

import (
	"errors"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type PendingRequests[T any] map[model.MsgCounterType]chan T

func (r PendingRequests[T]) Add(counter model.MsgCounterType, requestChannel chan T) {
	r[counter] = requestChannel
}

func (r PendingRequests[T]) Handle(counter model.MsgCounterType, data T) error {
	requestChannel, exists := r[counter]
	if exists {
		requestChannel <- data
		delete(r, counter)
		return nil
	}
	return errors.New("No pending request found")
}
