package commands

import (
	"strconv"

	"PanoptisMouthNew/bot/methods"
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func Search(botInstance *BotInstance, event *gateway.MessageCreateEvent, args []string) {
	firstArgument := args[0]
	userID := utilities.TransformMention(firstArgument)

	if firstArgument == "self" {
		userID = strconv.FormatUint(uint64(event.Message.Author.ID), 10)
	}

	isNumericArgument := utilities.IsUint64Type(userID)
	if !isNumericArgument {
		utilities.MalformedArguments(botInstance, event, 0, "number", "string")
		return
	}

	fetchedUser, okRequest := methods.FetchRblxAccountRw(botInstance, userID)

	if !okRequest {
		utilities.BadRequest(botInstance, event, "FetchAccount1")
		return
	}

	fetchedRblxUser, okRequest := methods.FetchRblxAccount(botInstance, strconv.FormatUint(fetchedUser.RobloxID, 10))

	if !okRequest {
		utilities.BadRequest(botInstance, event, "FetchAccount2")
		return
	}

	if fetchedRblxUser.Username == "" {
		utilities.NotFound(botInstance, event, "RoWifi User")
		return
	}

	rblxID := strconv.FormatUint(fetchedUser.RobloxID, 10)

	userAvatar, okRequest := methods.FetchAvatarURL(botInstance, rblxID)

	if okRequest {
		fetchedRblxUser.AvatarURL = userAvatar
	}

	userFriends, okRequest := methods.FetchRblxUserFriends(botInstance, rblxID)

	if okRequest {
		fetchedRblxUser.FriendCount = len(*userFriends)
	}

	utilities.SendRobloxUser(botInstance, event, fetchedRblxUser)
}
