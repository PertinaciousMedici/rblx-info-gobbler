package methods

import (
	"github.com/gorilla/websocket"
)

func (instance *ServerManager) cleanup(connection *websocket.Conn, failed bool) {
	if failed {
		server := instance.ServerInstance
		err := connection.Close()

		if err != nil {
			server.ServerChannels.ThrowException("cleanup", "Failed to close connection", err)
			return
		}
	}
}
