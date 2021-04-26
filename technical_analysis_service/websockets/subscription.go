package websockets

import (
	"encoding/json"
	"log"
	"time"

	"github.com/forexapi/technicalanalysis/exchangerate"
)

type SubscriptionMessage struct {
	MessageType string `json:"message"`
	Symbol      string `json:"data"`
}

var currencyPairs = [2]string{"EURUSD", "GBPUSD"}
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
	for i := 1; i > 0; i-- {
		newRate := exchangerate.GetLatestRate(pair, currencyPool.previousRate, 86400*int64(i)+48*60*60, 86300*int64(i)+48*60*60, "1")
		if newRate == nil {
			log.Println("Subscription: no data found")
			continue
		}
		newRateJSON, err := json.Marshal(newRate)
		if err != nil {
			log.Println(err)
			return
		}
		client.send <- newRateJSON
	}
	log.Printf("Client registered to %s pool\n", pair)
}

func unsubscribeFromPool(pair string, client *Client) {
	currencyPool := currencyPoolMap[pair]
	for index, pool := range *(client.pool) {
		if pool.currencyPair == currencyPool.currencyPair {
			poolLength := len(*(client.pool))
			(*(client.pool))[index] = (*(client.pool))[poolLength-1]
			(*(client.pool)) = (*(client.pool))[:poolLength-1]
			log.Printf("Client unregistered from %s pool\n", pair)
		}
	}
	currencyPool.unregister <- client
}

func unsubscribeFromAllPools(client *Client) {
	for _, pair := range currencyPairs {
		unsubscribeFromPool(pair, client)
	}
}
