package exchangerate

import (
	"time"
)

type ExchangeRate struct {
	Timestamp time.Time `json:"timestamp"`
	Open      float32   `json:"open"`
	High      float32   `json:"high"`
	Low       float32   `json:"low"`
	Close     float32   `json:"close"`
	Volume    float32   `json:"volume"`
}
