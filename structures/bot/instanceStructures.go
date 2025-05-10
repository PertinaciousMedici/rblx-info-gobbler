package botStructures

import (
	"net/http"
	"sync"

	serverStructures "PanoptisMouthNew/structures/server"
	"github.com/diamondburned/arikawa/v3/gateway"
	"github.com/diamondburned/arikawa/v3/state"
)

// Command
/*
 * Name is the identifier of the command, through which it is run.
 * Usage is the arguments necessary and formatting patterns for usage.
 * Description is a brief summary of the command's utility.
 * Category is the category of commands it belongs to.
 * ExpectedArguments is the minimum number of arguments to be passed down.
 * Execute is the function of the command itself.
 */
type Command struct {
	Name              string
	Usage             string
	Description       string
	Category          string
	ExpectedArguments int
	Execute           func(instance *BotInstance, event *gateway.MessageCreateEvent, args []string)
}

// BotInstance
/*
 * BotState is a pointer to a state.State, which allows communication with Discord.
 * HttpClient is the unified source of requests to the APIs.
 * CommandHooks is a map of all text chat commands functions.
 * BotStorage is for short-term iteration of embeds, a cache of sorts.
 * BotChannels is a wrapper for the helper functions' respective channels.
 * ServerInstance is a reference to the web server for ease of access.
 * InstanceMutex serves to ease the multithreading model, protects shared resources.
 * WaitGroup serves to ease the multithreading model through graceful shutdown.
 */
type BotInstance struct {
	BotState            *state.State
	HttpClient          *http.Client
	CommandHooks        map[string]Command
	CategorisedCommands map[string][]*Command
	BotStorage          *BotStorage
	BotChannels         *BotChannels
	ServerInstance      *serverStructures.Server
	InstanceMutex       sync.RWMutex
	WaitGroup           sync.WaitGroup
	GlobalDefines       map[string]string
}

// BotStorage
/*
 * StoredGroups holds references to a user's Roblox groups mapped to their username.
 * StoredFriends holds references to a user's Roblox friends mapped to their username.
 * StoredMutualG holds references to a user's Roblox groups mapped to the message request ID.
 * StoredMutualF holds references to a user's Roblox friends mapped to the message request ID.
 */
type BotStorage struct {
	StoredGroups  map[string]*FetchedUserGroupMemberships
	StoredFriends map[string]*FetchedRblxUserFriends
	StoredMutualG map[string]*FetchedUserGroupMemberships
	StoredMutualF map[string]*FetchedRblxUserFriends
}

// BotChannels
/*
 * BotChannels wraps the references to the channels for the helper functions.
 * ShutdownChannel sends a kill signal for graceful shutdown.
 * RejectionChannel passes BotRejection pointers.
 * OutputChannel passes BotOutput pointers.
 * WarningChannel passes BotWarning pointers.
 */
type BotChannels struct {
	ShutdownChannel  chan struct{}
	RejectionChannel chan *BotRejection
	OutputChannel    chan *BotOutput
	WarningChannel   chan *BotWarning
}

func (channels *BotChannels) ThrowException(caller string, content string, rejection error) {
	channels.RejectionChannel <- &BotRejection{
		caller,
		content,
		rejection,
	}
}

func (channels *BotChannels) ThrowWarning(caller string, content string) {
	channels.WarningChannel <- &BotWarning{
		caller,
		content,
	}
}

func (channels *BotChannels) PutMessage(caller string, content string) {
	channels.OutputChannel <- &BotOutput{
		caller,
		content,
	}
}
