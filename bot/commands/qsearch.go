package commands

import (
	"strconv"

	"PanoptisMouthNew/bot/methods"
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func QSearch(botInstance *BotInstance, event *gateway.MessageCreateEvent, args []string) {
	firstArgument := args[0]
	isNumericArgument := utilities.IsUint64Type(firstArgument)
	var userID string

	if isNumericArgument {
		userID = firstArgument
	} else {
		fetchedUsers, okRequest := methods.FetchRblxAccountsByUsernames(botInstance, []string{firstArgument})
		if !okRequest {
			utilities.BadRequest(botInstance, event, "FetchAccount1")
			return
		}

		if len(*fetchedUsers) < 1 {
			utilities.NotFound(botInstance, event, "Roblox User")
			return
		}

		userID = strconv.FormatUint((*fetchedUsers)[0].UserID, 10)
	}

	fetchedUser, okRequest := methods.FetchRblxAccount(botInstance, userID)

	if !okRequest {
		utilities.BadRequest(botInstance, event, "FetchAccount2")
		return
	}

	if fetchedUser.Username == "" {
		utilities.NotFound(botInstance, event, "Roblox User")
		return
	}

	userAvatar, okRequest := methods.FetchAvatarURL(botInstance, fetchedUser.UserID)

	if okRequest {
		fetchedUser.AvatarURL = userAvatar
	}

	userFriends, okRequest := methods.FetchRblxUserFriends(botInstance, userID)

	if okRequest {
		fetchedUser.FriendCount = len(*userFriends)
	}
	
	utilities.SendRobloxUser(botInstance, event, fetchedUser)
}
