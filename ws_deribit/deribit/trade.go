package deribit

import (
	"fmt"
	"time"

	"github.com/frankrap/deribit-api/models"
)

func Message(trade models.Trade) string {
	return fmt.Sprintf("Last price at %s for %s: %v\n",
		time.Unix(trade.Timestamp, 0).String(),
		trade.InstrumentName,
		trade.Price)
}