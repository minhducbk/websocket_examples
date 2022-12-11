package binance

import (
	"encoding/json"
	"fmt"
)

// Option WS Endpoints
const (
	baseWsMainURL          = "wss://nbstream.binance.com/eoptions/ws"
)


// WsMarketStatHandler handle websocket that push single market statistics for 24hr
type WsMarketStatHandler func(event *WsMarketStatEvent)


// WsMarketStatEvent define websocket market statistics event
type WsMarketStatEvent struct {
	Event              string `json:"e"`
	Time               int64  `json:"E"`
	Symbol             string `json:"s"`
	PriceChange        string `json:"p"`
	PriceChangePercent string `json:"P"`
	WeightedAvgPrice   string `json:"w"`
	PrevClosePrice     string `json:"x"`
	LastPrice          string `json:"c"`
	CloseQty           string `json:"Q"`
	BidPrice           string `json:"b"`
	BidQty             string `json:"B"`
	AskPrice           string `json:"a"`
	AskQty             string `json:"A"`
	OpenPrice          string `json:"o"`
	HighPrice          string `json:"h"`
	LowPrice           string `json:"l"`
	BaseVolume         string `json:"v"`
	QuoteVolume        string `json:"q"`
	OpenTime           string  `json:"O"`
	CloseTime          string  `json:"C"`
	FirstID            string  `json:"F"`
	LastID             string  `json:"L"`
	Count              int64  `json:"n"`
}


// WsMarketStatServe serve websocket that push 24hr statistics for single market every second
func WsTickerServe(symbol string, handler WsMarketStatHandler, errHandler ErrHandler) (err error) {
	endpoint := fmt.Sprintf("%s/%s@ticker", getWsEndpoint(), symbol)
	wsHandler := func(message []byte) {
		var event WsMarketStatEvent
		err := json.Unmarshal(message, &event)
		if err != nil {
			errHandler(err)
			return
		}
		handler(&event)
	}
	return wsServe(endpoint, wsHandler, errHandler)
}