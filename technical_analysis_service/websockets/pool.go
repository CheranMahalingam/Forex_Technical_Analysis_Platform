package websockets

import (
	"encoding/json"
	"log"
	"time"

	"github.com/forexapi/technicalanalysis/exchangerate"
)

type Pool struct {
	clients      map[*Client]bool
	broadcast    chan []byte
	register     chan *Client
	unregister   chan *Client
	currencyPair string
	previousRate *exchangerate.ExchangeRate
}

var lmao = [2]string{"EURUSD", "GBPUSD"}

func newPool(currencyPair string) *Pool {
	return &Pool{
		clients:      make(map[*Client]bool),
		broadcast:    make(chan []byte),
		register:     make(chan *Client),
		unregister:   make(chan *Client),
		currencyPair: currencyPair,
		previousRate: &exchangerate.ExchangeRate{},
	}
}

func (p *Pool) Run() {
	for {
		select {
		case client := <-p.register:
			p.clients[client] = true
			log.Println("client registered")
		case client := <-p.unregister:
			if _, ok := p.clients[client]; ok {
				delete(p.clients, client)
				//close(client.send)
				log.Println("client unregistered")
			}
		case message := <-p.broadcast:
			for client := range p.clients {
				select {
				case client.send <- message:
				default:
					delete(p.clients, client)
				}
			}
		}
	}
}

func (p *Pool) ExchangeRateOnInterval(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			newRate := exchangerate.GetLatestRate(p.currencyPair, p.previousRate, 120, 0, "1")
			log.Println(*exchangerate.CreateNewSymbolRate(&lmao, 120, 0, "1"))
			log.Println("Exchange Rate: latest", newRate)
			if newRate != nil {
				newRateJSON, err := json.Marshal(newRate)
				if err != nil {
					log.Println(err)
					return
				}
				p.broadcast <- newRateJSON
			}
		}
	}
}
