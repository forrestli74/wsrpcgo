package main

import (
	"net/http"
)

type RoomServer struct {
	connectionById map[string]*RoomConn
	history        History
}

func (rs *RoomServer) GetHandler() http.Handler {
	return roomServerHandler{rs: rs}
}

func (rs *RoomServer) AddConnection(id string) error {
	rs.connectionById[id] = &RoomConn{id: id}
	return nil
}

func NewRoomServer() *RoomServer {
	return &RoomServer{
		connectionById: make(map[string]*RoomConn),
		history:        CreateHistory(),
	}
}

type roomServerHandler struct {
	rs *RoomServer
}

func (rsh roomServerHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Access-Control-Allow-Origin", "*")
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, "Missing id", http.StatusBadRequest)
		return
	}
	if rc, ok := rsh.rs.connectionById[id]; ok {
		rc.Connect(rsh.rs, w, r)
	} else {
		http.Error(w, "Id not found", http.StatusBadRequest)
	}
}
