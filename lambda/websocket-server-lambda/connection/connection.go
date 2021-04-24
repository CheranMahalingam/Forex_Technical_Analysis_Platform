package connection

type Connection struct {
	ConnectionId string
	EurUsd       bool
	GbpUsd       bool
}

func CreateNewConnection(connectionId string) *Connection {
	return &Connection{ConnectionId: connectionId, EurUsd: false, GbpUsd: false}
}
