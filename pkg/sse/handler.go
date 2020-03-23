package sse

import "net/http"

func (s *serverEvents) Handler(w http.ResponseWriter, r *http.Request) {
	s.server.HTTPHandler(w, r)
}
