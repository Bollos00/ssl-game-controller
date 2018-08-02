package controller

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

// handle incoming web socket connections
func WsHandler(w http.ResponseWriter, r *http.Request) {
	u := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin:     func(*http.Request) bool { return true },
	}

	conn, err := u.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()
	defer log.Println("Client disconnected")

	log.Println("Client connected")

	go listenForNewEvents(conn)

	publishState(conn)
}

func publishState(conn *websocket.Conn) {
	for {
		b, err := json.Marshal(refBox.State)
		if err != nil {
			log.Println("Marshal error:", err)
		}

		err = conn.WriteMessage(websocket.TextMessage, b)
		if err != nil {
			log.Println("Could not write message.", err)
			return
		}

		// wait to be notified
		<-refBox.notifyUpdateState
	}
}

func listenForNewEvents(conn *websocket.Conn) {
	for {
		messageType, b, err := conn.ReadMessage()
		if err != nil || messageType != websocket.TextMessage {
			log.Println("Could not read message: ", err)
			return
		}

		handleNewEventMessage(b)
	}
}

func handleNewEventMessage(b []byte) {
	event := RefBoxEvent{}
	err := json.Unmarshal(b, &event)
	if err != nil {
		log.Println("Could not read event:", string(b), err)
		return
	}

	err = processEvent(&event)
	if err != nil {
		log.Println("Could not process event:", string(b), err)
		return
	}

	refBox.SaveState()
	refBox.Update(event.Command)
}
