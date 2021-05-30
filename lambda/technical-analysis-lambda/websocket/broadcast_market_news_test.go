package websocket

import (
	"technical-analysis-lambda/finance"
	"testing"
	"time"
)

func TestValidateNewMarketNews(t *testing.T) {
	// Test response when a nil pointer is passed
	headline := "words"
	isValid := ValidateNewMarketNews(nil, headline)
	if isValid != false {
		t.Errorf("Market news was invalid, got: %t", isValid)
	}

	// Test response when empty struct us passed
	isValid = ValidateNewMarketNews(&[]finance.NewsItem{}, headline)
	if isValid != false {
		t.Errorf("Market news was invalid, got: %t", isValid)
	}

	// Test response when headline is not unique
	duplicateNews := []finance.NewsItem{}
	currentDate := time.Now()
	news := finance.NewsItem{
		Timestamp: currentDate.Format("2006-01-02 15:04:05"),
		Headline:  "words",
		Image:     ".jpg",
		Source:    "forexLive",
		Summary:   "",
		NewsUrl:   ".com",
	}
	duplicateNews = append(duplicateNews, news)
	isValid = ValidateNewMarketNews(&duplicateNews, headline)
	if isValid != false {
		t.Errorf("Market news was invalid, got: %t", isValid)
	}

	// Test response when headline is unique
	duplicateNews[0].Headline = "Words"
	isValid = ValidateNewMarketNews(&duplicateNews, headline)
	if isValid != true {
		t.Errorf("Market news was valid, got: %t", isValid)
	}
}

func TestCreateCallbackNewsMessage(t *testing.T) {
	// Test whether websocket payload fields match for a message
	newsSlice := []finance.NewsItem{}
	currentDate := time.Now()
	news := finance.NewsItem{
		Timestamp: currentDate.Format("2006-01-02 15:04:05"),
		Headline:  "words",
		Image:     ".jpg",
		Source:    "ForexLive",
		Summary:   "",
		NewsUrl:   ".com",
	}
	newsSlice = append(newsSlice, news)

	// Check whether fields match
	payload := createCallbackNewsMessage(&newsSlice)
	if (*payload)[0].Timestamp != news.Timestamp {
		t.Errorf("Payload timestamp is invalid, got: %s, expected: %s", (*payload)[0].Timestamp, news.Timestamp)
	}
	if (*payload)[0].Headline != news.Headline {
		t.Errorf("Payload headline is invalid, got: %s, expected: %s", (*payload)[0].Headline, news.Headline)
	}
	if (*payload)[0].Image != news.Image {
		t.Errorf("Payload image is invalid, got: %s, expected: %s", (*payload)[0].Image, news.Image)
	}
	if (*payload)[0].Source != news.Source {
		t.Errorf("Payload source is invalid, got: %s, expected: %s", (*payload)[0].Source, news.Source)
	}
	if (*payload)[0].Summary != news.Summary {
		t.Errorf("Payload summary is invalid, got: %s, expected: %s", (*payload)[0].Summary, news.Summary)
	}
	if (*payload)[0].NewsUrl != news.NewsUrl {
		t.Errorf("Payload news url is invalid, got: %s, expected: %s", (*payload)[0].NewsUrl, news.NewsUrl)
	}
}
