package exchangerate

import (
	"math/rand"
	"time"
)

// Used to generate data for weekends to simulate live trading
// Useful to avoid being rate limited by finnhub api
func GenerateFakeData(symbol string, seconds int64) *map[string][]ExchangeRate {
	var newExchangeRateData []ExchangeRate
	dataLength := int(seconds / 60)
	var timestamp time.Time
	var open float32
	var close float32
	var high float32
	var low float32
	var volume float32
	for i := 0; i < dataLength; i++ {
		if i == 0 {
			open = (rand.Float32()*2-1)*0.001 + 1
		} else {
			open = (rand.Float32()*2-1)*0.001 + close
		}
		close = (rand.Float32()*2-1)*0.001 + open
		if open > close {
			high = rand.Float32()*0.001 + open
			low = close - rand.Float32()*0.001
		} else {
			high = rand.Float32()*0.001 + close
			low = open - rand.Float32()*0.001
		}
		timestamp = time.Unix(time.Now().Unix()-seconds+int64((i+1)*60), 0).UTC()
		volume = rand.Float32() * 100
		newExchangeRateData = append(newExchangeRateData, ExchangeRate{Timestamp: timestamp, Open: open, High: high, Low: low, Close: close, Volume: volume})
	}
	symbolData := map[string][]ExchangeRate{
		symbol: newExchangeRateData,
	}
	return &symbolData
}
