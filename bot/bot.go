package bot

import (
	"context"
	"log"
	"net/http"
	"sync"
	"time"

	"PanoptisMouthNew/bot/handlers"
	. "PanoptisMouthNew/structures/bot"
	"PanoptisMouthNew/structures/server"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

// NUMERIC CONSTANTS
const routine int = 1

func RunBot(server *serverStructures.Server) {
	globalDefines, toPanic, canRun := checkKeys()
	panicBot(!canRun, toPanic, server)
	time.Sleep(1 * time.Second)

	botState := state.New("Bot " + botToken)

	{
		botState.AddIntents(gateway.IntentGuilds)
		botState.AddIntents(gateway.IntentGuildMessages)
		botState.AddIntents(gateway.IntentDirectMessages)
		botState.AddIntents(gateway.IntentMessageContent)
	}

	botInstance := BotInstance{
		BotState: botState,
		HttpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
		CommandHooks:        make(map[string]Command),
		CategorisedCommands: nil,
		BotStorage: &BotStorage{
			StoredGroups:  make(map[string]*FetchedUserGroupMemberships),
			StoredFriends: make(map[string]*FetchedRblxUserFriends),
			StoredMutualG: make(map[string]*FetchedUserGroupMemberships),
			StoredMutualF: make(map[string]*FetchedRblxUserFriends),
		},
		BotChannels: &BotChannels{
			ShutdownChannel:  make(chan struct{}),
			RejectionChannel: make(chan *BotRejection, 64),
			OutputChannel:    make(chan *BotOutput, 64),
			WarningChannel:   make(chan *BotWarning, 64),
		},
		ServerInstance: server,
		InstanceMutex:  sync.RWMutex{},
		WaitGroup:      sync.WaitGroup{},
		GlobalDefines:  make(map[string]string),
	}

	go exceptionLogger(&botInstance)
	go warningLogger(&botInstance)
	go outputLogger(&botInstance)

	for key, value := range globalDefines {
		botInstance.GlobalDefines[key] = value
	}

	botController := handlers.BotManager{
		BotInstance: &botInstance,
	}

	SetHooks(&botInstance)

	botInstance.InstanceMutex.Lock()
	botInstance.BotState.AddHandler(botController.HandleInteractions)
	botInstance.BotState.AddHandler(botController.HandleMessages)
	botInstance.BotState.AddHandler(botController.ReadyEvent)
	botInstance.InstanceMutex.Unlock()

	err := botInstance.BotState.Open(context.Background())
	if err != nil {
		log.Panic()
	}
}
