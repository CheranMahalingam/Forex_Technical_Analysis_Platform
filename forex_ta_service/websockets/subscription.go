package websockets

import (
	"encoding/json"
	"log"
	"time"

	"github.com/forexapi/technicalanalysis/exchangerate"
)

type SubscriptionMessage struct {
	MessageType string `json:"messageType"`
	Symbol      string `json:"symbol"`
}

var currencyPairs = [4]string{"EURUSD", "GBPUSD", "USDCAD", "AUDUSD"}
var currencyPoolMap = make(map[string]*Pool)

func SetupCurrencyPools(interval time.Duration) {
	for _, currency := range currencyPairs {
		currencyPoolMap[currency] = newPool(currency)
		go currencyPoolMap[currency].Run()
		go currencyPoolMap[currency].ExchangeRateOnInterval(interval)
	}
	log.Println("Setup complete!")
}

func subscribeToPool(pair string, client *Client) {
	currencyPool := currencyPoolMap[pair]
	*(client.pool) = append(*(client.pool), *currencyPool)
	log.Println(client.pool)

	currencyPool.register <- client

	// When a new client subscribes to a pool they should be caught up with past data
	newRate := exchangerate.GetLatestRate(pair, 86400, "1")
	newRateJSON, err := json.Marshal(*newRate)
	if err != nil {
		log.Println(err)
		return
	}
	client.send <- newRateJSON
	log.Printf("Client registered to %s pool\n", pair)
}

func unsubscribeFromPool(pair string, client *Client) {
	currencyPool := currencyPoolMap[pair]
	for index, pool := range *(client.pool) {
		if pool.currencyPair == currencyPool.currencyPair {
			poolLength := len(*(client.pool))
			(*(client.pool))[index] = (*(client.pool))[poolLength-1]
			(*(client.pool)) = (*(client.pool))[:poolLength-1]
		}
	}
	currencyPool.unregister <- client
	log.Printf("Client unregistered from %s pool\n", pair)
}
