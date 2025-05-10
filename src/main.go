package main

import (
	"PanoptisMouthNew/bot"
	"PanoptisMouthNew/server"
)

func main() {
	serverInstance := server.RunServer()
	go bot.RunBot(serverInstance)
	select {}
}
