package core

import "net/http"

type ServerSideEventCreator interface {
	Handler(w http.ResponseWriter, r *http.Request)
	PublishMessage(string)
}

type ServerSideEvent struct {
	Message  string `json:"message"`
	Producer string `json:"producer"`
	Type     string `json:"type"`
	Data     string `json:"data"`
}
