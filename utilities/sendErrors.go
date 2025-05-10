package utilities

import (
	"fmt"

	. "PanoptisMouthNew/structures/bot"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func newErrorEmbed(errType string, description string) *discord.Embed {
	return &discord.Embed{
		Description: fmt.Sprintf("%s", description),
		Timestamp:   discord.Timestamp{},
		Color:       0xFF4C4C,
		Author: &discord.EmbedAuthor{
			Name: errType,
		},
	}
}

func sendEmbed(botInstance *BotInstance, event *gateway.MessageCreateEvent, caller string, embed *discord.Embed) {
	_, err := botInstance.BotState.SendEmbedReply(event.Message.ChannelID, event.Message.ID, *embed)
	if err != nil {
		botInstance.BotChannels.ThrowException(caller, "Could not send error embed", err)
	}
}

func InsufficientArguments(botInstance *BotInstance, event *gateway.MessageCreateEvent, expectedCount int, gotCount int, usage string) {
	expectedVsGot := fmt.Sprintf("Expected **%d**, but got **%d** arguments total.\n Command Syntax: `%s`", expectedCount, gotCount, usage)
	errorEmbed := newErrorEmbed(InsufficientArgumentsError, expectedVsGot)
	sendEmbed(botInstance, event, "InsufficientArguments", errorEmbed)
}

func MalformedArguments(botInstance *BotInstance, event *gateway.MessageCreateEvent, argumentIndex int, expectedType string, gotType string) {
	expectedVsGot := fmt.Sprintf("> Expected argument %d of type **%s** but got one of type **%s** instead.", argumentIndex+1, expectedType, gotType)
	errorEmbed := newErrorEmbed(MalformedArgumentsError, expectedVsGot)
	sendEmbed(botInstance, event, "MalformedArguments", errorEmbed)
}

func BadRequest(botInstance *BotInstance, event *gateway.MessageCreateEvent, onRequest string) {
	gotErrorOnRequest := fmt.Sprintf("Encountered an error on request to %s. %s", onRequest, TryAgainLater)
	errorEmbed := newErrorEmbed(BadRequestError, gotErrorOnRequest)
	sendEmbed(botInstance, event, "BadRequest", errorEmbed)
}

func NotFound(botInstance *BotInstance, event *gateway.MessageCreateEvent, notFoundType string) {
	whatWasNotFound := fmt.Sprintf("Request failed, **%s** not found.", notFoundType)
	errorEmbed := newErrorEmbed(NotFoundError, whatWasNotFound)
	sendEmbed(botInstance, event, "NotFound", errorEmbed)
}

func IllogicalRequest(botInstance *BotInstance, event *gateway.MessageCreateEvent, whatIsIllogical string, isArgumentRelated bool) {
	var illogicalText string

	if isArgumentRelated {
		illogicalText = fmt.Sprintf("Illogical request: %s, check your arguments.", whatIsIllogical)
	} else {
		illogicalText = fmt.Sprintf("Request failed: %s.", whatIsIllogical)
	}

	errorEmbed := newErrorEmbed(IllogicalError, illogicalText)
	sendEmbed(botInstance, event, "IllogicalRequest", errorEmbed)
}
