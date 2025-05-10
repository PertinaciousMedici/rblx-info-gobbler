package commands

import (
	"time"

	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func SoftArchive(botInstance *BotInstance, event *gateway.MessageCreateEvent, args []string) {
	_ = args
	var toBulkDelete []discord.MessageID
	var beforeID = event.Message.ID
	var totalDeleted int
	const messageLimit = 100

	allowedUsers := []discord.UserID{1308785655597498484, 475688942604124193, 590477605564841994}
	isAllowed := false

	for _, user := range allowedUsers {
		if user == event.Message.Author.ID {
			isAllowed = true
		}
	}

	if !isAllowed {
		utilities.IllogicalRequest(botInstance, event, "Not allowed", false)
		return
	}

	for {
		messages, err := botInstance.BotState.MessagesBefore(event.Message.ChannelID, beforeID, messageLimit)

		if err != nil {
			utilities.BadRequest(botInstance, event, "Fetch Messages")
			return
		}

		if len(messages) == 0 {
			break
		}

		for _, message := range messages {
			if message.Pinned {
				continue
			}

			if time.Since(message.Timestamp.Time()) < 14*24*time.Hour {
				toBulkDelete = append(toBulkDelete, message.ID)
				totalDeleted++
			}
		}

		beforeID = messages[len(messages)-1].ID
	}

	if len(toBulkDelete) >= 2 {
		err := botInstance.BotState.DeleteMessages(event.Message.ChannelID, toBulkDelete, "SOFT ARCHIVE")
		if err != nil {
			utilities.BadRequest(botInstance, event, "Delete Messages")
			return
		}
	} else {
		for _, messageID := range toBulkDelete {
			err := botInstance.BotState.DeleteMessage(event.Message.ChannelID, messageID, "SOFT ARCHIVE")
			if err != nil {
				utilities.BadRequest(botInstance, event, "Delete Messages")
				return
			}
		}
	}
}
