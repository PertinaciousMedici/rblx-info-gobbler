package methods

import (
	"bytes"
	"fmt"
	"net/http"

	"PanoptisMouthNew/bot/connection"
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
)

func FetchRblxAccountsByUsernames(botInstance *BotInstance, usernames []string) (*[]*FetchedUserByUsernames, bool) {
	fetchRequest := RequestQueryByUsernames{
		Usernames:          usernames,
		ExcludeBannedUsers: false,
	}

	var requestBody *bytes.Reader
	utilities.BuildRequestBody("FetchRblxAccountsByUsernames", botInstance, &fetchRequest, &requestBody)
	postRequest, err := http.NewRequest(connection.POST, connection.RblxUsernameSearchPost, requestBody)

	if err != nil {
		botInstance.BotChannels.ThrowException("FetchRblxAccountsByUsernames", "Failed to build the query HTTP request", err)
		return nil, false
	}

	postRequest.Header.Set("Content-Type", "application/json")
	postRequest.Header.Set("Accept", "application/json")

	var fetchResponse ResponseQueryByUsernames
	okRequest := utilities.ExecuteBotRequest("FetchRblxAccountsByUsernames", botInstance, postRequest, &fetchResponse)

	if !okRequest {
		return nil, false
	}

	return &fetchResponse.Users, true
}

func FetchRblxAccount(botInstance *BotInstance, userID string) (*FetchedUserById, bool) {
	apiPath := fmt.Sprintf(connection.RblxUserSearchGet, userID)
	getRequest, err := http.NewRequest(connection.GET, apiPath, nil)

	if err != nil {
		botInstance.BotChannels.ThrowException("FetchRblxAccount", "Failed to build the query HTTP request", err)
		return nil, false
	}

	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("x-api-key", botInstance.GlobalDefines["rblxAPIKey"])

	var fetchResponse FetchedUserById
	okRequest := utilities.ExecuteBotRequest("FetchRblxAccount", botInstance, getRequest, &fetchResponse)

	if !okRequest {
		return nil, false
	}

	return &fetchResponse, true
}

func FetchRblxPastUsernames(botInstance *BotInstance, userID string) []string {
	apiPath := fmt.Sprintf(connection.RblxPastUsernameSearchGet, userID)
	getRequest, err := http.NewRequest(connection.GET, apiPath, nil)

	if err != nil {
		botInstance.BotChannels.ThrowException("FetchRblxUsernames", "Failed to build the query HTTP request", err)
		return nil
	}

	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("x-api-key", botInstance.GlobalDefines["rblxAPIKey"])

	var fetchResponse struct {
		Data []struct {
			Name string `json:"name"`
		} `json:"data"`
	}

	okRequest := utilities.ExecuteBotRequest("FetchRblxUsernames", botInstance, getRequest, &fetchResponse)

	if !okRequest {
		return nil
	}

	var returnValue []string
	for _, entry := range fetchResponse.Data {
		returnValue = append(returnValue, entry.Name)
	}

	return returnValue
}

func FetchDiscAccountsRw(botInstance *BotInstance, userID string) (*[]RwFetchedUserByReverseSearch, bool) {
	apiPath := fmt.Sprintf(connection.RwReverseSearchGet, botInstance.GlobalDefines["homeGuild"], userID)
	getRequest, err := http.NewRequest(connection.GET, apiPath, nil)

	if err != nil {
		botInstance.BotChannels.ThrowException("FetchDiscAccountsRw", "Failed to build the query HTTP request", err)
		return nil, false
	}

	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("Authorization", "Bot "+botInstance.GlobalDefines["rwAPIKey"])

	var fetchResponse []RwFetchedUserByReverseSearch
	okRequest := utilities.ExecuteBotRequest("FetchDiscAccountsRw", botInstance, getRequest, &fetchResponse)

	if !okRequest {
		return nil, false
	}

	return &fetchResponse, true
}

func FetchRblxAccountRw(botInstance *BotInstance, discordID string) (*RwFetchedUserByRegularSearch, bool) {
	apiPath := fmt.Sprintf(connection.RwRegularSearchGet, botInstance.GlobalDefines["homeGuild"], discordID)
	getRequest, err := http.NewRequest(connection.GET, apiPath, nil)

	if err != nil {
		botInstance.BotChannels.ThrowException("FetchRblxAccountRw", "Failed to build the query HTTP request", err)
		return nil, false
	}

	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("Authorization", "Bot "+botInstance.GlobalDefines["rwAPIKey"])

	var fetchResponse RwFetchedUserByRegularSearch
	okRequest := utilities.ExecuteBotRequest("FetchRblxAccountRw", botInstance, getRequest, &fetchResponse)

	if !okRequest {
		return nil, false
	}

	return &fetchResponse, true
}

func FetchDiscAccountDisc(botInstance *BotInstance, userID string) (*DiscFetchedUserById, bool) {
	apiPath := fmt.Sprintf(connection.DsUserSearchGet, userID)
	getRequest, err := http.NewRequest(connection.GET, apiPath, nil)

	if err != nil {
		botInstance.BotChannels.ThrowException("FetchDiscAccountDisc", "Failed to build the query HTTP request", err)
		return nil, false
	}

	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("Authorization", "Bot "+botInstance.GlobalDefines["botToken"])

	var fetchResponse *DiscFetchedUserById
	okRequest := utilities.ExecuteBotRequest("FetchDiscAccountDisc", botInstance, getRequest, &fetchResponse)

	if !okRequest {
		return nil, false
	}

	return fetchResponse, true
}
