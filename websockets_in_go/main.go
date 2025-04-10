package main

import (
	"log"
	"net/http"
)

func SetupAPI() {
	manager := NewManager()
	http.Handle("/", http.FileServer(http.Dir("./frontend")))
	http.HandleFunc("/ws", manager.ServerWS)
}

func main() {
	SetupAPI()
	log.Fatal(http.ListenAndServe(":3000", nil))
}
