package methods

import (
	"fmt"
	"net/http"

	"PanoptisMouthNew/bot/connection"
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
)

func FetchAvatarURL(botInstance *BotInstance, userID string) (string, bool) {
	var apiPath = fmt.Sprintf(connection.RblxAvatarSearchGet, userID)
	getRequest, err := http.NewRequest(connection.GET, apiPath, nil)

	if err != nil {
		botInstance.BotChannels.ThrowException("FetchAvatarURL", "Failed to build the HTTP GET request", err)
		return "", false
	}

	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("x-api-key", botInstance.GlobalDefines["rblxAPIKey"])

	var fetchResponse FetchedRblxUserAvatar
	okRequest := utilities.ExecuteBotRequest("FetchAvatarURL", botInstance, getRequest, &fetchResponse)

	if okRequest && len(fetchResponse.QueryResult) > 0 {
		return fetchResponse.QueryResult[0].ImageURL, true
	}

	return "", false
}
