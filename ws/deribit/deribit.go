package deribit

import (
	"fmt"
	"log"
	"time"

	api "github.com/frankrap/deribit-api"
	"github.com/frankrap/deribit-api/models"
)

type IntermediateChannel struct {
	InputStocks         chan *models.BookSummary
	CurrencyToSellTrade map[string]chan *models.BookSummary
}

func Message(book models.BookSummary) string {
	return fmt.Sprintf("Product provide at %s expiry at xxx for %s: %v\n",
		time.UnixMilli(book.CreationTimestamp),
		book.InstrumentName,
		book.BidPrice)
}

type DeribitClient struct {
	*api.Client
	Result *IntermediateChannel
}

var BTC = "BTC"

func (client *DeribitClient) FlushPricesIntoChannelCmd() {

	// producer: read from websocket and send to channel
	go func() {
		// read from the websocket
		for {
			trade, err := client.ReadLastTradeCmd(BTC)
			if err != nil {
				break
			}
			client.Result.InputStocks <- trade
		}
		close(client.Result.InputStocks)
	}()
	// filter one kind of coin
	go func() {
		for trade := range client.Result.InputStocks {
			client.Result.CurrencyToSellTrade[BTC] <- trade
		}
		close(client.Result.CurrencyToSellTrade[BTC])
	}()

	// print the trades
	// for trade := range client.SellBTCTrades {
	// 	json, _ := json.Marshal(trade)
	// 	fmt.Println(string(json))
	// }
}
func SetupClient() *DeribitClient {
	cfg := &api.Configuration{
		Addr:          api.TestBaseURL,
		ApiKey:        "AsJTU16U",
		SecretKey:     "mM5_K8LVxztN6TjjYpv_cJVGQBvk4jglrEpqkw1b87U",
		AutoReconnect: true,
		DebugMode:     true,
	}
	client := api.New(cfg)

	client.GetTime()
	client.Test()
	result := &DeribitClient{
		Client: client,
		Result: &IntermediateChannel{
			InputStocks:         make(chan *models.BookSummary),
			CurrencyToSellTrade: make(map[string]chan *models.BookSummary),
		},
	}
	result.Result.CurrencyToSellTrade[BTC] = make(chan *models.BookSummary)
	return result
}

func (client *DeribitClient) SubscribeCmd() error {
	// Subscribe
	client.On("announcements", func(e *models.AnnouncementsNotification) {

	})
	client.On("book.ETH-PERPETUAL.100.1.100ms", func(e *models.OrderBookGroupNotification) {

	})
	client.On("book.BTC-PERPETUAL.100ms", func(e *models.OrderBookNotification) {
		fmt.Println("Order book ", e)
	})
	client.On("book.BTC-PERPETUAL.raw", func(e *models.OrderBookRawNotification) {
		fmt.Println("Order book raw ", e)
	})
	client.On("deribit_price_index.btc_usd", func(e *models.DeribitPriceIndexNotification) {

	})
	client.On("deribit_price_ranking.btc_usd", func(e *models.DeribitPriceRankingNotification) {

	})
	client.On("estimated_expiration_price.btc_usd", func(e *models.EstimatedExpirationPriceNotification) {

	})
	client.On("markprice.options.btc_usd", func(e *models.MarkpriceOptionsNotification) {

	})
	client.On("perpetual.BTC-PERPETUAL.raw", func(e *models.PerpetualNotification) {

	})
	client.On("quote.BTC-PERPETUAL", func(e *models.QuoteNotification) {

	})
	client.On("ticker.BTC-PERPETUAL.raw", func(e *models.TickerNotification) {

	})
	client.On("trades.BTC-PERPETUAL.raw", func(e *models.TradesNotification) {

	})

	client.On("user.changes.BTC-PERPETUAL.raw", func(e *models.UserChangesNotification) {

	})
	client.On("user.changes.future.BTC.raw", func(e *models.UserChangesNotification) {

	})
	client.On("user.orders.BTC-PERPETUAL.raw", func(e *models.UserOrderNotification) {

	})
	client.On("user.orders.future.BTC.100ms", func(e *models.UserOrderNotification) {

	})
	client.On("user.portfolio.btc", func(e *models.PortfolioNotification) {

	})
	client.On("user.trades.BTC-PERPETUAL.raw", func(e *models.UserTradesNotification) {

	})
	client.On("user.trades.future.BTC.100ms", func(e *models.UserTradesNotification) {

	})

	client.Subscribe([]string{
		//"announcements",
		//"book.BTC-PERPETUAL.none.10.100ms",	// none/1,2,5,10,25,100,250
		//"book.BTC-PERPETUAL.100ms",	// type: snapshot/change
		"book.BTC-PERPETUAL.raw",
		//"deribit_price_index.btc_usd",
		//"deribit_price_ranking.btc_usd",
		//"estimated_expiration_price.btc_usd",
		//"markprice.options.btc_usd",
		//"perpetual.BTC-PERPETUAL.raw",
		//"quote.BTC-PERPETUAL",
		//"ticker.BTC-PERPETUAL.raw",
		// "trades.BTC-PERPETUAL.raw",
		//"user.changes.BTC-PERPETUAL.raw",
		//"user.changes.future.BTC.raw",
		// "user.orders.BTC-PERPETUAL.raw",
		//"user.orders.future.BTC.100ms",
		//"user.portfolio.btc",
		//"user.trades.BTC-PERPETUAL.raw",
		//"user.trades.future.BTC.100ms",
	})

	return nil
}

func (client *DeribitClient) ReadLastTradeCmd(currency string) (*models.BookSummary, error) {
	if currency == "" {
		currency = BTC
	}
	getBookSummaryByCurrency, err := client.GetBookSummaryByCurrency(&models.GetBookSummaryByCurrencyParams{
		Currency: currency,
	})
	if err != nil {
		log.Printf("Error %v", err)
		return nil, err
	}
	book := &getBookSummaryByCurrency[0]
	fmt.Printf("Product provide at %s expiry at xxx for %s: %v\n",
		time.UnixMilli(book.CreationTimestamp),
		book.InstrumentName,
		book.BidPrice)
	return book, nil
}

func (client *DeribitClient) GetBookSummaryByCurrencyCmd(currency string) error {
	// GetBookSummaryByCurrency
	getBookSummaryByCurrencyParams := &models.GetBookSummaryByCurrencyParams{
		Currency: currency,
		Kind:     "future",
	}
	var getBookSummaryByCurrencyResult []models.BookSummary
	getBookSummaryByCurrencyResult, err := client.GetBookSummaryByCurrency(getBookSummaryByCurrencyParams)
	if err != nil {
		log.Printf("Error %v", err)
		return err
	}
	for _, bookSummary := range getBookSummaryByCurrencyResult {
		fmt.Printf("summary: %v \n", bookSummary)
	}
	log.Printf("getBookSummaryByCurrencyResult %v", getBookSummaryByCurrencyResult)
	return nil
}

func (client *DeribitClient) GetBookSummaryByInstrumentCmd(instrument string) error {
	if instrument == "" { 
		instrument = "BTC-16DEC22-14000-P"
	}
	// GetBookSummaryByInstrument
	getBookSummaryByInstrumentParams := &models.GetBookSummaryByInstrumentParams{
		InstrumentName: instrument,
	}
	var getBookSummaryByInstrumentResult []models.BookSummary
	getBookSummaryByInstrumentResult, err := client.GetBookSummaryByInstrument(getBookSummaryByInstrumentParams)
	if err != nil {
		log.Printf("Error %v", err)
		return err
	}
	for _, bookSummary := range getBookSummaryByInstrumentResult {
		fmt.Printf("summary: %v \n", bookSummary)
	}
	log.Printf("getBookSummaryByInstrumentResult %v", getBookSummaryByInstrumentResult)
	return nil
}

func (client *DeribitClient) GetInstrumentsCmd(currency string) error {
	// GetInstruments
	getInstrumentsParams := &models.GetInstrumentsParams{
		Currency: currency,
		Kind:     "option",
	}
	var getInstruments []models.Instrument
	getInstruments, err := client.GetInstruments(getInstrumentsParams)
	if err != nil {
		log.Printf("Error %v", err)
		return err
	}
	for _, getInstrument := range getInstruments {
		fmt.Printf("getInstrument: %v \n", getInstrument)
	}
	log.Printf("getInstruments %v", getInstruments)
	return nil
}

func (client *DeribitClient) GetOrderBookCmd() error {
	// GetOrderBook
	getOrderBookParams := &models.GetOrderBookParams{
		InstrumentName: "BTC-PERPETUAL",
		Depth:          5,
	}
	var getOrderBookResult models.GetOrderBookResponse
	getOrderBookResult, err := client.GetOrderBook(getOrderBookParams)
	if err != nil {
		log.Printf("Error %v", err)
		return err
	}
	log.Printf("getOrderBookResult %v", getOrderBookResult)
	return nil
}

func (client *DeribitClient) GetPositionCmd() error {
	// GetPosition
	getPositionParams := &models.GetPositionParams{
		InstrumentName: "BTC-PERPETUAL",
	}
	var getPositionResult models.Position
	getPositionResult, err := client.GetPosition(getPositionParams)
	if err != nil {
		log.Printf("Error %v", err)
		return err
	}
	log.Printf("getPositionResult %v", getPositionResult)
	return nil
}

func (client *DeribitClient) BuyCmd() error {
	// Buy
	guyParams := &models.BuyParams{
		InstrumentName: "BTC-PERPETUAL",
		Amount:         40,
		Price:          6000.0,
		Type:           "limit",
	}
	var buyResult models.BuyResponse
	buyResult, err := client.Buy(guyParams)
	if err != nil {
		log.Printf("Error %v", err)
		return err
	}
	log.Printf("buyResult %v", buyResult)
	return nil
}
