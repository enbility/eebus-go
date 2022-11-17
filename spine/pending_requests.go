package spine

import (
	"fmt"
	"sync"
	"time"

	"github.com/DerAndereAndi/eebus-go/spine/model"
)

type PendingRequests interface {
	Add(ski string, counter model.MsgCounterType, maxDelay time.Duration)
	SetData(ski string, counter model.MsgCounterType, data any) *ErrorType
	SetResult(ski string, counter model.MsgCounterType, errorResult *ErrorType) *ErrorType
	GetData(ski string, counter model.MsgCounterType) (any, *ErrorType)
	Remove(ski string, counter model.MsgCounterType) *ErrorType
}

type dataErrorPair struct {
	data        any
	errorResult *ErrorType
}

type request struct {
	ski       string // can not use device, as this is not available for the very first requests!
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

func (r *PendingRequestsImpl) Add(ski string, counter model.MsgCounterType, maxDelay time.Duration) {
	newRequest := &request{
		ski:     ski,
		counter: counter,
		// could be a performance problem in case of many requests
		response: make(chan *dataErrorPair, 1), // buffered, so that SetData will not block,
	}
	newRequest.countdown = time.AfterFunc(maxDelay, func() { newRequest.setTimeoutResult() })

	r.requestMap.Store(r.mapKey(ski, counter), newRequest)
}

func (r *PendingRequestsImpl) SetData(ski string, counter model.MsgCounterType, data any) *ErrorType {
	return r.setResponse(ski, counter, data, nil)
}

func (r *PendingRequestsImpl) SetResult(ski string, counter model.MsgCounterType, errorResult *ErrorType) *ErrorType {
	return r.setResponse(ski, counter, nil, errorResult)
}

func (r *PendingRequestsImpl) GetData(ski string, counter model.MsgCounterType) (any, *ErrorType) {
	request, err := r.getRequest(ski, counter)
	if err != nil {
		return nil, err
	}

	data := <-request.response
	r.removeRequest(request)

	return data.data, data.errorResult
}

func (r *PendingRequestsImpl) Remove(ski string, counter model.MsgCounterType) *ErrorType {
	request, err := r.getRequest(ski, counter)
	if err != nil {
		return err
	}
	r.removeRequest(request)
	return nil
}

func (r *PendingRequestsImpl) mapKey(ski string, counter model.MsgCounterType) string {
	return fmt.Sprintf("%s:%d", ski, counter)
}

func (r *PendingRequestsImpl) removeRequest(request *request) {
	request.countdown.Stop()
	r.requestMap.Delete(r.mapKey(request.ski, request.counter))
}

func (r *PendingRequestsImpl) getRequest(ski string, counter model.MsgCounterType) (*request, *ErrorType) {
	rq, exists := r.requestMap.Load(r.mapKey(ski, counter))
	if !exists {
		return nil, NewErrorTypeFromString(fmt.Sprintf("No pending request with message counter '%s' found", counter.String()))
	}

	return rq.(*request), nil
}

func (r *PendingRequestsImpl) setResponse(ski string, counter model.MsgCounterType, data any, errorResult *ErrorType) *ErrorType {

	request, err := r.getRequest(ski, counter)
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
