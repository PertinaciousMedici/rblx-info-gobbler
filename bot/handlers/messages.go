package handlers

import (
	"strings"

	"PanoptisMouthNew/utilities"
	"github.com/diamondburned/arikawa/v3/gateway"
)

func (instance *BotManager) HandleMessages(event *gateway.MessageCreateEvent) {
	botPrefix := instance.BotInstance.GlobalDefines["botPrefix"]
	content := event.Message.Content
	hasValidLength := len(content) >= len(botPrefix)
	startsWithPrefix := strings.HasPrefix(content, botPrefix)

	if startsWithPrefix {
		content = content[len(botPrefix):]
	}

	if hasValidLength && startsWithPrefix {
		var args []string
		var currentArg string

		for _, charRune := range content {
			if charRune == ' ' {
				if currentArg != "" {
					args = append(args, currentArg)
					currentArg = ""
				}
			} else {
				currentArg += string(charRune)
			}
		}

		if currentArg != "" {
			args = append(args, currentArg)
		}

		if len(args) == 0 {
			return
		}

		instance.HandleTextCommands(event, args)
	}
}

func (instance *BotManager) HandleTextCommands(event *gateway.MessageCreateEvent, args []string) {
	botInstance := instance.BotInstance

	command := args[0]
	args = args[1:]

	botInstance.InstanceMutex.RLock()
	cmd, ok := botInstance.CommandHooks[command]
	if ok {
		if len(args) >= cmd.ExpectedArguments {
			go cmd.Execute(botInstance, event, args)
		} else {
			utilities.InsufficientArguments(botInstance, event, cmd.ExpectedArguments, len(args), cmd.Usage)
		}
	}
	botInstance.InstanceMutex.RUnlock()
}
