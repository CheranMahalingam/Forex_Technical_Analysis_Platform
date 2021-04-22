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
	r := mux.NewRouter()
	r.HandleFunc("/ws", websockets.ServeWs)

	websockets.SetupCurrencyPools(120 * time.Second)

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
