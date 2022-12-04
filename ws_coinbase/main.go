package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

const Address = "wss://ws-feed.exchange.coinbase.com"

const (
	ChannelTicker   string = "ticker"
	TypeSubscribe   string = "subscribe"
	TypeUnSubscribe string = "unsubscribe"
)

type Message struct {
	Type       string   `json:"type"`
	ProductIDs []string `json:"product_ids"`
	Channels   []string `json:"channels"`
}

type Trade struct {
	Type    string    `json:"type"`
	Time    time.Time `json:"time"`
	Product string    `json:"product_id"`
	Price   string    `json:"price"`
}

func main() {
	// 1. websocket client connection
	conn, _, err := websocket.Dial(context.Background(), Address, nil)
	if err != nil {
		panic(err)
	}
	defer conn.Close(http.StatusOK, "connection closed")

	// 2. Subscribe
	sub := Message{
		Type:       TypeSubscribe,
		ProductIDs: []string{"BTC-USD"},
		Channels:   []string{ChannelTicker},
	}

	err = wsjson.Write(context.Background(), conn, sub)
	if err != nil {
		fmt.Printf("failed to subscribe to channel %s '%s': %v", Address, sub.ProductIDs[0], err)
		return
	}

	// 3. read from the websocket
	for {
		_, byteData, err := conn.Read(context.Background())
		if err != nil {
			fmt.Printf("Failed to read from WS channel %s '%s': %v", Address, sub.ProductIDs[0], err)
			break
		}
		var trade Trade
		err = json.Unmarshal(byteData, &trade)
		if err != nil {
			fmt.Printf("Failed to Unmarshal %v", err)
			return
		}
		fmt.Println("New price from channel: ", trade)
	}
}
