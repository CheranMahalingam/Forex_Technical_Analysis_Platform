package connection

type Connection struct {
	ConnectionId string
	EURUSD       bool
	GBPUSD       bool
}

func CreateNewConnection(connectionId string) *Connection {
	return &Connection{ConnectionId: connectionId, EURUSD: false, GBPUSD: false}
}
