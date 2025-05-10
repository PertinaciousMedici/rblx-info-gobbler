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

func HandleGroupIteration(botInstance *BotInstance, event *gateway.InteractionCreateEvent) {
	data := event.Data.(*discord.ButtonInteraction)
	customId := string(data.CustomID)

	var currentPage int
	var newPage = 0

	if strings.HasPrefix(customId, "groups_first") {
		newPage = 1
	} else if strings.HasPrefix(customId, "groups_prev") {
		_, err := fmt.Sscanf(customId, "groups_prev_%d", &currentPage)
		if err != nil {
			botInstance.BotChannels.ThrowException("HandleGroupIteration", "Failed to parse custom ID", err)
			return
		}
		newPage = currentPage - 1
	} else if strings.HasPrefix(customId, "groups_next") {
		_, err := fmt.Sscanf(customId, "groups_next_%d", &currentPage)
		if err != nil {
			botInstance.BotChannels.ThrowException("HandleGroupIteration", "Failed to parse custom ID", err)
			return
		}
		newPage = currentPage + 1
	}

	toExtract := event.Message.Embeds[0].Author.Name
	username := strings.TrimSuffix(toExtract, "'s Groups")

	botInstance.InstanceMutex.RLock()
	userGroups, ok := botInstance.BotStorage.StoredGroups[username]
	if !ok {
		return
	}
	botInstance.InstanceMutex.RUnlock()

	maxPages := utilities.CalculateMaxPages(design.MaximumFieldsPerPageGroups, len(userGroups.Groups))

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
		newPage,
		&userGroups.Groups,
		fieldBuilder,
		"",
	)

	utilities.EditIterator(
		botInstance,
		event,
		"groups",
		newPage,
		&groupsEmbed,
		"GroupIteration",
	)
}
