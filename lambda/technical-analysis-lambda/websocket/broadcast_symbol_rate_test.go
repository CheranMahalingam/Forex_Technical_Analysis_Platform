package websocket

import (
	"technical-analysis-lambda/finance"
	"testing"
	"time"
)

func TestCreateCallbackSymbolMessage(t *testing.T) {
	// Test validity of fields in payload
	financeItems := []finance.FinancialDataItem{}
	for i := 1; i < 5; i++ {
		currentDate := time.Now().Add(-time.Hour * 24 * time.Duration(i))
		eurusd := finance.SymbolData{Open: float32(i), High: float32(i), Low: float32(i), Close: float32(i), Volume: float32(i)}
		item := finance.FinancialDataItem{
			Date:      currentDate.Format("2006-01-02"),
			Timestamp: currentDate.Format("15:04:05"),
			EURUSD:    eurusd,
		}
		financeItems = append(financeItems, item)
	}

	payload, err := createCallbackSymbolMessage(&financeItems, "EURUSD")
	if err != nil {
		t.Errorf("Callback function threw error: %e", err)
	}
	for i := 1; i < 5; i++ {
		currentDate := time.Now().Add(-time.Hour * 24 * time.Duration(i))
		eurusd := finance.SymbolData{Open: float32(i), High: float32(i), Low: float32(i), Close: float32(i), Volume: float32(i)}
		message := (*payload)[i-1]

		if message.Timestamp != currentDate.Format("2006-01-02 15:04:05") {
			t.Errorf("Payload timestamp invalid, got: %s, expected: %s", message.Timestamp, currentDate.Format("2006-01-02 15:04:05"))
		}
		if message.Open != eurusd.Open {
			t.Errorf("Payload open price invalid, got: %f, expected: %f", message.Open, eurusd.Open)
		}
		if message.High != eurusd.High {
			t.Errorf("Payload high price invalid, got: %f, expected: %f", message.High, eurusd.High)
		}
		if message.Low != eurusd.Low {
			t.Errorf("Payload low price invalid, got: %f, expected: %f", message.Low, eurusd.Low)
		}
		if message.Close != eurusd.Close {
			t.Errorf("Payload close price invalid, got: %f, expected: %f", message.Close, eurusd.Close)
		}
		if message.Volume != eurusd.Volume {
			t.Errorf("Payload trading volume invalid, got: %f, expected: %f", message.Volume, eurusd.Volume)
		}
	}
}

func TestCheckIsRateValid(t *testing.T) {
	// Test response when empty struct is passed
	emptySymbolField := finance.SymbolData{}
	isValid := checkIsRateValid(&emptySymbolField)
	if isValid != false {
		t.Errorf("Invalid exchange rate, got: %t", isValid)
	}

	// Test response when one field is zero
	symbolField := finance.SymbolData{Open: 0, High: 1.1, Low: 2.1, Close: 3.1, Volume: 5}
	isValid = checkIsRateValid(&symbolField)
	if isValid != false {
		t.Errorf("Invalid exchange rate, got: %t", isValid)
	}

	// Test response when all fields are valid
	symbolField = finance.SymbolData{Open: 2.1, High: 1.1, Low: 2.1, Close: 3.1, Volume: 5}
	isValid = checkIsRateValid(&symbolField)
	if isValid != true {
		t.Errorf("Exchange rate was valid, got: %t", isValid)
	}
}

func TestGetSymbolRateField(t *testing.T) {
	// Test whether correct symbol field is retrieved
	financialData := finance.FinancialDataItem{}
	eurusd := finance.SymbolData{Open: 0, High: 1, Low: 1.2, Close: 5.1, Volume: 70}
	financialData.EURUSD = eurusd
	symbolField := getSymbolRateField("EURUSD", &financialData)
	if symbolField.Open != eurusd.Open {
		t.Errorf("Invalid open price, got: %f, expected: %f", symbolField.Open, eurusd.Open)
	}
	if symbolField.High != eurusd.High {
		t.Errorf("Invalid high price, got: %f, expected: %f", symbolField.High, eurusd.High)
	}
	if symbolField.Low != eurusd.Low {
		t.Errorf("Invalid Low price, got: %f, expected: %f", symbolField.Low, eurusd.Low)
	}
	if symbolField.Close != eurusd.Close {
		t.Errorf("Invalid close price, got: %f, expected: %f", symbolField.Close, eurusd.Close)
	}
	if symbolField.Volume != eurusd.Volume {
		t.Errorf("Invalid trade volume, got: %f, expected: %f", symbolField.Volume, eurusd.Volume)
	}

	// Test graceful failure when invalid symbol is given
	symbolField = getSymbolRateField("eUrUsD", &financialData)
	if symbolField != nil {
		t.Errorf("Invalid field, got: %v, expected: nil", symbolField)
	}
}
