package server

import (
	"log"
	"net/http"
	"sync"
	"sync/atomic"

	"PanoptisMouthNew/server/methods"
	serverStructures "PanoptisMouthNew/structures/server"
	"github.com/gorilla/websocket"
)

// SERVER CONSTANTS
const serveURL = ""
const servePort = ":6312"

// NUMERIC CONSTANTS
const routine int = 1

func RunServer() *serverStructures.Server {
	newServer := serverStructures.Server{
		ConnectionCount:   atomic.Int32{},
		ClientConnections: make(map[string]*serverStructures.Client, 64),
		ServerMutex:       sync.RWMutex{},
		ServerUpgrader:    &websocket.Upgrader{ReadBufferSize: 4096, WriteBufferSize: 4096},
		ServerChannels: &serverStructures.ServerChannels{
			ShutdownChannel:  make(chan struct{}, 3),
			RejectionChannel: make(chan *serverStructures.ServerRejection, 96),
			OutputChannel:    make(chan *serverStructures.ServerOutput, 96),
			WarningChannel:   make(chan *serverStructures.ServerWarning, 96),
		},
		WaitGroup: sync.WaitGroup{},
	}

	go exceptionLogger(&newServer)
	go warningLogger(&newServer)
	go outputLogger(&newServer)

	func(server *serverStructures.Server) {
		newServer.ServerChannels.PutMessage("Server", "Started rejection logger.")
		newServer.ServerChannels.PutMessage("Server", "Started warning logger.")
		newServer.ServerChannels.PutMessage("Server", "Started output logger.")
	}(&newServer)

	go serverListener(&newServer)

	return &newServer
}

func serverListener(server *serverStructures.Server) {
	server.ServerChannels.PutMessage("Server", "Server started successfully.")

	ServerManager := methods.ServerManager{
		ServerInstance: server,
	}

	http.HandleFunc("/handshake", ServerManager.HandshakeHandler)
	http.HandleFunc("/post", ServerManager.PostMessage)
	err := http.ListenAndServe(serveURL+servePort, nil)
	if err != nil {
		server.ServerChannels.ThrowException("Server", "Failed to serve and listen", err)
		server.ServerChannels.ShutdownChannel <- struct{}{}
		server.WaitGroup.Wait()
		log.Panic()
		return
	}
}
