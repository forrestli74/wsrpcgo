package main

import (
	proto "github.com/golang/protobuf/proto"
	"math/rand"
	"net/http"
	"time"
	"wsrpcgo/protobuf"
)

type RoomServer struct {
	setting        *tmp.RoomSetting
	connectionById map[string]*RoomConn
	history        History
	alive          bool
}

func (rs *RoomServer) GetHandler() http.Handler {
	return roomServerHandler{rs: rs}
}

func (rs *RoomServer) AddConnection(id string) error {
	rs.connectionById[id] = &RoomConn{id: id}
	return nil
}

func (rs *RoomServer) Close() {
	rs.alive = false
}

func NewRoomServer(setting *tmp.RoomSetting) (rs *RoomServer) {
	rs = &RoomServer{
		connectionById: make(map[string]*RoomConn),
		history:        CreateHistory(),
		setting:        setting,
		alive:          true,
	}
	period := setting.GetTickSetting().GetFrequencyMillis()
	if period != 0 {
		ticker := time.NewTicker(time.Duration(period) * time.Millisecond)
		go func() {
			randomBuffer := make([]byte, setting.GetTickSetting().GetSize())
			for range ticker.C {
				if !rs.alive {
					break
				}
				rand.Read(randomBuffer)
				tick := tmp.Command{
					Command: &tmp.Command_TickCommand{
						TickCommand: &tmp.TickCommand{
							RandomSeed: randomBuffer,
						},
					},
				}
				rawCommand, _ := proto.Marshal(&tick)
				rs.history.AppendCommand(rawCommand)
			}
		}()
	}
	return
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
