package exchangerate

import (
	"context"
	"net/http"
	"technical-analysis-lambda/finance"
	"testing"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

var currencyPairs = [4]string{"EURUSD", "GBPUSD", "USDJPY", "AUDCAD"}

type forexApiMock struct {
	mockForexCandles func(ctx context.Context, symbol string, resolution string, from int64, to int64) (finnhub.ForexCandles, *http.Response, error)
	mockGeneralNews  func(ctx context.Context, category string, localVarOptionals *finnhub.GeneralNewsOpts) ([]finnhub.News, *http.Response, error)
}

func (fam *forexApiMock) ForexCandles(ctx context.Context, symbol string, resolution string, from int64, to int64) (finnhub.ForexCandles, *http.Response, error) {
	if fam.mockForexCandles != nil {
		return fam.mockForexCandles(ctx, symbol, resolution, from, to)
	}
	return finnhub.ForexCandles{}, nil, nil
}

func (fam *forexApiMock) GeneralNews(ctx context.Context, category string, localVarOptionals *finnhub.GeneralNewsOpts) ([]finnhub.News, *http.Response, error) {
	if fam.mockGeneralNews != nil {
		return fam.mockGeneralNews(ctx, category, localVarOptionals)
	}
	return nil, nil, nil
}

func TestCreateNewSymbolRate(t *testing.T) {
	// Test response when Finnhub client responds with no data
	mock := &forexApiMock{
		mockForexCandles: func(ctx context.Context, symbol string, resolution string, from int64, to int64) (finnhub.ForexCandles, *http.Response, error) {
			return finnhub.ForexCandles{S: "no_data"}, nil, nil
		},
	}
	headline := "dummy value"

	ohlc, news, err := CreateNewSymbolRate(mock, &currencyPairs, 100, 101, "1", &headline)
	if ohlc != nil {
		t.Errorf("OHLC data retrieved was invalid, got: %v, expected: nil", *ohlc)
	}
	if news != nil {
		t.Errorf("News data was invalid, got: %v, expected: nil", *news)
	}
	if err != nil {
		t.Errorf("Received unexpected error: %e", err)
	}

	// Test response when Finnhub client provides a single invalid message
	// Invalid due to timestamp being out of specified range
	finnhubData := finnhub.ForexCandles{
		O: []float32{1.1},
		H: []float32{1.3},
		L: []float32{1.01},
		C: []float32{1.02},
		V: []float32{100},
		T: []float32{102},
		S: "ok",
	}
	mock = &forexApiMock{
		mockForexCandles: func(ctx context.Context, symbol string, resolution string, from int64, to int64) (finnhub.ForexCandles, *http.Response, error) {
			return finnhubData, nil, nil
		},
	}

	ohlc, news, err = CreateNewSymbolRate(mock, &currencyPairs, 100, 101, "1", &headline)
	if len(*ohlc) != 0 {
		t.Errorf("OHLC data retrieved was invalid, got length: %d, expected length: 0", len(*ohlc))
	}
	if len(*news) != 0 {
		t.Errorf("News data was invalid, got length: %d, expected length: 0", len(*news))
	}
	if err != nil {
		t.Errorf("Received unexpected error: %e", err)
	}

	// Test response when Finnhub client provides a single valid message
	currentTime := time.Now().Unix()
	finnhubData.T[0] = float32(currentTime - 700)

	ohlc, news, err = CreateNewSymbolRate(mock, &currencyPairs, 1000, 0, "1", &headline)
	if (*ohlc)[0].EURUSD.Open != finnhubData.O[0] {
		t.Errorf("OHLC data retrieved was invalid, got open price: %f, expected: %f", (*ohlc)[0].EURUSD.Open, finnhubData.O[0])
	}
	if (*ohlc)[0].EURUSD.High != finnhubData.H[0] {
		t.Errorf("OHLC data retrieved was invalid, got high price: %f, expected: %f", (*ohlc)[0].EURUSD.High, finnhubData.H[0])
	}
	if (*ohlc)[0].EURUSD.Low != finnhubData.L[0] {
		t.Errorf("OHLC data retrieved was invalid, got low price: %f, expected: %f", (*ohlc)[0].EURUSD.Low, finnhubData.L[0])
	}
	if (*ohlc)[0].EURUSD.Close != finnhubData.C[0] {
		t.Errorf("OHLC data retrieved was invalid, got close price: %f, expected: %f", (*ohlc)[0].EURUSD.Close, finnhubData.C[0])
	}
	if (*ohlc)[0].EURUSD.Volume != finnhubData.V[0] {
		t.Errorf("OHLC data retrieved was invalid, got trade volume: %f, expected: %f", (*ohlc)[0].EURUSD.Volume, finnhubData.V[0])
	}

	if len(*news) != 0 {
		t.Errorf("News data was invalid, got length: %d, expected length: 0", len(*news))
	}
	if err != nil {
		t.Errorf("Received unexpected error: %e", err)
	}

	// Test response when duplicate data is given
	finnhubData.O = append(finnhubData.O, 1.2)
	finnhubData.H = append(finnhubData.H, 2.1)
	finnhubData.L = append(finnhubData.L, 0.98)
	finnhubData.C = append(finnhubData.C, 1.1)
	finnhubData.T = append(finnhubData.T, float32(currentTime-700))
	finnhubData.V = append(finnhubData.V, 192)

	ohlc, news, err = CreateNewSymbolRate(mock, &currencyPairs, 1000, 0, "1", &headline)
	// The duplicate data should replace the original data
	// This was a design choice to deal with bugs in Finnhub where duplicate data is given
	// with greater trade volume
	if (*ohlc)[0].EURUSD.Open != finnhubData.O[1] {
		t.Errorf("OHLC data retrieved was invalid, got open price: %f, expected: %f", (*ohlc)[0].EURUSD.Open, finnhubData.O[0])
	}
	if (*ohlc)[0].EURUSD.High != finnhubData.H[1] {
		t.Errorf("OHLC data retrieved was invalid, got high price: %f, expected: %f", (*ohlc)[0].EURUSD.High, finnhubData.H[0])
	}
	if (*ohlc)[0].EURUSD.Low != finnhubData.L[1] {
		t.Errorf("OHLC data retrieved was invalid, got low price: %f, expected: %f", (*ohlc)[0].EURUSD.Low, finnhubData.L[0])
	}
	if (*ohlc)[0].EURUSD.Close != finnhubData.C[1] {
		t.Errorf("OHLC data retrieved was invalid, got close price: %f, expected: %f", (*ohlc)[0].EURUSD.Close, finnhubData.C[0])
	}
	if (*ohlc)[0].EURUSD.Volume != finnhubData.V[1] {
		t.Errorf("OHLC data retrieved was invalid, got trade volume: %f, expected: %f", (*ohlc)[0].EURUSD.Volume, finnhubData.V[0])
	}

	if len(*news) != 0 {
		t.Errorf("News data was invalid, got length: %d, expected length: 0", len(*news))
	}
	if err != nil {
		t.Errorf("Received unexpected error: %e", err)
	}

	// Test response when multiple messages are provided with no duplicates
	finnhubData.T[1] = float32(currentTime - 6000)
	finnhubData.O = append(finnhubData.O, 1.22)
	finnhubData.H = append(finnhubData.H, 2.12)
	finnhubData.L = append(finnhubData.L, 0.983)
	finnhubData.C = append(finnhubData.C, 1.11)
	finnhubData.T = append(finnhubData.T, float32(currentTime-8000))
	finnhubData.V = append(finnhubData.V, 1920)

	ohlc, news, err = CreateNewSymbolRate(mock, &currencyPairs, 10000, 0, "1", &headline)
	for i := 0; i < 3; i++ {
		if (*ohlc)[i].EURUSD.Open != finnhubData.O[i] {
			t.Errorf("OHLC data retrieved was invalid, got open price: %f, expected: %f", (*ohlc)[i].EURUSD.Open, finnhubData.O[i])
		}
		if (*ohlc)[i].EURUSD.High != finnhubData.H[i] {
			t.Errorf("OHLC data retrieved was invalid, got high price: %f, expected: %f", (*ohlc)[i].EURUSD.High, finnhubData.H[i])
		}
		if (*ohlc)[i].EURUSD.Low != finnhubData.L[i] {
			t.Errorf("OHLC data retrieved was invalid, got low price: %f, expected: %f", (*ohlc)[i].EURUSD.Low, finnhubData.L[i])
		}
		if (*ohlc)[i].EURUSD.Close != finnhubData.C[i] {
			t.Errorf("OHLC data retrieved was invalid, got close price: %f, expected: %f", (*ohlc)[i].EURUSD.Close, finnhubData.C[i])
		}
		if (*ohlc)[i].EURUSD.Volume != finnhubData.V[i] {
			t.Errorf("OHLC data retrieved was invalid, got trade volume: %f, expected: %f", (*ohlc)[i].EURUSD.Volume, finnhubData.V[i])
		}
	}

	if len(*news) != 0 {
		t.Errorf("News data was invalid, got length: %d, expected length: 0", len(*news))
	}
	if err != nil {
		t.Errorf("Received unexpected error: %e", err)
	}
}

func TestSearchTimeIndex(t *testing.T) {
	// Test no matches found
	symbolData := []finance.FinancialDataItem{}
	for i := 1; i < 5; i++ {
		currentDate := time.Now().Add(24 * time.Hour * time.Duration(i))
		newData := finance.FinancialDataItem{
			Date:      currentDate.Format("2006-01-02"),
			Timestamp: currentDate.Format("15:04:05"),
		}
		symbolData = append(symbolData, newData)
	}

	currentDate := time.Now()
	index := searchTimeIndex(&symbolData, currentDate.Format("2006-01-02"), currentDate.Format("15:04:05"))
	if index != -1 {
		t.Errorf("Invalid index, got: %d, expected: -1", index)
	}

	// Test date match without timstamp match
	futureDate := time.Now().Add(25 * time.Hour)
	index = searchTimeIndex(&symbolData, futureDate.Format("2006-01-02"), futureDate.Format("15:04:05"))
	if index != -1 {
		t.Errorf("Invalid index, got: %d, expected: -1", index)
	}

	// Test successful match
	futureDate = time.Now().Add(48 * time.Hour)
	index = searchTimeIndex(&symbolData, futureDate.Format("2006-01-02"), futureDate.Format("15:04:05"))
	if index != 1 {
		t.Errorf("Invalid index, got: %d, expected: 1", index)
	}
}
