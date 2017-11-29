// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// https://github.com/gorilla/websocket/tree/master/examples/chat
// This file contains original code or modifications
// Originaly under BSD 2-clause license

package main

// Hub maintains the set of active peers and broadcasts messages to the
// peers.
type Hub struct {
	// Registered peers.
	peers map[*Peer]bool

	// Inbound messages from the peers.
	broadcast chan []byte

	// Register requests from the peers.
	register chan *Peer

	// Unregister requests from peers.
	unregister chan *Peer
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Peer),
		unregister: make(chan *Peer),
		peers:      make(map[*Peer]bool),
	}
}

func (h *Hub) run() {
	for {
		select {
		case peer := <-h.register:
			h.peers[peer] = true
		case peer := <-h.unregister:
			if _, ok := h.peers[peer]; ok {
				delete(h.peers, peer)
				close(peer.send)
			}
		case message := <-h.broadcast:
			for peer := range h.peers {
				select {
				case peer.send <- message:
				default:
					close(peer.send)
					delete(h.peers, peer)
				}
			}
		}
	}
}
