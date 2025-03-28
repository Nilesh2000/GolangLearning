package main

import (
	"log"
	"net/http"
)

func main() {
	store := NewKVStore()
	s := NewServer(store)
	log.Fatal(http.ListenAndServe(":8080", s.mux))
}
