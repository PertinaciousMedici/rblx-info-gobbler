package interactions

import (
	"fmt"
	"strings"

	"PanoptisMouthNew/bot/connection"
	"PanoptisMouthNew/bot/design"
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func HandleMGroupIteration(botInstance *BotInstance, event *gateway.InteractionCreateEvent) {
	data := event.Data.(*discord.ButtonInteraction)
	customId := string(data.CustomID)

	var currentPage int
	var newPage = 0

	if strings.HasPrefix(customId, "mutualg_first") {
		newPage = 1
	} else if strings.HasPrefix(customId, "mutualg_prev") {
		_, err := fmt.Sscanf(customId, "mutualg_prev_%d", &currentPage)
		if err != nil {
			botInstance.BotChannels.ThrowException("HandleMGroupIteration", "Failed to parse custom ID", err)
			return
		}
		newPage = currentPage - 1
	} else if strings.HasPrefix(customId, "mutualg_next") {
		_, err := fmt.Sscanf(customId, "mutualg_next_%d", &currentPage)
		if err != nil {
			botInstance.BotChannels.ThrowException("HandleMGroupIteration", "Failed to parse custom ID", err)
			return
		}
		newPage = currentPage + 1
	}

	toExtract := event.Message.Embeds[0].Fields[0].Value
	toExtract = toExtract[1 : len(toExtract)-1]

	rawUsernames := event.Message.Embeds[0].Description

	botInstance.InstanceMutex.RLock()
	userMGroups, ok := botInstance.BotStorage.StoredMutualG[toExtract]
	if !ok {
		return
	}
	botInstance.InstanceMutex.RUnlock()

	username1, username2 := utilities.ExtractUsernames(rawUsernames)

	maxPages := utilities.CalculateMaxPages(design.MaximumFieldsPerPageGroups, len(userMGroups.Groups))

	if newPage == 0 {
		newPage = maxPages
	}

	if newPage < 1 {
		newPage = 1
	} else if newPage > maxPages {
		newPage = maxPages
	}

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
		toExtract,
		design.MaximumFieldsPerPageGroups,
		maxPages,
		newPage,
		&userMGroups.Groups,
		fieldBuilder,
		fmt.Sprintf("**%s** & **%s**", username1, username2),
	)

	utilities.EditIterator(
		botInstance,
		event,
		"mutualg",
		newPage,
		&groupsEmbed,
		"MutualGroupsIteration",
	)
}
