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

func MGroups(botInstance *BotInstance, event *gateway.MessageCreateEvent, args []string) {
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
		fetchedUser, okRequest := methods.FetchRblxAccount(botInstance, firstArgument)

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
		fetchedUser, okRequest := methods.FetchRblxAccount(botInstance, secondArgument)

		if !okRequest {
			utilities.BadRequest(botInstance, event, "FetchAccount4")
			return
		}

		if fetchedUser.Username == "" {
			utilities.NotFound(botInstance, event, "Roblox User 2")
			return
		}

		robloxUsername2 = fetchedUser.Username
	}

	fetchedMutuals, okRequest := methods.FetchMutualRblxGroups(botInstance, robloxID1, robloxID2)

	if !okRequest {
		utilities.BadRequest(botInstance, event, "Mutual Groups")
		return
	}

	if len(*fetchedMutuals) < 1 {
		utilities.NotFound(botInstance, event, "Mutual Groups")
		return
	}

	maxPages := utilities.CalculateMaxPages(design.MaximumFieldsPerPageGroups, len(*fetchedMutuals))
	uniqueIdentifier := strconv.FormatUint(uint64(event.Message.ID), 10)

	fieldBuilder := func(membership *FetchedUserMembership) discord.EmbedField {
		groupPageURL := fmt.Sprintf(connection.RblxGroupView, membership.Group.GroupID)
		groupMembershipText := fmt.Sprintf("[[%s]](%s)", membership.Group.GroupName, groupPageURL)

		embedField := discord.EmbedField{
			Value: fmt.Sprintf("**%s**\n**Member Count:** `%d`\n",
				groupMembershipText,
				membership.Group.MemberCount,
			),
		}

		return embedField
	}

	groupsEmbed := utilities.BuildPaginatedEmbed[FetchedUserMembership](
		"Mutual Groups",
		uniqueIdentifier,
		design.MaximumFieldsPerPageGroups,
		maxPages,
		1,
		fetchedMutuals,
		fieldBuilder,
		fmt.Sprintf("**%s** & **%s**", robloxUsername1, robloxUsername2),
	)

	utilities.SendIterator(
		botInstance,
		event.Message.ChannelID,
		"mutualg",
		1,
		&groupsEmbed,
		"MutualGroups",
	)

	botInstance.InstanceMutex.Lock()
	if _, ok := botInstance.BotStorage.StoredMutualG[uniqueIdentifier]; ok {
		delete(botInstance.BotStorage.StoredMutualG, uniqueIdentifier)
	}
	botInstance.BotStorage.StoredMutualG[uniqueIdentifier] = &FetchedUserGroupMemberships{
		Groups: *fetchedMutuals,
	}
	botInstance.InstanceMutex.Unlock()
}
