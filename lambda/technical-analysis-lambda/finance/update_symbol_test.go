package finance

import (
	"testing"
	"time"
)

func TestMergeNewsItems(t *testing.T) {
	// Test response when no news is provided
	symbolRate := []FinancialDataItem{}
	currentDate := time.Now()
	financialItem := FinancialDataItem{Date: currentDate.Format("2006-01-02"), Timestamp: currentDate.Format("15:04:05")}
	symbolRate = append(symbolRate, financialItem)
	mergedItem := mergeNewsItems(&symbolRate, &[]NewsItem{})

	if len(*mergedItem) != 1 {
		t.Errorf("Invalid merged item length, got %d, expected: 1", len(*mergedItem))
	}
	if (*mergedItem)[0].Date != currentDate.Format("2006-01-02") {
		t.Errorf("Invalid merged item date, got: %s, expected: %s", (*mergedItem)[0].Date, currentDate.Format("2006-01-02"))
	}
	if (*mergedItem)[0].Timestamp != currentDate.Format("15:04:05") {
		t.Errorf("Invalid merged item timestamp, got: %s, expected: %s", (*mergedItem)[0].Timestamp, currentDate.Format("15:04:05"))
	}

	// Test response when news is provided
	latestNews := []NewsItem{}
	article := NewsItem{
		Timestamp: currentDate.Format("2006-01-02 15:04:05"),
		Headline:  "words",
		Image:     ".jpg",
		Source:    "Forex Live",
		Summary:   "",
		NewsUrl:   ".com",
	}
	latestNews = append(latestNews, article)

	mergedItem = mergeNewsItems(&symbolRate, &latestNews)
	if len(*mergedItem) != 1 {
		t.Errorf("Invalid merged item length, got %d, expected: 1", len(*mergedItem))
	}
	if (*mergedItem)[0].Date != currentDate.Format("2006-01-02") {
		t.Errorf("Invalid merged item date, got: %s, expected: %s", (*mergedItem)[0].Date, currentDate.Format("2006-01-02"))
	}
	if (*mergedItem)[0].Timestamp != currentDate.Format("15:04:05") {
		t.Errorf("Invalid merged item timestamp, got: %s, expected: %s", (*mergedItem)[0].Timestamp, currentDate.Format("15:04:05"))
	}

	marketNews := (*mergedItem)[0].MarketNews
	if marketNews.Timestamp != article.Timestamp {
		t.Errorf("Invalid merged item market news timestamp, got: %s, expected: %s", marketNews.Timestamp, article.Timestamp)
	}
	if marketNews.Headline != article.Headline {
		t.Errorf("Invalid merged item market news headline, got: %s, expected: %s", marketNews.Headline, article.Headline)
	}
	if marketNews.Image != article.Image {
		t.Errorf("Invalid merged item market news image, got: %s, expected: %s", marketNews.Image, article.Image)
	}
	if marketNews.Source != article.Source {
		t.Errorf("Invalid merged item market news source, got: %s, expected: %s", marketNews.Source, article.Source)
	}
	if marketNews.Summary != article.Summary {
		t.Errorf("Invalid merged item market news summary, got: %s, expected: %s", marketNews.Summary, article.Summary)
	}
	if marketNews.NewsUrl != article.NewsUrl {
		t.Errorf("Invalid merged item market news url, got: %s, expected: %s", marketNews.NewsUrl, article.NewsUrl)
	}

	// Test response when there is more news than exchange rates
	additionalNews := NewsItem{
		Timestamp: currentDate.Add(-time.Hour * 24).Format("2006-01-02 15:04:05"),
		Headline:  "words2",
		Image:     ".png",
		Source:    "DailyFX",
		Summary:   "words",
		NewsUrl:   ".ca",
	}
	latestNews = append(latestNews, additionalNews)

	mergedItem = mergeNewsItems(&symbolRate, &latestNews)
	if len(*mergedItem) != 2 {
		t.Errorf("Invalid merged item length, got %d, expected: 2", len(*mergedItem))
	}

	marketNews = (*mergedItem)[1].MarketNews
	if marketNews.Timestamp != additionalNews.Timestamp {
		t.Errorf("Invalid merged item market news timestamp, got: %s, expected: %s", marketNews.Timestamp, article.Timestamp)
	}
	if marketNews.Headline != additionalNews.Headline {
		t.Errorf("Invalid merged item market news headline, got: %s, expected: %s", marketNews.Headline, article.Headline)
	}
	if marketNews.Image != additionalNews.Image {
		t.Errorf("Invalid merged item market news image, got: %s, expected: %s", marketNews.Image, article.Image)
	}
	if marketNews.Source != additionalNews.Source {
		t.Errorf("Invalid merged item market news source, got: %s, expected: %s", marketNews.Source, article.Source)
	}
	if marketNews.Summary != additionalNews.Summary {
		t.Errorf("Invalid merged item market news summary, got: %s, expected: %s", marketNews.Summary, article.Summary)
	}
	if marketNews.NewsUrl != additionalNews.NewsUrl {
		t.Errorf("Invalid merged item market news url, got: %s, expected: %s", marketNews.NewsUrl, article.NewsUrl)
	}
}
