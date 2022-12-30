package main

import (
	"encoding/json"
	"flag"
	"log"
	"math/rand"
	"net/http"
	_ "net/http/pprof"
	"time"

	"github.com/minhducbk/websocket_examples/ws/binance"
	"github.com/minhducbk/websocket_examples/ws/services"
)

var addr = flag.String("addr", ":8084", "http service address")

func serveHome(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "./chat/services/home.html")
}

func main() {
	flag.Parse()
	hub := services.NewHub()
	go hub.Run()

	// deribitClient := deribit.SetupClient()
	// deribitClient.SubscribeCmd()
	// go deribitClient.FlushPricesIntoChannelCmd()
	// wss://nbstream.binance.com/eoptions/ws/ETH-230127-1300-C@ticker
	middleChan := binance.ShowBinancePrice("ETH-230127-1300-C")

	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		services.ServeWs(hub, w, r, middleChan)
	})
	http.HandleFunc("/log", logHandler)
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

func logHandler(w http.ResponseWriter, r *http.Request) {
	ch := make(chan int)
	go func() {
		obj := make(map[string]float64)
		if err := json.NewDecoder(r.Body).Decode(&obj); err != nil {
			ch <- http.StatusBadRequest
			return
		}
		// simulation of a time consuming process like writing logs into db
		time.Sleep(time.Duration(rand.Intn(400)) * time.Millisecond)
		ch <- http.StatusOK
	}()

	select {
	case status := <-ch:
		w.WriteHeader(status)
	case <-time.After(200 * time.Millisecond):
		w.WriteHeader(http.StatusRequestTimeout)
	}
}
