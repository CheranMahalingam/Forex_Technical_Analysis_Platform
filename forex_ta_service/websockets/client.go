package websockets

import (
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 128
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  128,
	WriteBufferSize: 128,
	// We do not allow all origins to pass
	CheckOrigin: func(r *http.Request) bool {
		origin := r.Header.Get("Origin")
		if origin == "" {
			return true
		}
		originURL, err := url.Parse(origin)
		// Only requests from localhost:3000 should be accepted
		return err == nil && (strings.HasPrefix(originURL.Host, "localhost:3000"))
	},
}

type Client struct {
	pool *[]Pool
	conn *websocket.Conn
	send chan []byte
}

func newClient(conn *websocket.Conn) *Client {
	emptySlice := []Pool{}
	return &Client{pool: &emptySlice, conn: conn, send: make(chan []byte)}
}

func (c *Client) readPump() {
	defer func() {
		for i := range *c.pool {
			pool := (*c.pool)[i]
			pool.unregister <- c
		}
		c.conn.Close()
		log.Println("readPump: connection closed")
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline((time.Now().Add(pongWait)))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		var message SubscriptionMessage
		err := c.conn.ReadJSON(&message)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				unsubscribeFromAllPools(c)
				log.Printf("error: %v\n", err)
			}
			break
		}
		log.Println(message)
		log.Println("readPump: message broadcasted")
		if message.MessageType == "subscribe" {
			subscribeToPool(message.Symbol, c)
		} else if message.MessageType == "unsubscribe" {
			unsubscribeFromPool(message.Symbol, c)
		}
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		for i := range *c.pool {
			pool := (*c.pool)[i]
			pool.unregister <- c
		}
		c.conn.Close()
		log.Println("writePump: connection closed")
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				log.Println("writePump: close message")
				return
			}
			err := c.conn.WriteMessage(websocket.TextMessage, message)
			if err != nil {
				log.Println(err)
				return
			}
			log.Println("writePump: messages written")
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Println(err)
				return
			}
			log.Println("writePump: timer reset")
		}
	}
}

func ServeWs(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	client := newClient(conn)
	log.Println("ServeWs: new client registered", client)

	go client.writePump()
	go client.readPump()
}
