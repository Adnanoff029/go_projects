package main

import (
	"context"
	"log"
	"net/http"
)

func SetupAPI() {
	ctx := context.Background()
	manager := NewManager(ctx)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.ServerWS)
	http.HandleFunc("/login", manager.LoginHandler)
}

func main() {
	SetupAPI()
	log.Fatal(http.ListenAndServe(":3000", nil))
}
