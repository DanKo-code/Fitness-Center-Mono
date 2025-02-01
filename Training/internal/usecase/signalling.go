package internal

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/http"
	"time"
)

// AllRooms is the global hashmap for the server
var AllRooms RoomMap

// CreateRoomRequestHandler Create a Room and return roomID
func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	roomID := AllRooms.CreateRoom()

	type resp struct {
		RoomID string `json:"room_id"`
	}

	log.Println("CreateRoomRequestHandler: ", AllRooms.Map)
	json.NewEncoder(w).Encode(resp{RoomID: roomID})
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type broadcastMsg struct {
	Message map[string]interface{}
	RoomID  string
	Client  *websocket.Conn
}

var broadcast = make(chan broadcastMsg)

func broadcaster() {
	for {
		log.Println("start get from chanal broadcast: ")
		msg := <-broadcast
		log.Println("end get from chanal broadcast: ", msg)

		for _, client := range AllRooms.Map[msg.RoomID] {
			if client.Conn != msg.Client {
				client.Mutex.Lock()
				log.Println("start client.Conn.WriteJSON(msg.Message): ", msg.Message)
				err := client.Conn.WriteJSON(msg.Message)
				log.Println("end client.Conn.WriteJSON(msg.Message): ", msg.Message)
				client.Mutex.Unlock()

				if err != nil {
					log.Fatal(err)
					client.Conn.Close()
				}
			}
		}
	}
}

// JoinRoomRequestHandler will join the client in a particular room
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomID, ok := r.URL.Query()["roomID"]

	if !ok {
		log.Println("roomID missing in URL Parameters")
		return
	}

	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Web Socket Upgrade Error", err)
	}

	AllRooms.InsertIntoRoom(roomID[0], false, ws)

	go broadcaster()

	rand.Seed(time.Now().UnixNano())
	randFuncNumber := rand.Intn(100)
	for {

		log.Println("JoinRoomRequestHandler: %d", randFuncNumber)

		var msg broadcastMsg

		log.Println("Start ws.ReadJSON: ", msg)
		err := ws.ReadJSON(&msg.Message)
		if err != nil {
			log.Fatal("Read Error: ", err)
		}

		msg.Client = ws
		msg.RoomID = roomID[0]

		log.Println("End ws.ReadJSON: ", msg)

		broadcast <- msg
	}
}
