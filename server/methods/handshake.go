package methods

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	serverStructures "PanoptisMouthNew/structures/server"
	"github.com/gorilla/websocket"
)

const maxRetries = 5

func (instance *ServerManager) HandshakeHandler(responseWriter http.ResponseWriter, request *http.Request) {
	serverInstance := instance.ServerInstance
	serverUpgrader := serverInstance.ServerUpgrader
	connection, err := serverUpgrader.Upgrade(responseWriter, request, nil)
	var failed bool

	if err != nil {
		serverInstance.ServerChannels.ThrowException("Handshake", "Failed to upgrade connection", err)
		return
	}

	defer instance.cleanup(connection, failed)
	_, payload, err := connection.ReadMessage()

	if err != nil {
		serverInstance.ServerChannels.ThrowException("Handshake", "Failed to read message", err)
		failed = true
		return
	}

	var handshakePayload serverStructures.HandshakePayload

	jsonDecoder := json.NewDecoder(bytes.NewReader(payload))
	err = jsonDecoder.Decode(&handshakePayload)

	if err != nil {
		serverInstance.ServerChannels.ThrowException("Handshake", "Failed to decode payload", err)
		failed = true
		return
	}

	msNow := uint64(time.Now().UnixMilli())
	websocketLatency := msNow - handshakePayload.CurrentDate

	err = connection.WriteMessage(websocket.PingMessage, nil)

	if err != nil {
		serverInstance.ServerChannels.ThrowException("Handshake", "Failed to write ping message", err)
		failed = true
		return
	}

	connectionClient := serverStructures.Client{
		Discriminator: handshakePayload.ClientUsername,
		Connection:    connection,
		ClientMutex:   &sync.RWMutex{},
		WsLatency:     websocketLatency,
		Active:        atomic.Bool{},
	}
	connectionClient.Active.Store(true)

	{
		serverInstance.ServerMutex.Lock()
		serverInstance.ClientConnections[connectionClient.Discriminator] = &connectionClient
		serverInstance.ConnectionCount.Add(1)
		serverInstance.ServerMutex.Unlock()
	}

	go heartbeatClock(instance, &connectionClient)
}

func heartbeatClock(instance *ServerManager, client *serverStructures.Client) {
	server := instance.ServerInstance

	defer func() {
		server.ServerMutex.Lock()
		delete(server.ClientConnections, client.Discriminator)
		server.ConnectionCount.Add(-1)
		server.ServerMutex.Unlock()

		client.ClientMutex.Lock()
		client.Active.Store(false)
		client.ClientMutex.Unlock()

		server.ServerChannels.PutMessage("HeartbeatClock", "Client deemed inactive, connection dropped")
	}()

	var tryCounter int

	for {
		if tryCounter > maxRetries {
			break
		}

		time.Sleep(30 * time.Second)
		server.ServerChannels.PutMessage("HeartbeatClock", "Attempting heartbeat to client")

		if !client.Active.Load() {
			break
		}

		err := client.Connection.WriteMessage(websocket.PingMessage, nil)

		if websocket.IsCloseError(err,
			websocket.CloseNormalClosure,
			websocket.CloseGoingAway,
			websocket.CloseAbnormalClosure) {
			server.ServerChannels.PutMessage("HeartbeatClock", "Client disconnected from server")
			break
		}

		if err != nil {
			server.ServerChannels.ThrowException("HeartbeatClock", "Failed to write ping message", err)
			tryCounter++
			continue
		}

		readDeadline := time.Now().Add(5 * time.Second)
		err = client.Connection.SetReadDeadline(readDeadline)

		if err != nil {
			server.ServerChannels.ThrowException("HeartbeatClock", "Failed to set read deadline", err)
			tryCounter++
			continue
		}

		_, payload, err := client.Connection.ReadMessage()

		if err != nil {
			server.ServerChannels.ThrowException("HeartbeatClock", "Failed to read payload", err)
			tryCounter++
			continue
		}

		var heartbeatPayload serverStructures.HandshakePayload
		jsonDecoder := json.NewDecoder(bytes.NewReader(payload))
		err = jsonDecoder.Decode(&heartbeatPayload)

		if err != nil {
			server.ServerChannels.ThrowException("HeartbeatClock", "Failed to decode payload", err)
			tryCounter++
			continue
		}

		msNow := uint64(time.Now().UnixMilli())
		websocketLatency := msNow - heartbeatPayload.CurrentDate

		{
			client.ClientMutex.Lock()
			client.WsLatency = websocketLatency
			client.ClientMutex.Unlock()
		}

		heartbeatOutput := fmt.Sprintf("Client heartbeat successful, latency currently: %dms.", websocketLatency)
		server.ServerChannels.PutMessage("HeartbeatClock", heartbeatOutput)
		tryCounter = 0
	}
}
