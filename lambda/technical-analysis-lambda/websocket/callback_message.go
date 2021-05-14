package websocket

type CallbackSymbolMessage struct {
	Timestamp string  `json:"timestamp"`
	Open      float32 `json:"open"`
	High      float32 `json:"high"`
	Low       float32 `json:"low"`
	Close     float32 `json:"close"`
	Volume    float32 `json:"volume"`
}

type CallbackNewsMessage struct {
	Timestamp string `json:"timestamp"`
	Headline  string `json:"headline"`
	Image     string `json:"image"`
	Source    string `json:"source"`
	Summary   string `json:"summary"`
	NewsUrl   string `json:"newsurl"`
}
