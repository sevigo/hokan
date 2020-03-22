package sse

import (
	"context"

	serverevents "github.com/r3labs/sse"

	"github.com/sevigo/hokan/pkg/core"
)

type serverEvents struct {
	ctx    context.Context
	server *serverevents.Server
}

func New(ctx context.Context) core.ServerSideEventCreater {
	server := serverevents.New()

	return &serverEvents{
		ctx:    ctx,
		server: server,
	}
}

func (s *serverEvents) PublishMessage(msg string) {
	s.server.Publish("messages", &serverevents.Event{
		Data: []byte(msg),
	})
}
