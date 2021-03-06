// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// https://github.com/gorilla/websocket/tree/master/examples/chat
// This file contains original code or modifications
// Originaly under BSD 2-clause license

package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sort"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Peer is a middleman between the websocket connection and the hub.
type Peer struct {
	hub *Hub

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte

	bc *Blockchain
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Peer) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, data, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("error: %v", err)
			}
			break
		}

		var message Message
		err = json.Unmarshal(data, &message)
		if err != nil {
			log.Printf("error: received invalid message %s", data)
		}
		c.handleMessage(message)
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Peer) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, []byte{}); err != nil {
				return
			}
		}
	}
}

func (c *Peer) handleMessage(message Message) {
	switch message.Type {
	case QueryLatest:
		b := c.bc.LatestBlock()
		data, err := blocksMessageJSON([]Block{*b}, QueryLatest)
		if err != nil {
			log.Printf("error: %v\n", err)
			return
		}
		c.send <- data
	case QueryAll:
		data, err := blocksMessageJSON(c.bc.chain, QueryAll)
		if err != nil {
			log.Printf("error: %v\n", err)
			return
		}
		c.send <- data
	case ResponseBlockchain:
		c.handleBlockchainResponse(message)
	}
}

func (c *Peer) handleBlockchainResponse(message Message) {
	var blocks []Block
	err := json.Unmarshal([]byte(message.Data), &blocks)
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	if len(blocks) == 0 {
		log.Println("error: received blockchain is empty")
	}

	sort.Slice(blocks, func(l, r int) bool {
		return blocks[l].Index < blocks[r].Index
	})

	lbr := blocks[len(blocks)-1]
	lbh := c.bc.LatestBlock()

	if lbr.Index <= lbh.Index {
		return
	}

	// add new hash to blockchain
	if lbr.PreviousHash == lbh.Hash {
		err = c.bc.AddBlock(lbr)
		if err != nil {
			log.Printf("error: %v\n", err)
			return
		}
		data, err := blocksMessageJSON([]Block{lbr}, ResponseBlockchain)
		if err != nil {
			log.Printf("error: %v\n", err)
			return
		}
		c.hub.broadcast <- data
		return
	}

	// new block is two or more next
	// the received node has to query all blocks from other node
	if len(blocks) == 1 {
		data, err := blocksMessageJSON([]Block{lbr}, QueryAll)
		if err != nil {
			log.Printf("error: %v\n", err)
			return
		}
		c.hub.broadcast <- data
		return
	}

	// replace blockchain data
	err = c.bc.ReplaceBlocks(blocks, GenesisBlock())
	if err != nil {
		log.Printf("error: %v\n", err)
	}
	data, err := blocksMessageJSON([]Block{lbr}, ResponseBlockchain)
	if err != nil {
		log.Printf("error: %v\n", err)
		return
	}
	c.hub.broadcast <- data
}

// serveWs handles websocket requests from the peer.
func serveWs(hub *Hub, bc *Blockchain, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	initializePeer(conn, hub, bc)
}

func connectToPeer(hub *Hub, bc *Blockchain, peerURL chan string) {
	for {
		url, ok := <-peerURL
		if !ok {
			return
		}

		conn, _, err := websocket.DefaultDialer.Dial(url, nil)
		if err != nil {
			log.Println(err)
			continue
		}

		initializePeer(conn, hub, bc)
	}
}

func initializePeer(conn *websocket.Conn, hub *Hub, bc *Blockchain) {
	peer := &Peer{
		hub:  hub,
		conn: conn,
		send: make(chan []byte, 256),
		bc:   bc,
	}
	peer.hub.register <- peer

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go peer.writePump()
	go peer.readPump()

	m := Message{
		Type: QueryLatest,
	}
	data, _ := json.Marshal(m)
	peer.send <- data
}
