package sse

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"os"
	"os/signal"
	"syscall"

	serverevents "github.com/r3labs/sse/v2"

	"github.com/sevigo/hokan/pkg/core"
)

type serverEvents struct {
	ctx    context.Context
	server *serverevents.Server
}

func New(ctx context.Context) core.ServerSideEventCreator {
	server := serverevents.New()
	server.CreateStream("messages")

	s := &serverEvents{
		server: server,
		ctx:    ctx,
	}

	go s.handleServerStop()

	return s
}

func (s *serverEvents) handleServerStop() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)

	<-interrupt
	// TODO: this is not working proper
	// probably we need to close the connection on the client side
	// if http://localhost:8081/debug/events/ is active, the server can't shut down graceful
	log.Print("[SSE] Got interrupt ...")
	s.sendStopSignal()
	s.server.Close()
}

func (s *serverEvents) sendStopSignal() {
	data := &core.ServerSideEvent{
		Message:  "stop",
		Type:     "control",
		Producer: "server",
	}

	w := new(bytes.Buffer)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic(err)
	}

	s.server.Publish("messages", &serverevents.Event{
		Data: w.Bytes(),
	})
}

func (s *serverEvents) PublishMessage(msg string) {
	data := &core.ServerSideEvent{
		Message: msg,
	}

	w := new(bytes.Buffer)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		panic(err)
	}

	s.server.Publish("messages", &serverevents.Event{
		Data: w.Bytes(),
	})
}
