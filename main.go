package main

import (
	"net/http"
)

func main() {
}

func newApiServer(env Env, addr string) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("/blocks", env.GetBlocks)
	mux.HandleFunc("/mineBlock", env.MineBlock)
	mux.HandleFunc("/peers", env.GetPeers)
	mux.HandleFunc("/addPeer", env.AddPeer)

	return &http.Server{
		Handler: mux,
		Addr:    addr,
	}
}

func newP2PServer(env Env, addr string) *http.Server {
	return &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serveWs(env.hub, env.bc, w, r)
		}),
		Addr: addr,
	}
}
