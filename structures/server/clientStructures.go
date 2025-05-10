package serverStructures

import (
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// Client
/*
 * Client is an abstraction for a websocket connection and contains all information pertaining to it.
 * Discriminator is a user-chosen username specified at handshake time to discern clients.
 * Connection is the Gorilla websocket connection, upgraded upon access.
 * ClientMutex ensures no unexpected end of connection on the serverside.
 * WsLatency is a regularly updated latency measurement sent alongside the handshake.
 * Active determines whether a client is active, and is updated upon critical exceptions.
 */
type Client struct {
	Discriminator string
	Connection    *websocket.Conn
	ClientMutex   *sync.RWMutex
	WsLatency     uint64
	Active        atomic.Bool
}

// MessagePostRequest
/*
 * MessageURL is the link to the message flagged by the Eye.
 * MessageContent is the text content flagged by the Eye.
 * AuthorName is the name of the author of the message.
 * AuthorID is the user snowflake of the author of the message flagged by the Eye.
 * GuildName is the name of the guild where the message was sent.
 * GuildID is the guild snowflake of the server where the message was sent.
 */
type MessagePostRequest struct {
	MessageURL     string `json:"url"`
	MessageContent string `json:"messageContent"`
	AuthorName     string `json:"authorName"`
	AuthorID       string `json:"authorId"`
	GuildName      string `json:"guildName"`
	GuildID        string `json:"guildId"`
}

// HandshakePayload
/*
 * ClientUsername is a user-defined discriminator sent upon handshake.
 * CurrentDate is the UNIX timestamp for websocket latency calculation.
 */
type HandshakePayload struct {
	ClientUsername string `json:"clientUsername"`
	CurrentDate    uint64 `json:"currentDate"`
}
