package broadcast

import (
	"testing"
	"time"
)

func TestCreateCallbackSymbolMessage(t *testing.T) {
	// Test validity of fields in symbol payload for a single message
	currentDate := time.Now()
	newSymbolDate := []exchangeRateTable{}
	newOhlcData := ohlcData{Open: 1, High: 5.5, Low: 3.1, Close: 1.0, Volume: 5}
	newSymbolDate = append(newSymbolDate, exchangeRateTable{Date: currentDate.Format("2006-01-02"), Timestamp: currentDate.Format("15:04:05"), EURUSD: newOhlcData})

	exchangeRatePayload := createCallbackSymbolMessage(&newSymbolDate, "EURUSD")
	for _, message := range *exchangeRatePayload {
		timestamp := currentDate.Format("2006-01-02") + " " + currentDate.Format("15:04:05")
		if message.Timestamp != timestamp {
			t.Errorf("Invalid symbol message payload timestamp, got: %s, expected: %s", message.Timestamp, timestamp)
		}
		if message.Open != newOhlcData.Open {
			t.Errorf("Invalid symbol message payload open price, got: %f, expected: %f", message.Open, newOhlcData.Open)
		}
		if message.High != newOhlcData.High {
			t.Errorf("Invalid symbol message payload high price, got: %f, expected: %f", message.High, newOhlcData.High)
		}
		if message.Close != newOhlcData.Close {
			t.Errorf("Invalid symbol message payload close price, got: %f, expected: %f", message.Close, newOhlcData.Close)
		}
		if message.Low != newOhlcData.Low {
			t.Errorf("Invalid symbol message payload low price, got: %f, expected: %f", message.Low, newOhlcData.Low)
		}
		if message.Volume != newOhlcData.Volume {
			t.Errorf("Invalid symbol message payload trade volume, got: %f, expected: %f", message.Volume, newOhlcData.Volume)
		}
	}

	// Test validity of fields in symbol payload for multiple messages
	for i := 0; i < 3; i++ {
		currentDate = time.Now().Add(-24 * time.Hour * time.Duration(i+1))
		newOhlcData = ohlcData{
			Open:   float32(i+1) * 0.1,
			High:   float32(i+1) * 0.1,
			Low:    float32(i+1) * 0.1,
			Close:  float32(i+1) * 0.1,
			Volume: float32(i+1) * 0.1,
		}
		newSymbolDate = append(newSymbolDate, exchangeRateTable{Date: currentDate.Format("2006-01-02"), Timestamp: currentDate.Format("15:04:05"), EURUSD: newOhlcData})
	}
	exchangeRatePayload = createCallbackSymbolMessage(&newSymbolDate, "EURUSD")
	for index, message := range *exchangeRatePayload {
		timestamp := newSymbolDate[index].Date + " " + newSymbolDate[index].Timestamp
		if message.Timestamp != timestamp {
			t.Errorf("Invalid symbol message payload timestamp, got: %s, expected: %s", message.Timestamp, timestamp)
		}
		if message.Open != newSymbolDate[index].EURUSD.Open {
			t.Errorf("Invalid symbol message payload open price, got: %f, expected: %f", message.Open, newOhlcData.Open)
		}
		if message.High != newSymbolDate[index].EURUSD.High {
			t.Errorf("Invalid symbol message payload high price, got: %f, expected: %f", message.High, newOhlcData.High)
		}
		if message.Close != newSymbolDate[index].EURUSD.Close {
			t.Errorf("Invalid symbol message payload close price, got: %f, expected: %f", message.Close, newOhlcData.Close)
		}
		if message.Low != newSymbolDate[index].EURUSD.Low {
			t.Errorf("Invalid symbol message payload low price, got: %f, expected: %f", message.Low, newOhlcData.Low)
		}
		if message.Volume != newSymbolDate[index].EURUSD.Volume {
			t.Errorf("Invalid symbol message payload trade volume, got: %f, expected: %f", message.Volume, newOhlcData.Volume)
		}
	}

	// Test failure when invalid ohlc data is passed
	currentDate = time.Now()
	newSymbolDate = []exchangeRateTable{}
	newOhlcData = ohlcData{Open: 1, High: 0, Low: 3.1, Close: 1.0, Volume: 5}
	newSymbolDate = append(newSymbolDate, exchangeRateTable{Date: currentDate.Format("2006-01-02"), Timestamp: currentDate.Format("15:04:05"), EURUSD: newOhlcData})

	exchangeRatePayload = createCallbackSymbolMessage(&newSymbolDate, "EURUSD")
	if len(*exchangeRatePayload) != 0 {
		t.Errorf("Invalid symbol message payload length, got: %d, expected: 0", len(*exchangeRatePayload))
	}
}

func TestCreateCallbackInferenceMessage(t *testing.T) {
	// Test validity of fields in inference websocket payload
	currentDate := time.Now()
	inferencePayload := []inferenceTable{}
	prediction := []float32{1.0, 1.1, 1.3, 1.4}
	newInferenceTable := inferenceTable{Latest: currentDate.Format("2006-01-02"), EURUSDInference: prediction}
	inferencePayload = append(inferencePayload, newInferenceTable)

	inferenceMessage := createCallbackInferenceMessage(&inferencePayload, "EURUSD")
	if inferenceMessage.Date != currentDate.Format("2006-01-02") {
		t.Errorf("Inference message payload date was invalid, got: %s, expected: %s", inferenceMessage.Date, currentDate.Format("2006-01-02"))
	}
	// Validate that predictions have not been tampered
	for index, pred := range inferenceMessage.Inference {
		if pred != prediction[index] {
			t.Errorf("Inference message payload inference was invalid, got: %f, expected: %f", pred, prediction[index])
		}
	}
}
