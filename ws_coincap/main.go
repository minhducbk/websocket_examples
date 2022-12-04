package main

import (
	"encoding/json"
	"fmt"

	"github.com/gorilla/websocket"
)

type Trade struct {
	Exchange  string  `json:"exchange"`
	Base      string  `json:"base"`
	Quote     string  `json:"quote"`
	Direction string  `json:"direction"`
	Price     float64 `json:"price"`
	Volume    int64   `json:"volume"`
	Timestamp int64   `json:"timestamp"`
	PriceUsd  float64 `json:"priceUsd"`
}

func main() {

	// web socket is

	// websocket source
	c, _, err := websocket.DefaultDialer.Dial("wss://ws.coincap.io/trades/binance", nil)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	if err != nil {
		panic(err)
	}
	defer c.Close()

	// create the input channel
	inputStocks := make(chan Trade)
	dogecoin := make(chan Trade)

	// producer: read from websocket and send to channel
	go func() {
		// read from the websocket
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				break
			}
			// unmarshal the message
			var trade Trade
			json.Unmarshal(message, &trade)
			// send the trade to the channel
			inputStocks <- trade
		}
		close(inputStocks)
	}()
	// filter one kind of coin
	go func() {
		for trade := range inputStocks {
			if trade.Base == "dogecoin" && trade.Quote == "tether" {
				dogecoin <- trade
			}
		}
		close(dogecoin)
	}()

	// print the trades
	for trade := range dogecoin {
		json, _ := json.Marshal(trade)
		fmt.Println(string(json))
	}
}