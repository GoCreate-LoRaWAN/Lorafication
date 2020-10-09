package server

import "net/http"

func (s *Server) Notify(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("test"))
}
