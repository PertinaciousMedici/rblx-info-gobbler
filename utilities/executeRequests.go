package utilities

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"time"

	botStructures "PanoptisMouthNew/structures/bot"
)

func BuildRequestBody(caller string, botInstance *botStructures.BotInstance, source interface{}, target **bytes.Reader) bool {
	var bytesBuffer bytes.Buffer

	JSONEncoder := json.NewEncoder(&bytesBuffer)
	JSONEncoder.SetIndent("", "  ")
	err := JSONEncoder.Encode(source)

	if err != nil {
		botInstance.BotChannels.ThrowException(caller, "Failed to encode the query request", err)
		return false
	}

	bytesReader := bytes.NewReader(bytesBuffer.Bytes())
	*target = bytesReader

	return true
}

func ExecuteBotRequest(caller string, botInstance *botStructures.BotInstance, request *http.Request, targetVariable interface{}) bool {
	time.Sleep(time.Duration(500) * time.Millisecond)
	response, err := botInstance.HttpClient.Do(request)

	if err != nil {
		botInstance.BotChannels.ThrowException(caller, "Failed to execute the HTTP request", err)
		return false
	}

	defer func() {
		err = response.Body.Close()
		if err != nil {
			botInstance.BotChannels.ThrowException(caller, "Failed to close response body", err)
		}
	}()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		botInstance.BotChannels.ThrowException(caller, "Failed to read response body", err)
		return false
	}

	err = json.NewDecoder(bytes.NewReader(body)).Decode(targetVariable)

	if err != nil {
		botInstance.BotChannels.ThrowException(caller, "Failed to decode the JSON response", err)
		return false
	}

	return true
}
