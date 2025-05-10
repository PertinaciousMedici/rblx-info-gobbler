package server

import (
	"fmt"
	"log"

	serverStructures "PanoptisMouthNew/structures/server"
)

func exceptionLogger(server *serverStructures.Server) {
	server.WaitGroup.Add(routine)

	for {
		select {
		case rejection := <-server.ServerChannels.RejectionChannel:
			caller := rejection.Caller
			content := rejection.Content
			err := rejection.Rejection.Error()
			logOutput := fmt.Sprintf("\x1b[1;31m[%s | ERROR]:\x1b[31m %s: %s", caller, content, err)
			log.Println(logOutput)
		case _ = <-server.ServerChannels.ShutdownChannel:
			logOutput := fmt.Sprintf("\x1b[1;32m[Server | SYS]:\x1b[32m Rejection logger shutting down.\x1b[0m")
			log.Println(logOutput)
			server.WaitGroup.Done()
			server.ServerChannels.ShutdownChannel <- struct{}{}
			return
		}
	}
}

func outputLogger(server *serverStructures.Server) {
	server.WaitGroup.Add(routine)

	for {
		select {
		case output := <-server.ServerChannels.OutputChannel:
			caller := output.Caller
			content := output.Content
			logOutput := fmt.Sprintf("\x1b[1;32m[%s | SYS]:\x1b[32m %s\x1b[0m", caller, content)
			log.Println(logOutput)
		case _ = <-server.ServerChannels.ShutdownChannel:
			logOutput := fmt.Sprintf("\x1b[1;32m[Server | SYS]:\x1b[32m Output logger shutting down.\x1b[0m")
			log.Println(logOutput)
			server.WaitGroup.Done()
			server.ServerChannels.ShutdownChannel <- struct{}{}
			return
		}
	}
}

func warningLogger(server *serverStructures.Server) {
	server.WaitGroup.Add(routine)

	for {
		select {
		case warning := <-server.ServerChannels.WarningChannel:
			caller := warning.Caller
			content := warning.Content
			logOutput := fmt.Sprintf("\x1b[1;33m[%s | WARN]:\x1b[33m %s\x1b[0m", caller, content)
			log.Println(logOutput)
		case _ = <-server.ServerChannels.ShutdownChannel:
			logOutput := fmt.Sprintf("\x1b[1;32m[Server | SYS]:\x1b[32m Warning logger shutting down.\x1b[0m")
			log.Println(logOutput)
			server.WaitGroup.Done()
			server.ServerChannels.ShutdownChannel <- struct{}{}
			return
		}
	}
}
