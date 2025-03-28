package main

import "net/http"

type Server struct {
	store *KVStore
	mux   *http.ServeMux
}

func NewServer(store *KVStore) *Server {
	s := &Server{
		store: store,
		mux:   http.NewServeMux(),
	}
	s.routes()
	return s
}

func (s *Server) routes() {
	s.mux.HandleFunc("/set", s.SetHandler)
	s.mux.HandleFunc("/get", s.GetHandler)
	s.mux.HandleFunc("/delete", s.DeleteHandler)
	s.mux.HandleFunc("/keys", s.KeysHandler)
}
