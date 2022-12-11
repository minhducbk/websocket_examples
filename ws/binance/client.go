package binance

import (
	"encoding/json"
	"fmt"
)

// wsDepthHandler := func(event *binance.WsDepthEvent) {
//     fmt.Println(event)
// }
// errHandler := func(err error) {
//     fmt.Println(err)
// }

func ShowBinancePrice(symbol string) chan []byte {
	outputChannel := make(chan []byte)
	go func() {
		wsDepthHandler := func(event *WsMarketStatEvent) {
			fmt.Println(event)
			res2B, _ := json.Marshal(event)
			outputChannel <- res2B
			fmt.Println("WsMarketStatEvent JSON: ", string(res2B))
			fmt.Println("------")
		}
		errHandler := func(err error) {
			fmt.Println("err: ", err)
		}
		err := WsTickerServe(symbol, wsDepthHandler, errHandler)
		if err != nil {
			fmt.Println("err: ", err)
			return
		}
	}()

	fmt.Println("End of ShowBinancePrice ")
	return outputChannel
}
