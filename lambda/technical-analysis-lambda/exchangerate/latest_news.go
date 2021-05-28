package exchangerate

import (
	"context"
	"errors"
	"log"
	"technical-analysis-lambda/finance"
	"time"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

func getMarketNews(client *finnhub.DefaultApiService, authContext *context.Context, headline *string) (*[]finance.NewsItem, error) {
	forexNews, _, err := client.GeneralNews(*authContext, "forex", nil)
	if err != nil {
		log.Println("Unable to get latest news")
		return nil, errors.New("news request failed")
	}
	latestNews := []finance.NewsItem{}

	if headline == nil && len(forexNews) > 0 {
		article := forexNews[0]
		timestamp := time.Unix(article.Datetime, 0).UTC()
		formattedTime := timestamp.Format("2006-01-02 15:04:05")
		newsItem := finance.NewsItem{
			Timestamp: formattedTime,
			Headline:  article.Headline,
			Image:     article.Image,
			Source:    article.Source,
			Summary:   article.Summary,
			NewsUrl:   article.Url,
		}
		latestNews = append(latestNews, newsItem)
	} else if len(forexNews) > 0 {
		for _, article := range forexNews {
			if article.Headline == *headline && len(latestNews) == 0 {
				timestamp := time.Unix(article.Datetime, 0).UTC()
				formattedTime := timestamp.Format("2006-01-02 15:04:05")
				newsItem := finance.NewsItem{
					Timestamp: formattedTime,
					Headline:  article.Headline,
					Image:     article.Image,
					Source:    article.Source,
					Summary:   article.Summary,
					NewsUrl:   article.Url,
				}
				latestNews = append(latestNews, newsItem)
				break
			} else if article.Headline == *headline {
				break
			}
			timestamp := time.Unix(article.Datetime, 0).UTC()
			formattedTime := timestamp.Format("2006-01-02 15:04:05")
			newsItem := finance.NewsItem{
				Timestamp: formattedTime,
				Headline:  article.Headline,
				Image:     article.Image,
				Source:    article.Source,
				Summary:   article.Summary,
				NewsUrl:   article.Url,
			}
			latestNews = append(latestNews, newsItem)
		}
	}

	return &latestNews, nil
}
