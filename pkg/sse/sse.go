package sse

import (
	"bytes"
	"context"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	serverevents "github.com/r3labs/sse/v2"
	log "github.com/sirupsen/logrus"

	"github.com/sevigo/hokan/pkg/core"
)

type serverEvents struct {
	ctx           context.Context
	server        *serverevents.Server
	backupResults chan core.BackupResult
}

func New(ctx context.Context, results chan core.BackupResult) core.ServerSideEventCreator {
	server := serverevents.New()
	server.CreateStream("messages")

	s := &serverEvents{
		server:        server,
		ctx:           ctx,
		backupResults: results,
	}

	go s.handleServerStop()
	go s.debugBackupResult()

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

func (s *serverEvents) debugBackupResult() {
	for {
		select {
		case <-s.ctx.Done():
			log.Printf("sss.debugBackupResult(): event-stream canceled")
			return
		case result := <-s.backupResults:
			if result.Error != nil {
				log.Errorf("backup.Result(): [%s] %s", result.Error, result.Message)
			} else {
				log.Printf("backup.Result(): %s", result.Message)
			}
		}
	}
}
