package finance

type SymbolData struct {
	Open   float32
	High   float32
	Low    float32
	Close  float32
	Volume float32
}

type NewsItem struct {
	Timestamp string
	Headline  string
	Image     string
	Source    string
	Summary   string
	NewsUrl   string
}

type FinancialDataItem struct {
	Date       string
	Timestamp  string
	EURUSD     SymbolData
	GBPUSD     SymbolData
	MarketNews NewsItem
}
