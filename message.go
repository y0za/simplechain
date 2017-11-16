package main

import "encoding/json"

type MessageType int

const (
	QUERY_LATEST MessageType = iota
	QUERY_ALL
	RESPONSE_BLOCKCHAIN
)

// Message is sent through websocket
type Message struct {
	Type MessageType `json:"type"`
	Data string      `json:"data"`
}

func blocksMessageJSON(blocks []Block, mt MessageType) ([]byte, error) {
	bd, err := json.Marshal(blocks)
	if err != nil {
		return nil, err
	}

	m := Message{
		Type: mt,
		Data: string(bd),
	}
	return json.Marshal(m)
}
