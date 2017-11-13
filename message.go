package main

type MessageType int

const (
	QUERY_LATEST MessageType = iota
	QUERY_ALL
	RESPONSE_BLOCKCHAIN
)

// Message is sent through websocket
type Message struct {
	Type MessageType `json:"type"`
	Data []byte      `json:"data"`
}
