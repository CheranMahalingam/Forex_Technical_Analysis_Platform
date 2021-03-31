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
				client.send <- message
				/*default:
				delete(p.clients, client)
				close(client.send)*/
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
			newRate := exchangerate.GetLatestRate(p.currencyPair, p.previousRate, 60, 0, "1")
			log.Println(p.previousRate, "Exchange Rate: latest", p.currencyPair)
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
