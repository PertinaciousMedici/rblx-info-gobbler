package bot

import (
	"fmt"
	"log"

	botStructures "PanoptisMouthNew/structures/bot"
)

func exceptionLogger(bot *botStructures.BotInstance) {
	bot.WaitGroup.Add(routine)
	bot.BotChannels.PutMessage("Bot", "Started rejection logger.")

	for {
		select {
		case rejection := <-bot.BotChannels.RejectionChannel:
			caller := rejection.Caller
			content := rejection.Content
			err := rejection.Rejection.Error()
			logOutput := fmt.Sprintf("\x1b[1;31m[%s | ERROR]:\x1b[31m %s: %s", caller, content, err)
			log.Println(logOutput)
		case _ = <-bot.BotChannels.ShutdownChannel:
			logOutput := fmt.Sprintf("\x1b[1;32m[Bot | SYS]:\x1b[32m Rejection logger shutting down.\x1b[0m")
			log.Println(logOutput)
			bot.WaitGroup.Done()
			bot.BotChannels.ShutdownChannel <- struct{}{}
			return
		}
	}
}

func outputLogger(bot *botStructures.BotInstance) {
	bot.WaitGroup.Add(routine)
	bot.BotChannels.PutMessage("Bot", "Started output logger.")

	for {
		select {
		case output := <-bot.BotChannels.OutputChannel:
			caller := output.Caller
			content := output.Content
			logOutput := fmt.Sprintf("\x1b[1;32m[%s | SYS]:\x1b[32m %s\x1b[0m", caller, content)
			log.Println(logOutput)
		case _ = <-bot.BotChannels.ShutdownChannel:
			logOutput := fmt.Sprintf("\x1b[1;32m[Bot | SYS]:\x1b[32m Output logger shutting down.\x1b[0m")
			log.Println(logOutput)
			bot.WaitGroup.Done()
			bot.BotChannels.ShutdownChannel <- struct{}{}
			return
		}
	}
}

func warningLogger(bot *botStructures.BotInstance) {
	bot.WaitGroup.Add(routine)
	bot.BotChannels.PutMessage("Bot", "Started warning logger.")

	for {
		select {
		case warning := <-bot.BotChannels.WarningChannel:
			caller := warning.Caller
			content := warning.Content
			logOutput := fmt.Sprintf("\x1b[1;33m[%s | WARN]:\x1b[33m %s\x1b[0m", caller, content)
			log.Println(logOutput)
		case _ = <-bot.BotChannels.ShutdownChannel:
			logOutput := fmt.Sprintf("\x1b[1;32m[Bot | SYS]:\x1b[32m Warning logger shutting down.\x1b[0m")
			log.Println(logOutput)
			bot.WaitGroup.Done()
			bot.BotChannels.ShutdownChannel <- struct{}{}
			return
		}
	}
}
