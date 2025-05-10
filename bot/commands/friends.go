package commands

import (
	"fmt"
	"strconv"

	"PanoptisMouthNew/bot/connection"
	"PanoptisMouthNew/bot/design"
	"PanoptisMouthNew/bot/methods"
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func Friends(botInstance *BotInstance, event *gateway.MessageCreateEvent, args []string) {
	firstArgument := args[0]
	isNumericArgument := utilities.IsUint64Type(firstArgument)
	var userID string
	var username string

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
		username = firstArgument
	}

	if username == "" {
		fetchedUser, okRequest := methods.FetchRblxAccount(botInstance, userID)

		if !okRequest {
			utilities.BadRequest(botInstance, event, "FetchAccount2")
			return
		}

		if fetchedUser.Username == "" {
			utilities.NotFound(botInstance, event, "Roblox User")
			return
		}

		username = fetchedUser.Username
	}

	fetchedUserFriends, okRequest := methods.FetchRblxUserFriends(botInstance, userID)

	if !okRequest {
		utilities.BadRequest(botInstance, event, "FetchFriends1")
		return
	}

	if len(*fetchedUserFriends) < 1 {
		utilities.NotFound(botInstance, event, "User Friends")
		return
	}

	maxPages := utilities.CalculateMaxPages(design.MaximumFieldsPerPageFriends, len(*fetchedUserFriends))

	fieldBuilder := func(friendship *FetchedRblxFriendship) discord.EmbedField {
		userProfileURL := fmt.Sprintf(connection.RblxUserView, friendship.UserID)
		userProfileText := fmt.Sprintf("[[%s]](%s)", friendship.Username, userProfileURL)

		embedField := discord.EmbedField{
			Value: fmt.Sprintf("**%s**\n**Display Name:** `%s`\n", userProfileText, friendship.DisplayName),
		}

		return embedField
	}

	friendsEmbed := utilities.BuildPaginatedEmbed[FetchedRblxFriendship](
		fmt.Sprintf("%s's Friends", username),
		"",
		design.MaximumFieldsPerPageFriends,
		maxPages,
		1,
		fetchedUserFriends,
		fieldBuilder,
		"",
	)

	utilities.SendIterator(
		botInstance,
		event.Message.ChannelID,
		"friends",
		1,
		&friendsEmbed,
		"Friends",
	)

	botInstance.InstanceMutex.Lock()
	if _, ok := botInstance.BotStorage.StoredFriends[username]; ok {
		delete(botInstance.BotStorage.StoredFriends, username)
	}
	botInstance.BotStorage.StoredFriends[username] = &FetchedRblxUserFriends{
		Friends: *fetchedUserFriends,
	}
	botInstance.InstanceMutex.Unlock()
}
