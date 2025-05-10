package utilities

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"PanoptisMouthNew/bot/connection"
	. "PanoptisMouthNew/structures/bot"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func SendHelpEmbed(botInstance *BotInstance, event *gateway.MessageCreateEvent, mappedCommands map[string][]*Command) {
	var embedFields []discord.EmbedField

	var categoryNames []string
	for categoryName := range mappedCommands {
		categoryNames = append(categoryNames, categoryName)
	}
	sort.Strings(categoryNames)

	for _, categoryName := range categoryNames {
		commandCategory := mappedCommands[categoryName]
		var fieldData []string
		for _, command := range commandCategory {
			fieldData = append(fieldData, fmt.Sprintf("`%s`", command.Name))
		}

		newField := discord.EmbedField{
			Name:  categoryName,
			Value: strings.Join(fieldData, ", "),
		}

		embedFields = append(embedFields, newField)
	}

	newEmbed := discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: "Commands",
		},
		Color:     0x808080,
		Timestamp: discord.Timestamp{},
		Fields:    embedFields,
	}

	sendEmbed(botInstance, event, "SendHelpEmbed", &newEmbed)
}

func SendCommandUsage(botInstance *BotInstance, event *gateway.MessageCreateEvent, command string, usage string, description string, category string) {
	newEmbed := &discord.Embed{
		Author: &discord.EmbedAuthor{
			Name: "Help | Command",
		},
		Color: 0x808080,
		Description: fmt.Sprintf(
			"### [%s]\n\n\n**[Category]:** `%s`\n**[Syntax]:** `%s`\n**[Description]:** `%s`.",
			command,
			category,
			usage,
			description,
		),
		Timestamp: discord.Timestamp{},
	}
	sendEmbed(botInstance, event, "SendCommandUsage", newEmbed)
}

func SendRobloxUser(botInstance *BotInstance, event *gateway.MessageCreateEvent, fetchedUser *FetchedUserById) {
	var hasPremium = "No"

	if fetchedUser.Premium {
		hasPremium = "Yes"
	}

	parsedTime, err := time.Parse(time.RFC3339, fetchedUser.CreationDate)
	if err != nil {
		botInstance.BotChannels.ThrowException("SendRobloxUser", "Failed to parse time", err)
	}
	formattedTime := parsedTime.Format("January 2, 2006 at 3:04 PM (MST)")

	responseEmbed := discord.Embed{
		Title:     fmt.Sprintf("%s's Profile", fetchedUser.Username),
		URL:       fmt.Sprintf("https://www.roblox.com/users/%s/profile", fetchedUser.UserID),
		Timestamp: discord.Timestamp{},
		Color:     0x808080,
		Footer: &discord.EmbedFooter{
			Text: "Rule, Alatheia!",
			Icon: AlatheiaEmoji,
		},
		Thumbnail: &discord.EmbedThumbnail{
			URL: fetchedUser.AvatarURL,
		},
		Fields: []discord.EmbedField{
			{
				Name:  "Roblox ID",
				Value: fetchedUser.UserID,
			},
			{
				Name:  "Username",
				Value: fetchedUser.Username,
			},
			{
				Name:  "Display Name",
				Value: fetchedUser.DisplayName,
			},
			{
				Name:  "Creation Date",
				Value: formattedTime,
			},
			{
				Name:  "Friend Count",
				Value: fmt.Sprintf("%d", fetchedUser.FriendCount),
			},
			{
				Name:  "Description",
				Value: fetchedUser.Description,
			},
			{
				Name:  "Premium?",
				Value: hasPremium,
			},
		},
	}

	sendEmbed(botInstance, event, "SendRobloxUser", &responseEmbed)
}

func SendDiscordUser(botInstance *BotInstance, event *gateway.MessageCreateEvent, fetchedUser *DiscFetchedUserById) {
	if fetchedUser == nil {
		return
	}

	hasURL := fetchedUser.Avatar != ""
	var avatarURL string

	if !hasURL {
		idAsUint, _ := strconv.Atoi(fetchedUser.UserID)
		avatarURL = fmt.Sprintf(DefaultDiscPfp, idAsUint%5)
	} else {
		avatarURL = fmt.Sprintf(connection.DsAvatarSearchGet, fetchedUser.UserID, fetchedUser.Avatar)
	}

	creationDate := deriveCreationDate(fetchedUser.UserID)

	userFlags, hasFlags := checkFlags(fetchedUser.Flags)
	flagsTxt := strings.Join(userFlags, ",\n")

	if !hasFlags {
		flagsTxt = "None"
	}

	responseEmbed := discord.Embed{
		Timestamp: discord.Timestamp{},
		Color:     0x808080,
		Footer: &discord.EmbedFooter{
			Text: "Rule, Alatheia!",
			Icon: AlatheiaEmoji,
		},
		Thumbnail: &discord.EmbedThumbnail{
			URL: avatarURL,
		},
		Author: &discord.EmbedAuthor{
			Name: fetchedUser.Username,
			Icon: avatarURL,
		},
		Fields: []discord.EmbedField{
			{
				Name:  "User ID",
				Value: fetchedUser.UserID,
			},
			{
				Name:  "Creation Date",
				Value: creationDate,
			},
			{
				Name:  "Flags",
				Value: flagsTxt,
			},
		},
	}

	sendEmbed(botInstance, event, "SendDiscordUser", &responseEmbed)
}
