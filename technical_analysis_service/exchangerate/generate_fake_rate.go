package exchangerate

import (
	"math/rand"
	"time"
)

// Used to generate data for weekends to simulate live trading
// Useful to avoid being rate limited by finnhub api
func GenerateFakeData(prevRate *ExchangeRate, symbol string, startSeconds int64, endSeconds int64) map[string][]ExchangeRate {
	var newExchangeRateData []ExchangeRate
	dataLength := int((startSeconds - endSeconds) / 60)
	for i := 0; i < dataLength; i++ {
		if prevRate.Open == 0 {
			prevRate.Open = (rand.Float32()*2-1)*0.001 + 1
		} else {
			prevRate.Open = (rand.Float32()*2-1)*0.001 + prevRate.Close
		}
		prevRate.Close = (rand.Float32()*2-1)*0.001 + prevRate.Open
		if prevRate.Open > prevRate.Close {
			prevRate.High = rand.Float32()*0.001 + prevRate.Open
			prevRate.Low = prevRate.Close - rand.Float32()*0.001
		} else {
			prevRate.High = rand.Float32()*0.001 + prevRate.Close
			prevRate.Low = prevRate.Open - rand.Float32()*0.001
		}
		prevRate.Timestamp = time.Unix(time.Now().Unix()-startSeconds+int64((i+1)*60), 0).UTC()
		prevRate.Volume = rand.Float32() * 10
		newExchangeRateData = append(newExchangeRateData, *prevRate)
	}
	symbolData := map[string][]ExchangeRate{
		symbol: newExchangeRateData,
	}
	return symbolData
}
