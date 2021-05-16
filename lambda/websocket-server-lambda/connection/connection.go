package connection

type Connection struct {
	ConnectionId    string
	EURUSD          bool
	GBPUSD          bool
	USDJPY          bool
	AUDCAD          bool
	EURUSDInference bool
	GBPUSDInference bool
	USDJPYInference bool
	AUDCADInference bool
	MarketNews      bool
}

func CreateNewConnection(connectionId string) *Connection {
	return &Connection{
		ConnectionId:    connectionId,
		EURUSD:          false,
		GBPUSD:          false,
		USDJPY:          false,
		AUDCAD:          false,
		EURUSDInference: false,
		GBPUSDInference: false,
		USDJPYInference: false,
		AUDCADInference: false,
		MarketNews:      false,
	}
}
