package connection

type Connection struct {
	ConnectionId    string
	EURUSD          bool
	GBPUSD          bool
	EURUSDInference bool
	GBPUSDInference bool
	MarketNews      bool
}

func CreateNewConnection(connectionId string) *Connection {
	return &Connection{
		ConnectionId:    connectionId,
		EURUSD:          false,
		GBPUSD:          false,
		EURUSDInference: false,
		GBPUSDInference: false,
		MarketNews:      false,
	}
}
