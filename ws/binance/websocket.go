package binance

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var (
	// WebsocketTimeout is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WsTimeout = time.Second * 60
	// WebsocketKeepalive enables sending ping/pong messages to check the connection stability
	WsKeepalive = false

	// PingFrame is an interval for sending ping/pong messages if WebsocketKeepalive is enabled
	WsPingFrame = time.Minute * 5

	// A single connection is only valid for 24 hours; expect to be disconnected at the 24 hour mark
	WsConnectionTime = time.Hour * 24

	WsReadLimit = 655350
)

// WsHandler handle raw websocket message
type WsHandler func(message []byte)

// ErrHandler handles errors
type ErrHandler func(err error)

// WsConfig webservice configuration
type WsConfig struct {
	Endpoint          string
}

// getWsEndpoint return the base endpoint of the WS according the UseTestnet flag
func getWsEndpoint() string {
	return baseWsMainURL
}

func newWsConfig(endpoint string) *WsConfig {
	return &WsConfig{
		Endpoint: endpoint,
	}
}

func NewWSConnection(endpoint string) (*websocket.Conn, error) {
	dailer := websocket.Dialer{
		Proxy:             http.ProxyFromEnvironment,
		HandshakeTimeout:  45 * time.Second,
		EnableCompression: false,
	}
	c, _, err := dailer.Dial(endpoint, nil)
	if err != nil {
		return nil, err
	}
	return c, nil
}

var wsServe = func(endpoint string, handler WsHandler, errHandler ErrHandler) (err error) {
Daily:
	connectionTicker := time.NewTicker(WsConnectionTime)
	defer connectionTicker.Stop()

	// Get Default config first
	c, err := NewWSConnection(endpoint)
	if err != nil {
		return err
	}
	fmt.Println(" Work in ", endpoint)
	c.SetReadLimit(int64(WsReadLimit))
	stopC := make(chan struct{})

	go func() {
		// TODO: let's see if we need this
		// replyPing(c, PingFrame)

		if WsKeepalive {
			keepAlive(c, WsTimeout)
		}
		// Wait for the stopC channel to be closed.  We do that in a
		// separate goroutine because ReadMessage is a blocking
		// operation.
		silent := false
		go func() {
			<-stopC
			c.Close()
		}()

		for {

			_, message, err := c.ReadMessage()
			if err != nil {
				if !silent {
					errHandler(err)
				}
				return
			}
			handler(message)
		}
	}()
	
	// Handle daily re-create WS connection
	fmt.Println("Still wait ticker")
	<-connectionTicker.C
	stopC <- struct{}{}
	goto Daily
}

func keepAlive(c *websocket.Conn, timeout time.Duration) {
	ticker := time.NewTicker(timeout)

	lastResponse := time.Now()
	c.SetPongHandler(func(msg string) error {
		lastResponse = time.Now()
		return nil
	})

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PingMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
			if time.Since(lastResponse) > timeout {
				c.Close()
				return
			}
		}
	}()
}

func replyPing(c *websocket.Conn, pingFrame time.Duration) {
	ticker := time.NewTicker(pingFrame)

	go func() {
		defer ticker.Stop()
		for {
			deadline := time.Now().Add(10 * time.Second)
			err := c.WriteControl(websocket.PongMessage, []byte{}, deadline)
			if err != nil {
				return
			}
			<-ticker.C
		}
	}()
}
