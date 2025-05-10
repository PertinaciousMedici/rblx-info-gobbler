package methods

import (
	"fmt"
	"net/http"

	"PanoptisMouthNew/bot/connection"
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
)

func FetchRblxUserFriends(botInstance *BotInstance, userID string) (*[]*FetchedRblxFriendship, bool) {
	var apiPath = fmt.Sprintf(connection.RblxUserFriendsSearchGet, userID)
	getRequest, err := http.NewRequest(connection.GET, apiPath, nil)

	if err != nil {
		botInstance.BotChannels.ThrowException("FetchRblxUserFriends", "Failed to build the HTTP GET request", err)
		return nil, false
	}

	getRequest.Header.Set("Content-Type", "application/json")
	getRequest.Header.Set("Accept", "application/json")
	getRequest.Header.Set("x-api-key", botInstance.GlobalDefines["rblxAPIKey"])

	var fetchResponse FetchedRblxUserFriends
	okRequest := utilities.ExecuteBotRequest("FetchRblxUserFriends", botInstance, getRequest, &fetchResponse)

	if !okRequest {
		return nil, false
	}

	return &fetchResponse.Friends, true
}

func FetchMutualRblxFriends(botInstance *BotInstance, userID1 string, userID2 string) (*[]*FetchedRblxFriendship, bool) {
	user1Friends, queryOk := FetchRblxUserFriends(botInstance, userID1)

	if !queryOk {
		return nil, false
	}

	user2Friends, queryOk := FetchRblxUserFriends(botInstance, userID2)

	if !queryOk {
		return nil, false
	}

	friendMap := make(map[uint64]*FetchedRblxFriendship)

	for _, friend := range *user1Friends {
		friendMap[friend.UserID] = friend
	}

	var mutualFriends []*FetchedRblxFriendship

	for _, friend := range *user2Friends {
		if _, exists := friendMap[friend.UserID]; exists {
			mutualFriends = append(mutualFriends, friend)
		}
	}

	return &mutualFriends, true
}
