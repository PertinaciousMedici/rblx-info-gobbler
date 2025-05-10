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

func MFriends(botInstance *BotInstance, event *gateway.MessageCreateEvent, args []string) {
	if args[0] == args[1] {
		utilities.IllogicalRequest(botInstance, event, "Arguments 1 and 2 are identical", true)
		return
	}

	firstArgument := args[0]
	isNumericArgument1 := utilities.IsUint64Type(firstArgument)
	var robloxID1 string
	var robloxUsername1 string

	if isNumericArgument1 {
		robloxID1 = firstArgument
	} else {
		fetchedUsers, okRequest := methods.FetchRblxAccountsByUsernames(botInstance, []string{firstArgument})

		if !okRequest {
			utilities.BadRequest(botInstance, event, "FetchAccount1")
			return
		}

		if len(*fetchedUsers) < 1 {
			utilities.NotFound(botInstance, event, "Roblox User 1")
			return
		}

		robloxID1 = strconv.FormatUint((*fetchedUsers)[0].UserID, 10)
		robloxUsername1 = firstArgument
	}

	if robloxUsername1 == "" {
		fetchedUser, okRequest := methods.FetchRblxAccount(botInstance, robloxID1)

		if !okRequest {
			utilities.BadRequest(botInstance, event, "FetchAccount2")
			return
		}

		if fetchedUser.Username == "" {
			utilities.NotFound(botInstance, event, "Roblox User 1")
			return
		}

		robloxUsername1 = fetchedUser.Username
	}

	secondArgument := args[1]
	isNumericArgument2 := utilities.IsUint64Type(secondArgument)
	var robloxID2 string
	var robloxUsername2 string

	if isNumericArgument2 {
		robloxID2 = secondArgument
	} else {
		fetchedUsers, okRequest := methods.FetchRblxAccountsByUsernames(botInstance, []string{secondArgument})

		if !okRequest {
			utilities.BadRequest(botInstance, event, "FetchAccount3")
			return
		}

		if len(*fetchedUsers) < 1 {
			utilities.NotFound(botInstance, event, "Roblox User 2")
			return
		}

		robloxID2 = strconv.FormatUint((*fetchedUsers)[0].UserID, 10)
		robloxUsername2 = secondArgument
	}

	if robloxUsername2 == "" {
		fetchedUser, okRequest := methods.FetchRblxAccount(botInstance, robloxID2)

		if !okRequest {
			utilities.BadRequest(botInstance, event, "FetchAccount4")
		}

		if fetchedUser.Username == "" {
			utilities.NotFound(botInstance, event, "Roblox User 2")
			return
		}
	}

	fetchedMutuals, okRequest := methods.FetchMutualRblxFriends(botInstance, robloxID1, robloxID2)

	if !okRequest {
		utilities.BadRequest(botInstance, event, "FetchMutual")
		return
	}

	if len(*fetchedMutuals) < 1 {
		utilities.NotFound(botInstance, event, "Mutual Friends")
		return
	}

	maxPages := utilities.CalculateMaxPages(design.MaximumFieldsPerPageFriends, len(*fetchedMutuals))
	uniqueIdentifier := strconv.FormatUint(uint64(event.Message.ID), 10)

	fieldBuilder := func(friendship *FetchedRblxFriendship) discord.EmbedField {
		userProfileURL := fmt.Sprintf(connection.RblxUserView, friendship.UserID)
		userProfileText := fmt.Sprintf("[[%s]](%s)", friendship.Username, userProfileURL)

		embedField := discord.EmbedField{
			Value: fmt.Sprintf("**%s**\n**Display Name:** `%s`\n", userProfileText, friendship.DisplayName),
		}

		return embedField
	}

	friendsEmbed := utilities.BuildPaginatedEmbed[FetchedRblxFriendship](
		"Mutual Friends",
		uniqueIdentifier,
		design.MaximumFieldsPerPageFriends,
		maxPages,
		1,
		fetchedMutuals,
		fieldBuilder,
		fmt.Sprintf("**%s** & **%s**", robloxUsername1, robloxUsername2),
	)

	utilities.SendIterator(
		botInstance,
		event.Message.ChannelID,
		"mutualf",
		1,
		&friendsEmbed,
		"MutualFriends",
	)

	botInstance.InstanceMutex.Lock()
	if _, ok := botInstance.BotStorage.StoredMutualF[uniqueIdentifier]; ok {
		delete(botInstance.BotStorage.StoredMutualF, uniqueIdentifier)
	}
	botInstance.BotStorage.StoredMutualF[uniqueIdentifier] = &FetchedRblxUserFriends{
		Friends: *fetchedMutuals,
	}
	botInstance.InstanceMutex.Unlock()
}
