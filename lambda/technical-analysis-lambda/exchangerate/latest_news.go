package exchangerate

import (
	"context"
	"errors"
	"log"
	"technical-analysis-lambda/finance"
	"time"
)

// Queries the Finnhub API to get forex news
// Filters out news that have already been stored in db
func getMarketNews(fc FinnhubClient, authContext *context.Context, headline *string) (*[]finance.NewsItem, error) {
	forexNews, _, err := fc.GeneralNews(*authContext, "forex", nil)
	if err != nil {
		log.Println("Unable to get latest news")
		return nil, errors.New("news request failed")
	}
	latestNews := []finance.NewsItem{}

	if headline == nil && len(forexNews) > 0 {
		// If there is no latest headline, the newest headline from Finnhub is used
		// as the new headline
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
		// Keep looping through headlines until the headline matches the previous latest
		// headline stored in DynamoDB
		for _, article := range forexNews {
			if article.Headline == *headline && len(latestNews) == 0 {
				// If there is no new market news the old news is stored
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
