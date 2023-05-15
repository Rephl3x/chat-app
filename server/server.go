package server

import (
	"log"
	"net/http"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var clients = make(map[websocket.Conn]bool) // connected clients
var broadcast = make(chan Message)          // broadcast channel

// Message struct
type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

func HandleConnections(w http.ResponseWriter, r *http.Request) {
	c, err := websocket.Accept(w, r, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
	})
	if err != nil {
		log.Println(err)
		return
	}
	defer c.Close(websocket.StatusInternalError, "Internal Server Error")

	clients[c] = true

	for {
		var msg Message
		err := wsjson.Read(r.Context(), c, &msg)
		if err != nil {
			log.Println(err)
			delete(clients, c)
			return
		}

		broadcast <- msg
	}
}

func HandleMessages() {
	for {
		msg := <-broadcast
		for client := range clients {
			err := wsjson.Write(client, msg)
			if err != nil {
				log.Println(err)
				client.Close(websocket.StatusInternalError, "Internal Server Error")
				delete(clients, client)
			}
		}
	}
}
