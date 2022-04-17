package spine

import "github.com/DerAndereAndi/eebus-go/spine/model"

type PendingRequests[T any] map[model.MsgCounterType]chan T

func (r PendingRequests[T]) Add(counter model.MsgCounterType, requestChannel chan T) {
	r[counter] = requestChannel
}

func (r PendingRequests[T]) Handle(counter model.MsgCounterType, data T) {
	requestChannel, exists := r[counter]
	if exists {
		requestChannel <- data
		delete(r, counter)
	}
}
