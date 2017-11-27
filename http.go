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

// mineJSON used in MineBlock for reading request body
type mineJSON struct {
	Data string `json:"data"`
}

func (e *Env) GetBlocks(w http.ResponseWriter, r *http.Request) {
	err := json.NewEncoder(w).Encode(e.bc.chain)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
	}
}

func (e *Env) MineBlock(w http.ResponseWriter, r *http.Request) {
	var mj mineJSON
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&mj)
	if err != nil {
		http.Error(w, http.StatusText(400), 400)
		return
	}

	b := NextBlock(*e.bc.LatestBlock(), mj.Data)
	err = e.bc.AddBlock(b)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}

	data, err := blocksMessageJSON([]Block{b}, QueryLatest)
	if err != nil {
		http.Error(w, http.StatusText(500), 500)
		return
	}
	e.hub.broadcast <- data
}

func (e *Env) GetPeers(w http.ResponseWriter, r *http.Request) {
	peers := make([]string, 0, len(e.hub.clients))
	for c, _ := range e.hub.clients {
		peers = append(peers, c.conn.RemoteAddr().String())
	}

	err := json.NewEncoder(w).Encode(peers)
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
