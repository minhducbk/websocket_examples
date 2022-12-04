package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/minhducbk/websocket_examples/ws_deribit/deribit"
	"github.com/minhducbk/websocket_examples/ws_deribit/services"
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
	deribitClient := deribit.SetupClient()
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		services.ServeWs(hub, w, r, deribitClient)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
