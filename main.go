package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	bc := NewBlockchain(GenesisBlock())
	hub := newHub()
	go hub.run()
	peerURL := make(chan string)
	env := &Env{bc, hub, peerURL}

	httpPort := getEnv("HTTP_PORT", "3001")
	p2pPort := getEnv("P2P_PORT", "6001")
	peersRaw := getEnv("PEERS", "")
	peers := strings.Split(peersRaw, ",")

	api := newApiServer(env, ":"+httpPort)
	p2p := newApiServer(env, ":"+p2pPort)

	go func() {
		if err := api.ListenAndServe(); err != http.ErrServerClosed {
			log.Print(err)
		}
	}()
	go func() {
		if err := p2p.ListenAndServe(); err != http.ErrServerClosed {
			log.Print(err)
		}
	}()
	go connectToPeer(hub, bc, peerURL)

	for _, p := range peers {
		if p != "" {
			peerURL <- p
		}
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop

	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	api.Shutdown(ctx)
	p2p.Shutdown(ctx)
	close(peerURL)
}

func newApiServer(env *Env, addr string) *http.Server {
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

func newP2PServer(env *Env, addr string) *http.Server {
	return &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			serveWs(env.hub, env.bc, w, r)
		}),
		Addr: addr,
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
