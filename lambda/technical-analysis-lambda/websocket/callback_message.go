package websocket

type CallbackMessageData struct {
	Symbol    string
	Timestamp string
	Open      float32
	High      float32
	Low       float32
	Close     float32
	Volume    float32
}
