package serverStructures

import (
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// Server
/*
 * Server is an abstraction that wraps the process of the REST API.
 * ConnectionCount is the total number of clients connected to the server.
 * ClientConnections is a map that stores pointers to the clients currently connected to the server.
 * ServerMutex ensures safe concurrency, to avoid access of already deleted clients.
 * ServerUpgrader upgrades a regular HTTP request to a Gorilla websocket.
 * ServerChannels is a wrapper for the helper functions' respective channels.
 * WaitGroup eases the management of the concurrency model.
 */
type Server struct {
	ConnectionCount   atomic.Int32
	ClientConnections map[string]*Client
	ServerMutex       sync.RWMutex
	ServerUpgrader    *websocket.Upgrader
	ServerChannels    *ServerChannels
	WaitGroup         sync.WaitGroup
}

// ServerChannels
/*
 * ServerChannels wraps the references to the channels for the helper functions.
 * ShutdownChannel sends a kill signal for graceful shutdown.
 * RejectionChannel passes ServerRejection pointers.
 * OutputChannel passes ServerOutput pointers.
 * WarningChannel passes ServerWarning pointers.
 * PostChannel passes MessagePostRequest pointers (Server -> BotInstance).
 */
type ServerChannels struct {
	ShutdownChannel  chan struct{}
	RejectionChannel chan *ServerRejection
	OutputChannel    chan *ServerOutput
	WarningChannel   chan *ServerWarning
	PostChannel      chan *MessagePostRequest
}

func (channels *ServerChannels) ThrowException(caller string, content string, rejection error) {
	channels.RejectionChannel <- &ServerRejection{
		caller,
		content,
		rejection,
	}
}

func (channels *ServerChannels) ThrowWarning(caller string, content string) {
	channels.WarningChannel <- &ServerWarning{
		caller,
		content,
	}
}

func (channels *ServerChannels) PutMessage(caller string, content string) {
	channels.OutputChannel <- &ServerOutput{
		caller,
		content,
	}
}
