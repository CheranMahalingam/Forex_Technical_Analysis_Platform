package main

import (
	"log"
	"net/http"
	"time"

	"github.com/forexapi/technicalanalysis/websockets"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"github.com/rs/cors"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	pool := websockets.NewPool("EUR_USD")
	go pool.Run()
	go pool.ExchangeRateOnInterval(10 * time.Second)
	r := mux.NewRouter()
	r.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		websockets.ServeWs(pool, w, r)
	})

	corsOpts := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000/"},
		AllowedMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
			http.MethodOptions,
			http.MethodHead,
		},
		AllowedHeaders: []string{
			"Accept",
			"Content-Type",
			"Origin",
		},
	})

	err := http.ListenAndServe(":8080", corsOpts.Handler(r))
	if err != nil {
		log.Fatal(err)
	}
}
