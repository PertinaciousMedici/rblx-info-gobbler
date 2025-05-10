package methods

import (
	"bytes"
	"encoding/json"
	"net/http"

	serverStructures "PanoptisMouthNew/structures/server"
)

func (instance *ServerManager) PostMessage(responseWriter http.ResponseWriter, request *http.Request) {
	serverInstance := instance.ServerInstance
	serverUpgrader := serverInstance.ServerUpgrader
	connection, err := serverUpgrader.Upgrade(responseWriter, request, nil)
	var failed bool

	if err != nil {
		serverInstance.ServerChannels.ThrowException("PostMessage", "Failed to upgrade connection", err)
		return
	}

	defer instance.cleanup(connection, failed)

	_, payload, err := connection.ReadMessage()

	if err != nil {
		serverInstance.ServerChannels.ThrowException("PostMessage", "Failed to read message", err)
		failed = true
		return
	}

	var postPayload serverStructures.MessagePostRequest

	jsonDecoder := json.NewDecoder(bytes.NewReader(payload))
	err = jsonDecoder.Decode(&postPayload)

	if err != nil {
		serverInstance.ServerChannels.ThrowException("PostMessage", "Failed to decode payload", err)
		failed = true
		return
	}

	serverInstance.ServerChannels.PostChannel <- &postPayload
}
