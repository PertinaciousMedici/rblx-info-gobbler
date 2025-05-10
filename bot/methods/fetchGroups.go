package methods

import (
	"fmt"
	"net/http"

	"PanoptisMouthNew/bot/connection"
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
)

// FetchRblxGroupInfo -> Should be rewritten with better structs.
func FetchRblxGroupInfo(botInstance *BotInstance, groupID string) (*FetchedFullRblxGroup, bool) {
	var apiPath = fmt.Sprintf(connection.RblxGroupSearchGet, groupID)
	getRequest, err := http.NewRequest(connection.GET, apiPath, nil)

	if err != nil {
		botInstance.BotChannels.ThrowException("FetchRblxGroups", "Failed to build the HTTP GET request", err)
		return nil, false
	}

	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("x-api-key", botInstance.GlobalDefines["rblxAPIKey"])

	var fetchResponse FetchedFullRblxGroup
	okRequest := utilities.ExecuteBotRequest("FetchRblxGroups", botInstance, getRequest, &fetchResponse)

	if !okRequest {
		return nil, false
	}

	return &fetchResponse, true
}

func FetchRblxUserGroups(botInstance *BotInstance, userID string) (*[]*FetchedUserMembership, bool) {
	var apiPath = fmt.Sprintf(connection.RblxUserGroupsSearchGet, userID)
	getRequest, err := http.NewRequest(connection.GET, apiPath, nil)

	if err != nil {
		botInstance.BotChannels.ThrowException("FetchRblxUserGroups", "Failed to build the query HTTP request", err)
		return nil, false
	}

	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("x-api-key", botInstance.GlobalDefines["rblxAPIKey"])

	var fetchResponse FetchedUserGroupMemberships
	okRequest := utilities.ExecuteBotRequest("FetchRblxUserGroups", botInstance, getRequest, &fetchResponse)

	if !okRequest {
		return nil, false
	}

	return &fetchResponse.Groups, true
}

func FetchMutualRblxGroups(botInstance *BotInstance, userID1 string, userID2 string) (*[]*FetchedUserMembership, bool) {
	user1Memberships, queryOk := FetchRblxUserGroups(botInstance, userID1)

	if !queryOk {
		return nil, false
	}

	user2Memberships, queryOk := FetchRblxUserGroups(botInstance, userID2)

	if !queryOk {
		return nil, false
	}

	groupMap := make(map[uint64]*FetchedUserMembership)

	for _, membership := range *user1Memberships {
		groupMap[membership.Group.GroupID] = membership
	}

	var mutualGroups []*FetchedUserMembership

	for _, membership := range *user2Memberships {
		if _, exists := groupMap[membership.Group.GroupID]; exists {
			mutualGroups = append(mutualGroups, membership)
		}
	}

	return &mutualGroups, true
}
