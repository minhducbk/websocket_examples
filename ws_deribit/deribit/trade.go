package deribit

import (
	"fmt"
	"time"

	"github.com/frankrap/deribit-api/models"
)

type IntermediateChannel struct {
	InputStocks         chan *models.Trade
	CurrencyToSellTrade map[string]chan *models.Trade
}

func Message(trade models.Trade) string {
	return fmt.Sprintf("Last price at %s for %s: %v\n",
		time.UnixMilli(trade.Timestamp),
		trade.InstrumentName,
		trade.Price)
}
