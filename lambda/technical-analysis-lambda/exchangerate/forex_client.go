package exchangerate

import (
	"context"
	"net/http"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go"
)

type FinnhubClient interface {
	ForexCandles(ctx context.Context, symbol string, resolution string, from int64, to int64) (finnhub.ForexCandles, *http.Response, error)
	GeneralNews(ctx context.Context, category string, localVarOptionals *finnhub.GeneralNewsOpts) ([]finnhub.News, *http.Response, error)
}

type ForexApi struct {
	client *finnhub.DefaultApiService
}

func (fa *ForexApi) ForexCandles(ctx context.Context, symbol string, resolution string, from int64, to int64) (finnhub.ForexCandles, *http.Response, error) {
	return fa.client.ForexCandles(ctx, symbol, resolution, from, to)
}

func (fa *ForexApi) GeneralNews(ctx context.Context, category string, localVarOptionals *finnhub.GeneralNewsOpts) ([]finnhub.News, *http.Response, error) {
	return fa.client.GeneralNews(ctx, category, localVarOptionals)
}

func NewForexApi() *ForexApi {
	return &ForexApi{client: finnhub.NewAPIClient(finnhub.NewConfiguration()).DefaultApi}
}
