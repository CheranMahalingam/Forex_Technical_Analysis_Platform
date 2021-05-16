package subscription

type CallbackMessageData struct {
	Timestamp string  `json:"timestamp"`
	Open      float32 `json:"open"`
	High      float32 `json:"high"`
	Low       float32 `json:"low"`
	Close     float32 `json:"close"`
	Volume    float32 `json:"volume"`
}

type CallbackMessageInference struct {
	Inference []float32 `json:"inference"`
	Date      string    `json:"date"`
}

type CallbackMessageNews struct {
	MarketNews newsItem
}
