package handlers

import (
	"fmt"
	"regexp"
	"strings"

	"PanoptisMouthNew/bot/interactions"
	"github.com/diamondburned/arikawa/v3/discord"
	"github.com/diamondburned/arikawa/v3/gateway"
)

var categories = []string{"mutualf", "mutualg", "friends", "groups"}

func (instance *BotManager) HandleInteractions(event *gateway.InteractionCreateEvent) {
	typeOfInteraction := event.Data.InteractionType()

	switch typeOfInteraction {
	case discord.ComponentInteractionType:
		instance.HandleIterativeEmbeds(event)
		break
	default:
		return
	}
}

func (instance *BotManager) HandleIterativeEmbeds(event *gateway.InteractionCreateEvent) {
	eventData := event.Data.(*discord.ButtonInteraction)
	customId := string(eventData.CustomID)

	joinedCategories := strings.Join(categories, "|")
	categoryPattern := fmt.Sprintf("(%s)", joinedCategories)

	regex := regexp.MustCompile("^" + categoryPattern + "_")
	matches := regex.FindStringSubmatch(customId)

	if matches != nil {
		category := matches[1]

		switch category {
		case "friends":
			interactions.HandleFriendIteration(instance.BotInstance, event)
			break
		case "groups":
			interactions.HandleGroupIteration(instance.BotInstance, event)
			break
		case "mutualf":
			interactions.HandleMFriendIteration(instance.BotInstance, event)
			break
		case "mutualg":
			interactions.HandleMGroupIteration(instance.BotInstance, event)
			return
		default:
			return
		}
	}
}
