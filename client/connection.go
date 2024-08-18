package client

import (
	"context"
	"errors"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/oklog/ulid/v2"
	"log"
	"sync"
)

// WebsocketConnectionManager is a websocket manager
type WebsocketConnectionManager struct {
	connections map[ulid.ULID]*WebsocketConnection
	mutex       *sync.Mutex // Protects connections map
}

// NewWebsocketConnection will create a websocket connection and register it with manager
func (wcm *WebsocketConnectionManager) NewWebsocketConnection(ctx context.Context, url string) (*WebsocketConnection, error) {

	newConn, err := NewWebsocketConnection(ctx, url)
	if err != nil {
		return nil, err
	}
	wcm.connections[newConn.ConnId] = newConn
	wcm.mutex.Lock()
	defer wcm.mutex.Unlock()
	// TODO (bharat-rajani):Figure out how to maintain heartbeats and keepalives through manager and reconnect if needed
	return newConn, nil
}

// Close will close and unregister all the websocket connections associated with the manager
func (wcm *WebsocketConnectionManager) Close() error {
	wcm.mutex.Lock()
	defer wcm.mutex.Unlock()
	var errs []error
	for _, conn := range wcm.connections {
		if err := conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}
	return errors.Join(errs...)
}

// CloseByID will close and unregister the websocket connection from manager by id
func (wcm *WebsocketConnectionManager) CloseByID(id ulid.ULID) error {
	wcm.mutex.Lock()
	defer wcm.mutex.Unlock()
	conn, ok := wcm.connections[id]
	if !ok {
		return errors.New("connection not found")
	}
	if err := conn.Close(); err != nil {
		return err
	}
	return nil
}

func NewWebsocketConnectionManager() *WebsocketConnectionManager {
	return &WebsocketConnectionManager{connections: make(map[ulid.ULID]*WebsocketConnection), mutex: new(sync.Mutex)}
}

// WebsocketConnection contains underlying websocket connection with concurrent safe usage and reconnect mechanism
type WebsocketConnection struct {
	url    string
	conn   *websocket.Conn
	ConnId ulid.ULID
	mutex  *sync.Mutex // Protects concurrent operations on underlying connection
	alive  bool
}

func (wc *WebsocketConnection) WriteMessage(ctx context.Context, messageType int, data []byte) error {
	if !wc.alive {
		return errors.New("websocket connection is dead")
	}
	wc.mutex.Lock()
	defer wc.mutex.Unlock()
	if deadline, ok := ctx.Deadline(); ok {
		if err := wc.conn.SetWriteDeadline(deadline); err != nil {
			return err
		}
	}
	return wc.conn.WriteMessage(messageType, data)
}

func (wc *WebsocketConnection) ReadMessage(c chan []byte) error {
	defer close(c)
	for {
		if !wc.alive {
			return errors.New("websocket connection is dead")
		}
		messageType, message, err := wc.conn.ReadMessage()
		if err != nil {
			return err
		}

		switch messageType {
		case websocket.CloseMessage:
			log.Println("WS websocket.CloseMessage received")
			return nil
		case websocket.TextMessage:
			c <- message
		case websocket.BinaryMessage:
		default:
		}
	}
	return nil
}

func (wc *WebsocketConnection) Reconnect() error {
	if wc.url == "" {
		return errors.New("websocket connection url is empty")
	}
	if wc.alive {
		return errors.New("websocket connection is already alive")
	}
	// TODO (bharat-rajani):should we save context in websocket connection object (anti-pattern whose use-cases are niche)?
	conn, err := dial(context.TODO(), wc.url)
	if err != nil {
		return fmt.Errorf("error while reconnecting %w", err)
	}
	wc.conn = conn
	return nil
}

func (wc *WebsocketConnection) Close() error {
	wc.mutex.Lock()
	defer wc.mutex.Unlock()

	err := wc.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
	if err != nil {
		log.Println("WS write error: ", err)
		return err
	}
	err = wc.conn.Close()
	if err != nil {
		log.Println("WS close error: ", err)
		return err
	}
	return nil
}

func dial(ctx context.Context, url string) (*websocket.Conn, error) {
	conn, r, err := websocket.DefaultDialer.DialContext(ctx, url, nil)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()
	return conn, nil
}

func NewWebsocketConnection(ctx context.Context, url string) (*WebsocketConnection, error) {

	conn, err := dial(ctx, url)
	if err != nil {
		return nil, err
	}
	wc := &WebsocketConnection{
		url:    url,
		conn:   conn,
		ConnId: ulid.Make(),
		mutex:  &sync.Mutex{},
		alive:  true,
	}
	return wc, nil
}
