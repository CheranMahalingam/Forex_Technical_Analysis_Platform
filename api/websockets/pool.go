package websockets

import (
	"encoding/json"
	"log"
	"time"

	"github.com/forexapi/technicalanalysis/controllers"
)

type Pool struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	pair       string
}

func NewPool(currencyPair string) *Pool {
	return &Pool{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (p *Pool) Run() {
	for {
		select {
		case client := <-p.register:
			p.clients[client] = true
		case client := <-p.unregister:
			if _, ok := p.clients[client]; ok {
				delete(p.clients, client)
				close(client.send)
				log.Println("client unregistered")
			}
		case message := <-p.broadcast:
			for client := range p.clients {
				select {
				case client.send <- message:
				default:
					delete(p.clients, client)
					close(client.send)
				}
			}
		}
	}
}

func (p *Pool) ExchangeRateOnInterval(intervalSeconds time.Duration) {
	ticker := time.NewTicker(intervalSeconds)
	defer func() {
		ticker.Stop()
	}()
	for {
		select {
		case <-ticker.C:
			newRate := controllers.GetLatestRate("EUR_USD", 60, "1")
			newRateJSON, err := json.Marshal(*newRate)
			if err != nil {
				log.Println(err)
				return
			}
			p.broadcast <- newRateJSON
		}
	}
}
