package utilities

import (
	"fmt"

	. "PanoptisMouthNew/structures/bot"
	"github.com/diamondburned/arikawa/v3/api"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func BuildPaginatedEmbed[Type any](
	title string,
	requestID string,
	itemsPerPage int,
	maxPage int,
	currentPage int,
	data *[]*Type,
	formatField func(*Type) discord.EmbedField,
	extraDescription string,
) discord.Embed {
	var untypedArray []interface{}

	for _, element := range *data {
		untypedArray = append(untypedArray, element)
	}

	paginatedArray := PaginateArray(untypedArray, itemsPerPage, currentPage)

	var fields []discord.EmbedField

	for _, element := range paginatedArray {
		field := formatField(element.(*Type))
		fields = append(fields, field)
	}

	returnEmbed := discord.Embed{
		Timestamp:   discord.NowTimestamp(),
		Color:       0xF4860A,
		Description: "*Press the **buttons** below to iterate*\n",
		Footer: &discord.EmbedFooter{
			Text: fmt.Sprintf("Page %d/%d", currentPage, maxPage),
		},
		Author: &discord.EmbedAuthor{
			Name: title,
		},
		Fields: fields,
	}

	if requestID != "" {
		returnEmbed.Fields = append([]discord.EmbedField{
			{
				Value: fmt.Sprintf("`%s`", requestID),
			},
		}, returnEmbed.Fields...)
	}

	if extraDescription != "" {
		returnEmbed.Description = fmt.Sprintf("%s%s\n", returnEmbed.Description, extraDescription)
	}

	return returnEmbed
}

func SendIterator(
	botInstance *BotInstance,
	channelID discord.ChannelID,
	identifier string,
	currentPage int,
	embed *discord.Embed,
	label string,
) {
	complexMessage := BuildIterator(identifier, currentPage, embed)
	_, err := botInstance.BotState.SendMessageComplex(channelID, *complexMessage)
	if err != nil {
		botInstance.BotChannels.ThrowException(label, "Could not send iterator", err)
	}
}

func EditIterator(
	botInstance *BotInstance,
	event *gateway.InteractionCreateEvent,
	identifier string,
	currentPage int,
	embed *discord.Embed,
	label string,
) {
	complexMessage := BuildIterator(identifier, currentPage, embed)
	response := api.InteractionResponse{
		Type: api.UpdateMessage,
		Data: &api.InteractionResponseData{
			Embeds:     &complexMessage.Embeds,
			Components: &complexMessage.Components,
		},
	}

	err := botInstance.BotState.RespondInteraction(event.ID, event.Token, response)
	if err != nil {
		botInstance.BotChannels.ThrowException(label, "Could not edit iterator", err)
	}
}
