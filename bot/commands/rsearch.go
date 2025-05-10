package commands

import (
	"strconv"
	"time"

	"PanoptisMouthNew/bot/methods"
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func RSearch(botInstance *BotInstance, event *gateway.MessageCreateEvent, args []string) {
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

		userToPick := (*fetchedUsers)[0]
		userID = strconv.FormatUint(userToPick.UserID, 10)
	}

	fetchedAccounts, okRequest := methods.FetchDiscAccountsRw(botInstance, userID)

	if !okRequest {
		utilities.BadRequest(botInstance, event, "FetchAccount2")
		return
	}

	if len(*fetchedAccounts) < 1 {
		utilities.NotFound(botInstance, event, "RoWifi User")
		return
	}

	for _, account := range *fetchedAccounts {
		discordUser, okRequest := methods.FetchDiscAccountDisc(botInstance, account.DiscordID)

		if !okRequest {
			utilities.BadRequest(botInstance, event, "FetchAccount3")
			continue
		}

		utilities.SendDiscordUser(botInstance, event, discordUser)
		time.Sleep(500 * time.Millisecond)
	}
}
