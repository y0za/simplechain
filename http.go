package main

import (
	"encoding/json"
	"net/http"
)

// Env inject dependencies to http handler
type Env struct {
	bc  *Blockchain
	hub *Hub
}

func (e *Env) GetBlocks(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(e.bc.chain)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}
