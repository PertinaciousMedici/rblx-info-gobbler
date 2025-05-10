package handlers

import (
	"fmt"

	"github.com/diamondburned/arikawa/v3/gateway"
)

func (instance *BotManager) ReadyEvent(event *gateway.ReadyEvent) {
	_ = event
	
	botInstance := instance.BotInstance
	self, err := botInstance.BotState.Me()

	if err != nil {
		botInstance.BotChannels.ThrowException("ReadyEvent", "Failed to fetch self", err)
	}

	out := fmt.Sprintf("Launched bot instance as %s", self.Username)
	botInstance.BotChannels.PutMessage("ReadyEvent", out)
}
