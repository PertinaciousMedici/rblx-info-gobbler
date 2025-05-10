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

func HandleFriendIteration(botInstance *BotInstance, event *gateway.InteractionCreateEvent) {
	data := event.Data.(*discord.ButtonInteraction)
	customId := string(data.CustomID)

	var currentPage int
	var newPage = 0

	if strings.HasPrefix(customId, "friends_first") {
		newPage = 1
	} else if strings.HasPrefix(customId, "friends_prev") {
		_, err := fmt.Sscanf(customId, "friends_prev_%d", &currentPage)
		if err != nil {
			botInstance.BotChannels.ThrowException("HandleFriendIteration", "Failed to parse custom ID", err)
			return
		}
		newPage = currentPage - 1
	} else if strings.HasPrefix(customId, "friends_next") {
		_, err := fmt.Sscanf(customId, "friends_next_%d", &currentPage)
		if err != nil {
			botInstance.BotChannels.ThrowException("HandleFriendIteration", "Failed to parse custom ID", err)
			return
		}
		newPage = currentPage + 1
	}

	toExtract := event.Message.Embeds[0].Author.Name
	username := strings.TrimSuffix(toExtract, "'s Friends")

	botInstance.InstanceMutex.RLock()
	userFriends, ok := botInstance.BotStorage.StoredFriends[username]
	if !ok {
		return
	}
	botInstance.InstanceMutex.RUnlock()

	maxPages := utilities.CalculateMaxPages(design.MaximumFieldsPerPageFriends, len(userFriends.Friends))

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
		fmt.Sprintf("%s's Friends", username),
		"",
		design.MaximumFieldsPerPageFriends,
		maxPages,
		newPage,
		&userFriends.Friends,
		fieldBuilder,
		"",
	)

	utilities.EditIterator(
		botInstance,
		event,
		"friends",
		newPage,
		&friendsEmbed,
		"FriendIteration",
	)
}
