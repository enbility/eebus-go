package spine

import (
	"errors"
	"fmt"
	"time"

	"github.com/DerAndereAndi/eebus-go/spine/model"
	"github.com/DerAndereAndi/eebus-go/util"
)

type dataErrorPair struct {
	data        any
	errorResult *ErrorType
}

type PendingRequests map[model.MsgCounterType]chan *dataErrorPair

func (r PendingRequests) Add(counter model.MsgCounterType) {
	requestChannel := make(chan *dataErrorPair, 1) // buffered, so that SetData will not block
	r[counter] = requestChannel
}

func (r PendingRequests) SetData(counter model.MsgCounterType, data any) *ErrorType {
	requestChannel, err := r.getEmptyChannel(counter)
	if err != nil {
		return err
	}
	requestChannel <- &dataErrorPair{data: data, errorResult: nil}
	return nil
}

func (r PendingRequests) SetResult(counter model.MsgCounterType, errorResult *ErrorType) *ErrorType {
	requestChannel, err := r.getEmptyChannel(counter)
	if err != nil {
		return err
	}
	requestChannel <- &dataErrorPair{data: nil, errorResult: errorResult}
	return nil
}

func (r PendingRequests) GetData(counter model.MsgCounterType, maxDelay time.Duration) RequestDataResult {
	requestChannel, exists := r[counter]
	if !exists {
		return *NewRequestDataResultError(errors.New("No pending request found"))
	}

	data := util.ReceiveWithTimeout(requestChannel, maxDelay)
	delete(r, counter)
	if data == nil {
		return *NewRequestDataResult(nil, NewErrorType(model.ErrorNumberTypeTimeout, "Timeout occured"))
	}
	return *NewRequestDataResult(data.data, data.errorResult)
}

func (r PendingRequests) getEmptyChannel(counter model.MsgCounterType) (chan *dataErrorPair, *ErrorType) {
	requestChannel, exists := r[counter]
	if !exists {
		return nil, NewErrorTypeFromString("No pending request found")
	}
	if len(requestChannel) > 0 {
		return nil, NewErrorTypeFromString(fmt.Sprintf("the Data or Result for the request (MsgCounter: %s) was already set!", &counter))
	}

	return requestChannel, nil
}
