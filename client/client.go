package client

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/xrpscan/xrpl-go/models/methods"
	"github.com/xrpscan/xrpl-go/pkg/ripplexcodec"
	"log"
)

//var DefaultClient = &Client{}

const (
	RippledApiV1      = 1
	RippledApiV2      = 2
	DefaultApiVersion = RippledApiV2
)

// TODO client config
type ClientConfig struct {
}

type Client struct {
	connection     *WebsocketConnection
	requestManager *RequestManager
	requestQueue   chan methods.AccountInfoRequest
}

func NewClient(ctx context.Context, conn *WebsocketConnection, clientConfig ClientConfig) *Client {
	return &Client{connection: conn, requestManager: NewRequestManager()}
}

func (c *Client) Do(req methods.AccountInfoRequest) (any, error) {
	return c.do(req)
}

var testHookClientDoResult func(retres any, reterr error)

func (c *Client) do(request methods.AccountInfoRequest) (retres any, reterr error) {
	if testHookClientDoResult != nil {
		defer func() { testHookClientDoResult(retres, reterr) }()
	}

	// client request is registered with the req manager,
	// manager maintains the response mappings
	// once response is received, check if it's id mapping,
	// if present then send the response to associated channel where do(er) would be listening and close the channel
	id, request, respChan, cancelF := c.requestManager.RegisterRequest(request)
	request.BaseRequest.ApiVersion = DefaultApiVersion
	request.BaseRequest.Command = request.Command()
	classicAddress, err := EnsureClassicAddress(request.Account)
	if err != nil {
		return methods.AccountInfoResponse{}, err
	}
	request.Account = classicAddress
	data, err := json.Marshal(request)
	if err != nil {
		// cancel to notify the request manager of key deletion
		if cancelF != nil {
			cancelF()
		}
		return nil, fmt.Errorf("cannot send request due to json marshal error: %w", err)
	}
	err = c.connection.WriteMessage(request.Context(), websocket.TextMessage, data)
	if err != nil {
		if cancelF != nil {
			cancelF()
		}
		return nil, fmt.Errorf("websocket error: %w", err)
	}

	log.Println(id)
	respFn := <-respChan
	return respFn()
}

// GetXRPBalance
func (c *Client) GetXRPBalance(ctx context.Context, address string, options ...methods.AccountInfoRequestOption) (json.Number, error) {
	xrpRequest := methods.AccountInfoRequest{
		Account: address,
	}

	for _, o := range options {
		o(&xrpRequest)
	}

	resp, err := c.Do(xrpRequest)
	if err != nil {
		return "", err
	}
	fmt.Println(resp)
	return "", nil
}

func EnsureClassicAddress(address string) (string, error) {
	if ripplexcodec.IsValidXAddress(address) {
		address, tag, _ := ripplexcodec.XAddressToClassicAddress(address)
		if tag != nil {
			return "", errors.New("this command does not support the use of a tag. Use an address without a tag")
		}
		return address, nil
	}
	return address, nil
}

func (c *Client) ReadMessages() error {
	log.Println("started reading")
	comm := make(chan []byte)

	go func() {
		err := c.connection.ReadMessage(comm)
		if err != nil {
			log.Println(err)
		}
	}()

	for {
		select {
		case data, ok := <-comm:
			fmt.Println("read from comm chan")
			if !ok {
				return errors.New("client communication channel with ReadMessage is closed")
			}
			var mp map[string]any
			if err := json.Unmarshal(data, &mp); err != nil {
				log.Println("json.Unmarshal error: ", err)
			}
			id, ok := mp["id"].(string)
			fmt.Println(mp)
			if !ok {
				log.Println("received response has no string id")
			}
			err := c.requestManager.ResolveRequest(id, data)
			if err != nil {
				log.Println(err.Error())
			}
		}
	}
}
