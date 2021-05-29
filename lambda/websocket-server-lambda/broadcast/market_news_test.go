package broadcast

import (
	"strconv"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/service/dynamodb/expression"
)

func TestInitialMarketNews(t *testing.T) {
	// Test DynamoDB query input parameters with valid information
	currentDate := time.Now().Add(-24 * time.Hour)
	keyCond := expression.Key("Date").Equal(expression.Value(currentDate.Format("2006-01-02")))
	dbInput, err := InitialMarketNews(keyCond)
	if err != nil {
		t.Errorf("DynamoDB query was invalid, failed with error: %e", err)
	}

	tableName := dbInput.TableName
	attributeNames := dbInput.ExpressionAttributeNames
	attributeValues := dbInput.ExpressionAttributeValues
	if *tableName != "SymbolRateTable" {
		t.Errorf("DynamoDB query was invalid, table name was %s, expected: SymbolRateTable", *tableName)
	}
	if *attributeNames["#0"] != "Date" {
		t.Errorf("DynamoDB query was invalid, attribute 0 was %s, expected: Date", *attributeNames["#0"])
	}
	if *attributeNames["#1"] != "MarketNews" {
		t.Errorf("DynamoDB query was invalid, attribute 1 was %s, expected: MarketNews", *attributeNames["#1"])
	}
	if *attributeValues[":0"].S != currentDate.Format("2006-01-02") {
		t.Errorf("DynamoDB query was invalid, expected attribute value was %s, expected: %s", *attributeValues[":0"].S, currentDate.Format("2006-01-02"))
	}
}

func TestFilterMarketNews(t *testing.T) {
	// Test whether a single duplicate is filtered out from headlines
	testNews := []CallbackMessageNews{}

	// Setup for news headlines
	for i := 0; i < 4; i++ {
		news := marketNewsTable{
			Timestamp: strconv.Itoa(i) + ":00",
			Headline:  strconv.Itoa(i) + "words",
			Image:     strconv.Itoa(i) + ".jpg",
			Source:    strconv.Itoa(i) + "forexlive",
			Summary:   strconv.Itoa(i) + "happened",
			NewsUrl:   strconv.Itoa(i) + ".com",
		}
		testNews = append(testNews, CallbackMessageNews{MarketNews: news})
	}
	duplicateNews := marketNewsTable{
		Timestamp: "3" + ":00",
		Headline:  "3" + "words",
		Image:     "3" + ".jpg",
		Source:    "3" + "forexlive",
		Summary:   "3" + "happened",
		NewsUrl:   "3" + ".com",
	}
	testNews = append(testNews, CallbackMessageNews{MarketNews: duplicateNews})

	newsList := filterMarketNews(&testNews)
	// Ensure there are only 4 headlines remaining
	if len(*newsList) != 4 {
		t.Errorf("Filtered market news invalid, got slice of length: %d, expected: 4", len(*newsList))
	}
	// Validate the headlines of the news
	for index, news := range *newsList {
		if news.MarketNews.Headline != strconv.Itoa(index)+"words" {
			t.Errorf("Filtered market news invalid, got headline: %s, expected: %s", news.MarketNews.Headline, strconv.Itoa(index)+"words")
		}
	}

	// Test whether several duplicate headlines can be filtered out
	testNews = []CallbackMessageNews{}

	// Setup for news headlines
	for i := 0; i < 2; i++ {
		news := marketNewsTable{
			Timestamp: strconv.Itoa(i) + ":00",
			Headline:  strconv.Itoa(i) + "words",
			Image:     strconv.Itoa(i) + ".jpg",
			Source:    strconv.Itoa(i) + "forexlive",
			Summary:   strconv.Itoa(i) + "happened",
			NewsUrl:   strconv.Itoa(i) + ".com",
		}
		testNews = append(testNews, CallbackMessageNews{MarketNews: news})
	}
	for i := 0; i < 5; i++ {
		duplicateNews := marketNewsTable{
			Timestamp: "1" + ":00",
			Headline:  "1" + "words",
			Image:     "1" + ".jpg",
			Source:    "1" + "forexlive",
			Summary:   "1" + "happened",
			NewsUrl:   "1" + ".com",
		}
		testNews = append(testNews, CallbackMessageNews{MarketNews: duplicateNews})
	}

	newsList = filterMarketNews(&testNews)
	// Ensure there are only 4 headlines remaining
	if len(*newsList) != 2 {
		t.Errorf("Filtered market news invalid, got slice of length: %d, expected: 2", len(*newsList))
	}
	// Validate the headlines of the news
	for index, news := range *newsList {
		if news.MarketNews.Headline != strconv.Itoa(index)+"words" {
			t.Errorf("Filtered market news invalid, got headline: %s, expected: %s", news.MarketNews.Headline, strconv.Itoa(index)+"words")
		}
	}
}
