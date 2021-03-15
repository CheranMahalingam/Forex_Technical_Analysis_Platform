package main

import (
	"flag"
	"log"
	"net/http"

	"github.com/forexapi/technicalanalysis/websockets"
	"github.com/gorilla/mux"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	pool := websockets.NewPool()
	go pool.Run()
	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWs(pool, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal(err)
	}
}
