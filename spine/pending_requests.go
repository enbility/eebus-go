package spine

import (
	"fmt"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type PendingRequests interface {
	Add(counter model.MsgCounterType, maxDelay time.Duration)
	SetData(counter model.MsgCounterType, data any) *ErrorType
	SetResult(counter model.MsgCounterType, errorResult *ErrorType) *ErrorType
	GetData(counter model.MsgCounterType) RequestDataResult
	Remove(counter model.MsgCounterType) *ErrorType
}

type dataErrorPair struct {
	data        any
	errorResult *ErrorType
}

type request struct {
	countdown *time.Timer
	response  chan *dataErrorPair
}

type PendingRequestsImpl struct {
	mu         sync.Mutex
	requestMap map[model.MsgCounterType]*request
}

func NewPendingRequest() PendingRequests {
	return &PendingRequestsImpl{
		requestMap: make(map[model.MsgCounterType]*request),
	}
}

func (r *PendingRequestsImpl) Add(counter model.MsgCounterType, maxDelay time.Duration) {
	r.mu.Lock()
	defer r.mu.Unlock()

	newRequest := &request{
		countdown: time.AfterFunc(maxDelay, func() { r.setTimeoutResult(counter) }),
		// could be a performance problem in case of many requests
		response: make(chan *dataErrorPair, 1), // buffered, so that SetData will not block,
	}

	r.requestMap[counter] = newRequest
}

func (r *PendingRequestsImpl) SetData(counter model.MsgCounterType, data any) *ErrorType {
	return r.setResponse(counter, data, nil)
}

func (r *PendingRequestsImpl) SetResult(counter model.MsgCounterType, errorResult *ErrorType) *ErrorType {
	return r.setResponse(counter, nil, errorResult)
}

func (r *PendingRequestsImpl) GetData(counter model.MsgCounterType) RequestDataResult {
	request, err := r.getRequest(counter)
	if err != nil {
		return *NewRequestDataResult(nil, err)
	}

	data := <-request.response
	err = r.Remove(counter)

	if err != nil {
		return *NewRequestDataResult(nil, err)
	}
	return *NewRequestDataResult(data.data, data.errorResult)
}

func (r *PendingRequestsImpl) Remove(counter model.MsgCounterType) *ErrorType {
	r.mu.Lock()
	defer r.mu.Unlock()

	request, err := r.getRequest(counter)
	if err != nil {
		return err
	}
	request.countdown.Stop()

	delete(r.requestMap, counter)
	return nil
}

func (r *PendingRequestsImpl) getRequest(counter model.MsgCounterType) (*request, *ErrorType) {
	request, exists := r.requestMap[counter]
	if !exists {
		return nil, NewErrorTypeFromString(fmt.Sprintf("No pending request with message counter '%s' found", counter.String()))
	}

	return request, nil
}

func (r *PendingRequestsImpl) setTimeoutResult(counter model.MsgCounterType) {
	r.mu.Lock()
	defer r.mu.Unlock()

	request, err := r.getRequest(counter)
	if err == nil {
		if len(request.response) == 0 {
			errorResult := NewErrorType(model.ErrorNumberTypeTimeout, fmt.Sprintf("the request with the message counter '%s' timed out", counter.String()))
			request.response <- &dataErrorPair{data: nil, errorResult: errorResult}
		}
	}
}

func (r *PendingRequestsImpl) setResponse(counter model.MsgCounterType, data any, errorResult *ErrorType) *ErrorType {
	r.mu.Lock()
	defer r.mu.Unlock()

	request, err := r.getRequest(counter)
	if err != nil {
		return err
	}
	if len(request.response) > 0 {
		return NewErrorTypeFromString(fmt.Sprintf("the Data or Result for the request (MsgCounter: %s) was already set!", &counter))
	}

	request.countdown.Stop()
	request.response <- &dataErrorPair{data: data, errorResult: errorResult}
	return nil
}
