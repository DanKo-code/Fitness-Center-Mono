package model

import (
	"Training/pkg/logger"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

type Participant struct {
	Conn  *websocket.Conn
	Mutex *sync.Mutex
}

type RoomMapKey struct {
	RoomId   uuid.UUID
	ClientId uuid.UUID
	CoachId  uuid.UUID
}

type RoomMap struct {
	Mutex sync.RWMutex
	Map   map[RoomMapKey][]Participant
}

func (r *RoomMap) Init() {
	r.Map = make(map[RoomMapKey][]Participant)
}

func (r *RoomMap) Get(roomKey RoomMapKey) []Participant {
	r.Mutex.RLock()
	defer r.Mutex.RUnlock()

	return r.Map[roomKey]
}

func (r *RoomMap) InitRoom(roomKey RoomMapKey) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	logger.Logger.Info("Init room: %v", roomKey)
	r.Map[roomKey] = []Participant{}
}

func (r *RoomMap) InsertIntoRoom(roomKey RoomMapKey, conn *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	p := Participant{conn, &sync.Mutex{}}

	log.Println("Inserting into Room with RoomID: ", roomKey.RoomId)
	r.Map[roomKey] = append(r.Map[roomKey], p)
}

func (r *RoomMap) DeleteRoom(roomKey RoomMapKey) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	delete(r.Map, roomKey)
}
