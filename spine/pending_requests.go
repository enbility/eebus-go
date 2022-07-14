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
	GetData(counter model.MsgCounterType) (any, *ErrorType)
	Remove(counter model.MsgCounterType) *ErrorType
}

type dataErrorPair struct {
	data        any
	errorResult *ErrorType
}

type request struct {
	counter   model.MsgCounterType
	countdown *time.Timer
	response  chan *dataErrorPair
}

func (r *request) setTimeoutResult() {
	if len(r.response) == 0 {
		errorResult := NewErrorType(model.ErrorNumberTypeTimeout, fmt.Sprintf("the request with the message counter '%s' timed out", r.counter.String()))
		r.response <- &dataErrorPair{data: nil, errorResult: errorResult}
	}
}

type PendingRequestsImpl struct {
	requestMap sync.Map
}

func NewPendingRequest() PendingRequests {
	return &PendingRequestsImpl{
		requestMap: sync.Map{},
	}
}

func (r *PendingRequestsImpl) Add(counter model.MsgCounterType, maxDelay time.Duration) {
	newRequest := &request{
		counter: counter,
		// could be a performance problem in case of many requests
		response: make(chan *dataErrorPair, 1), // buffered, so that SetData will not block,
	}
	newRequest.countdown = time.AfterFunc(maxDelay, func() { newRequest.setTimeoutResult() })

	r.requestMap.Store(counter, newRequest)
}

func (r *PendingRequestsImpl) SetData(counter model.MsgCounterType, data any) *ErrorType {
	return r.setResponse(counter, data, nil)
}

func (r *PendingRequestsImpl) SetResult(counter model.MsgCounterType, errorResult *ErrorType) *ErrorType {
	return r.setResponse(counter, nil, errorResult)
}

func (r *PendingRequestsImpl) GetData(counter model.MsgCounterType) (any, *ErrorType) {
	request, err := r.getRequest(counter)
	if err != nil {
		return nil, err
	}

	data := <-request.response
	r.removeRequest(request)

	return data.data, data.errorResult
}

func (r *PendingRequestsImpl) Remove(counter model.MsgCounterType) *ErrorType {
	request, err := r.getRequest(counter)
	if err != nil {
		return err
	}
	r.removeRequest(request)
	return nil
}

func (r *PendingRequestsImpl) removeRequest(request *request) {
	request.countdown.Stop()
	r.requestMap.Delete(request.counter)
}

func (r *PendingRequestsImpl) getRequest(counter model.MsgCounterType) (*request, *ErrorType) {
	rq, exists := r.requestMap.Load(counter)
	if !exists {
		return nil, NewErrorTypeFromString(fmt.Sprintf("No pending request with message counter '%s' found", counter.String()))
	}

	return rq.(*request), nil
}

func (r *PendingRequestsImpl) setResponse(counter model.MsgCounterType, data any, errorResult *ErrorType) *ErrorType {

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
