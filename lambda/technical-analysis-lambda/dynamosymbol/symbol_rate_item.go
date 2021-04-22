package dynamosymbol

type SymbolData struct {
	Open   float32
	High   float32
	Low    float32
	Close  float32
	Volume float32
}

type SymbolRateItem struct {
	Date      string
	Timestamp string
	EurUsd    SymbolData
	GbpUsd    SymbolData
}
