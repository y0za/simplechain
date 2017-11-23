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

// peerJSON used in AddPeer for reading request body
type peerJSON struct {
	Peer string `json:"peer"`
}

func (e *Env) GetBlocks(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(e.bc.chain)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func (e *Env) AddPeer(w http.ResponseWriter, r *http.Request) {
	var pj peerJSON
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&pj)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	connectToPeer(e.hub, e.bc, pj.Peer)
}
