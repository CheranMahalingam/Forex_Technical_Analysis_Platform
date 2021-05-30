package exchangerate

import (
	"context"
	"net/http"
	"os"
	"testing"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

type newsApiMock struct {
	mockForexCandles func(ctx context.Context, symbol string, resolution string, from int64, to int64) (finnhub.ForexCandles, *http.Response, error)
	mockGeneralNews  func(ctx context.Context, category string, localVarOptionals *finnhub.GeneralNewsOpts) ([]finnhub.News, *http.Response, error)
}

func (nam *newsApiMock) ForexCandles(ctx context.Context, symbol string, resolution string, from int64, to int64) (finnhub.ForexCandles, *http.Response, error) {
	if nam.mockForexCandles != nil {
		return nam.mockForexCandles(ctx, symbol, resolution, from, to)
	}
	return finnhub.ForexCandles{}, nil, nil
}

func (nam *newsApiMock) GeneralNews(ctx context.Context, category string, localVarOptionals *finnhub.GeneralNewsOpts) ([]finnhub.News, *http.Response, error) {
	if nam.mockGeneralNews != nil {
		return nam.mockGeneralNews(ctx, category, localVarOptionals)
	}
	return nil, nil, nil
}

func TestGetMarketNews(t *testing.T) {
	// Test response when Finnhub client returns no market news
	mock := &newsApiMock{
		mockGeneralNews: func(ctx context.Context, category string, localVarOptionals *finnhub.GeneralNewsOpts) ([]finnhub.News, *http.Response, error) {
			return []finnhub.News{}, nil, nil
		},
	}
	auth := context.WithValue(context.Background(), finnhub.ContextAPIKey, finnhub.APIKey{
		Key: os.Getenv("PROJECT_API_KEY"),
	})
	headline := "words"

	news, err := getMarketNews(mock, &auth, &headline)
	if err != nil {
		t.Errorf("Unexpected error: %e", err)
	}
	if len(*news) != 0 {
		t.Errorf("Invalid market news, got length: %d, expected: 0", len(*news))
	}

	// Test response with duplicate headline but single message
	newNews := []finnhub.News{}
	article := finnhub.News{Headline: "words"}
	newNews = append(newNews, article)
	mock = &newsApiMock{
		mockGeneralNews: func(ctx context.Context, category string, localVarOptionals *finnhub.GeneralNewsOpts) ([]finnhub.News, *http.Response, error) {
			return newNews, nil, nil
		},
	}

	// Since new headline matches passed headline the function should not return
	// the duplicate news, however, since it is a single message it will be returned
	news, err = getMarketNews(mock, &auth, &headline)
	if err != nil {
		t.Errorf("Unexpected error: %e", err)
	}
	if (*news)[0].Headline != "words" {
		t.Errorf("Invalid market news, got headline: %s, expected headline: words", (*news)[0].Headline)
	}

	// Test response with duplicate headline but multiple messages
	newNews = []finnhub.News{}
	article1 := finnhub.News{Headline: "word"}
	article2 := finnhub.News{Headline: "words"}
	newNews = append(newNews, article1)
	newNews = append(newNews, article2)
	mock = &newsApiMock{
		mockGeneralNews: func(ctx context.Context, category string, localVarOptionals *finnhub.GeneralNewsOpts) ([]finnhub.News, *http.Response, error) {
			return newNews, nil, nil
		},
	}

	// Since there are multiple messages the first one will be copied, however, since
	// the second message matches the headline it will be rejected
	news, err = getMarketNews(mock, &auth, &headline)
	if err != nil {
		t.Errorf("Unexpected error: %e", err)
	}
	if len(*news) != 1 {
		t.Errorf("Invalid market news, got length: %d, expected: 1", len(*news))
	}
	if (*news)[0].Headline != "word" {
		t.Errorf("Invalid market news, got headline: %s, expected headline: word", (*news)[0].Headline)
	}

	// Test response with no headline passed
	// Only first message should be copied all other news is ignored
	news, err = getMarketNews(mock, &auth, nil)
	if err != nil {
		t.Errorf("Unexpected error: %e", err)
	}
	if len(*news) != 1 {
		t.Errorf("Invalid market news, got length: %d, expected: 1", len(*news))
	}
	if (*news)[0].Headline != "word" {
		t.Errorf("Invalid market news, got headline: %s, expected headline: word", (*news)[0].Headline)
	}

	// Test response with headline passed but multiple valid headlines
	// All news that have a unique headline will be copied
	newNews[1].Headline = "word2"
	news, err = getMarketNews(mock, &auth, &headline)
	if err != nil {
		t.Errorf("Unexpected error: %e", err)
	}
	if len(*news) != 2 {
		t.Errorf("Invalid market news, got length: %d, expected: 1", len(*news))
	}
	if (*news)[0].Headline != "word" {
		t.Errorf("Invalid market news, got headline: %s, expected headline: word", (*news)[0].Headline)
	}
	if (*news)[1].Headline != "word2" {
		t.Errorf("Invalid market news, got headline: %s, expected headline: word2", (*news)[1].Headline)
	}
}
