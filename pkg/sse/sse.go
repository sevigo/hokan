package sse

import (
	serverevents "github.com/r3labs/sse"

	"github.com/sevigo/hokan/pkg/core"
)

type serverEvents struct {
	server *serverevents.Server
}

func New() core.ServerSideEventCreator {
	server := serverevents.New()
	server.CreateStream("messages")

	return &serverEvents{
		server: server,
	}
}

func (s *serverEvents) PublishMessage(msg string) {
	s.server.Publish("messages", &serverevents.Event{
		Data: []byte(msg),
	})
}
