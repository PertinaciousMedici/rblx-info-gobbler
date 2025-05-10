package commands

import (
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/utilities"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func Help(botInstance *BotInstance, event *gateway.MessageCreateEvent, args []string) {
	hasTarget := len(args) >= 1

	if hasTarget {
		cmd, ok := botInstance.CommandHooks[args[0]]

		if !ok {
			utilities.IllogicalRequest(botInstance, event, "Command does not exist", true)
			return
		}

		utilities.SendCommandUsage(botInstance, event, cmd.Name, cmd.Usage, cmd.Description, cmd.Category)
		return
	}

	if botInstance.CategorisedCommands == nil {
		categorisedCommands := make(map[string][]*Command)
		for idx := range botInstance.CommandHooks {
			command := botInstance.CommandHooks[idx]
			cmdCategory := command.Category
			categorisedCommands[cmdCategory] = append(categorisedCommands[cmdCategory], &command)
		}

		botInstance.InstanceMutex.Lock()
		botInstance.CategorisedCommands = categorisedCommands
		botInstance.InstanceMutex.Unlock()

		if len(categorisedCommands) == 0 {
			return
		}
	}

	categorisedCommands := botInstance.CategorisedCommands
	utilities.SendHelpEmbed(botInstance, event, categorisedCommands)
}
