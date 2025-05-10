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

func Groups(botInstance *BotInstance, event *gateway.MessageCreateEvent, args []string) {
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

	fetchedUserGroups, okRequest := methods.FetchRblxUserGroups(botInstance, userID)

	if !okRequest {
		utilities.BadRequest(botInstance, event, "FetchGroups1")
		return
	}

	if len(*fetchedUserGroups) < 1 {
		utilities.NotFound(botInstance, event, "User Groups")
		return
	}

	maxPages := utilities.CalculateMaxPages(design.MaximumFieldsPerPageGroups, len(*fetchedUserGroups))

	fieldBuilder := func(membership *FetchedUserMembership) discord.EmbedField {
		groupPageURL := fmt.Sprintf(connection.RblxGroupView, membership.Group.GroupID)
		groupMembershipText := fmt.Sprintf("[[%s]](%s)", membership.Group.GroupName, groupPageURL)

		embedField := discord.EmbedField{
			Value: fmt.Sprintf("**%s**\n**Member Count:** `%d`\n**Rank ID:** `%d`\n **Rank Name:** `%s`",
				groupMembershipText,
				membership.Group.MemberCount,
				membership.Role.RankID,
				membership.Role.RankName,
			),
		}

		return embedField
	}

	groupsEmbed := utilities.BuildPaginatedEmbed[FetchedUserMembership](
		fmt.Sprintf("%s's Groups", username),
		"",
		design.MaximumFieldsPerPageGroups,
		maxPages,
		1,
		fetchedUserGroups,
		fieldBuilder,
		"",
	)

	utilities.SendIterator(
		botInstance,
		event.Message.ChannelID,
		"groups",
		1,
		&groupsEmbed,
		"Groups",
	)

	botInstance.InstanceMutex.Lock()
	if _, ok := botInstance.BotStorage.StoredGroups[username]; ok {
		delete(botInstance.BotStorage.StoredGroups, username)
	}
	botInstance.BotStorage.StoredGroups[username] = &FetchedUserGroupMemberships{
		Groups: *fetchedUserGroups,
	}
	botInstance.InstanceMutex.Unlock()
}
