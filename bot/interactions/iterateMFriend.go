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

func HandleMFriendIteration(botInstance *BotInstance, event *gateway.InteractionCreateEvent) {
	data := event.Data.(*discord.ButtonInteraction)
	customId := string(data.CustomID)

	var currentPage int
	var newPage = 0

	if strings.HasPrefix(customId, "mutualf_first") {
		newPage = 1
	} else if strings.HasPrefix(customId, "mutualf_prev") {
		_, err := fmt.Sscanf(customId, "mutualf_prev_%d", &currentPage)
		if err != nil {
			botInstance.BotChannels.ThrowException("HandleMFriendIteration", "Failed to parse custom ID", err)
			return
		}
		newPage = currentPage - 1
	} else if strings.HasPrefix(customId, "mutualf_next") {
		_, err := fmt.Sscanf(customId, "mutualf_next_%d", &currentPage)
		if err != nil {
			botInstance.BotChannels.ThrowException("HandleMFriendIteration", "Failed to parse custom ID", err)
			return
		}
		newPage = currentPage + 1
	}

	toExtract := event.Message.Embeds[0].Fields[0].Value
	toExtract = toExtract[1 : len(toExtract)-1]

	rawUsernames := event.Message.Embeds[0].Description

	botInstance.InstanceMutex.RLock()
	userMFriends, ok := botInstance.BotStorage.StoredMutualF[toExtract]
	if !ok {
		return
	}
	botInstance.InstanceMutex.RUnlock()

	username1, username2 := utilities.ExtractUsernames(rawUsernames)

	maxPages := utilities.CalculateMaxPages(design.MaximumFieldsPerPageFriends, len(userMFriends.Friends))

	if newPage == 0 {
		newPage = maxPages
	}

	if newPage < 1 {
		newPage = 1
	} else if newPage > maxPages {
		newPage = maxPages
	}

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
		toExtract,
		design.MaximumFieldsPerPageFriends,
		maxPages,
		newPage,
		&userMFriends.Friends,
		fieldBuilder,
		fmt.Sprintf("**%s** & **%s**", username1, username2),
	)

	utilities.EditIterator(
		botInstance,
		event,
		"mutualf",
		newPage,
		&friendsEmbed,
		"MutualFriendsIteration",
	)
}
