package main

import (
	"log"
	"net/http"

	"./server"
)

func main() {
	http.HandleFunc("/ws", server.HandleConnections)
	go server.HandleMessages()

	// Start the server on localhost port 8080 and log any errors
	log.Println("http server started on :8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
