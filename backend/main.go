package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	initDB()
	defer db.Close()

	router := mux.NewRouter()

	// REST
	router.HandleFunc("/messages", getSessions).Methods("GET")
	router.HandleFunc("/messages", createSession).Methods("POST")
	router.HandleFunc("/messages/{id}", deleteSession).Methods("DELETE")

	// WebSocket
	router.HandleFunc("/ws", handleConnections)
	go handleMessages()

	log.Println("Servidor rodando em :8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
