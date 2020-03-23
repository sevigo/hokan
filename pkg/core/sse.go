package core

import "net/http"

type ServerSideEventCreator interface {
	Handler(w http.ResponseWriter, r *http.Request)
	PublishMessage(string)
}
