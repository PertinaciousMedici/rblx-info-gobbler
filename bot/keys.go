package bot

import (
	"fmt"
	"log"
	"os"

	"PanoptisMouthNew/bot/design"
	botStructures "PanoptisMouthNew/structures/bot"
	serverStructures "PanoptisMouthNew/structures/server"
)

var botToken = os.Getenv("BOT_TOKEN")
var logChannel = os.Getenv("LOG_CHANNEL")
var botPrefix = os.Getenv("BOT_PREFIX")
var homeGuild = os.Getenv("HOME_GUILD")
var rwAPIKey = os.Getenv("ROWIFI_API_KEY")
var rblxAPIKey = os.Getenv("ROBLOX_API_KEY")

func checkKeys() (map[string]string, []string, bool) {
	var environmentVariables = []botStructures.VariableType{
		{"botToken", botToken, false, false},
		{"logChannel", logChannel, true, false},
		{"botPrefix", botPrefix, false, false},
		{"homeGuild", homeGuild, true, false},
		{"rwAPIKey", rwAPIKey, false, false},
		{"rblxAPIKey", rblxAPIKey, false, false},
	}

	OkVariables, notOkVariables, situationOk := design.RunChecks(environmentVariables)

	return OkVariables, notOkVariables, situationOk
}

func panicBot(shouldPanic bool, toPanic []string, server *serverStructures.Server) {
	if shouldPanic {
		for _, missingValue := range toPanic {
			out := fmt.Sprintf("\x1b[1;31m[ERROR]:\x1b[31m Missing variable: %s\x1b[0m", missingValue)
			log.Println(out)
		}

		server.ServerChannels.ShutdownChannel <- struct{}{}
		server.WaitGroup.Wait()
	}
}
