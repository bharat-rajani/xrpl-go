package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/oklog/ulid/v2"
	"github.com/xrpscan/xrpl-go/models/methods"
	"sync"
	"time"
)

const (
	defaultRequestTimeout = time.Minute * 5
)

type RequestManager struct {
	m *sync.Mutex // mutex for responseChannels
	// TODO (bharat-rajani):can use concurrent safe map here to avoid mutex misuse
	responseChannels map[ulid.ULID]chan ResponseTupleFn
}

func NewRequestManager() *RequestManager {
	return &RequestManager{
		responseChannels: make(map[ulid.ULID]chan ResponseTupleFn),
		m:                new(sync.Mutex),
	}
}

type ResponseTupleFn func() (*methods.AccountInfoResponse, error)

// CreateRequest
func (rm *RequestManager) RegisterRequest(request methods.AccountInfoRequest) (
	ulid.ULID, methods.AccountInfoRequest, <-chan ResponseTupleFn, context.CancelFunc) {

	_, ok := request.Context().Deadline()
	var ctx context.Context
	var cancelF context.CancelFunc
	// if input request has no timeout specified then we should derive our own context
	// also, we should not (cannot) inject this derived context in request
	if !ok {
		ctx, cancelF = context.WithTimeout(request.Context(), defaultRequestTimeout)
	}
	id := ulid.Make()
	respChan := make(chan ResponseTupleFn)

	// write to responseChannel map safely
	rm.m.Lock()
	rm.responseChannels[id] = respChan
	rm.m.Unlock()

	go func(ctx context.Context) {
		for {
			select {
			case <-ctx.Done():
				// cleanup
				rm.m.Lock()
				c := rm.responseChannels[id]
				c <- func() (*methods.AccountInfoResponse, error) {
					return nil, errors.New("timeout: no response received from ripple servers")
				}
				delete(rm.responseChannels, id)
				rm.m.Unlock()
			default:

			}
		}
	}(ctx)

	// inject the id in the request
	request.BaseRequest.Id = id.String()
	return id, request, respChan, cancelF
}

func (rm *RequestManager) ResolveRequest(id string, data []byte) error {

	uid, err := ulid.Parse(id)
	if err != nil {
		return err
	}

	rm.m.Lock()
	ch, ok := rm.responseChannels[uid]
	if !ok {
		return fmt.Errorf("no associated request for id %s", id)
	}
	rm.m.Unlock()

	// id found make sure to cleanup
	defer func() {
		close(ch)
		rm.m.Lock()
		delete(rm.responseChannels, uid)
		rm.m.Unlock()
	}()

	var accountInfoResponse methods.AccountInfoResponse
	err = json.Unmarshal(data, &accountInfoResponse)
	if err != nil {
		ch <- func() (*methods.AccountInfoResponse, error) {
			return nil, err
		}
		return nil // important return
	}
	ch <- func() (*methods.AccountInfoResponse, error) {
		return &accountInfoResponse, nil
	}
	return nil
}
