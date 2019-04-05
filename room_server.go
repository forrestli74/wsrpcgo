package main

import (
	"context"
	"math/rand"
	"net/http"
	"time"

	proto "github.com/golang/protobuf/proto"
	tmp "github.com/lijiaqigreat/wsrpcgo/protobuf"
)

/*
RoomServer ...

*/
type RoomServer struct {
	setting        *tmp.RoomSetting
	connectionByID map[string]*RoomConn
	history        History
	alive          bool
}

/*
Debug ...
*/
func (rs *RoomServer) Debug(c context.Context, request *tmp.DebugRequest) (*tmp.DebugResponse, error) {
	return nil, nil
}

/*
CreateRoom ...
*/
func (rs *RoomServer) CreateRoom(c context.Context, request *tmp.CreateRoomRequest) (*tmp.CreateRoomResponse, error) {
	return nil, nil
}

/*
AddWriter ...
*/
func (rs *RoomServer) AddWriter(c context.Context, request *tmp.AddWriterRequest) (*tmp.AddWriterResponse, error) {
	for _, id := range request.ProposedIds {
		if _, ok := rs.connectionByID[id]; !ok {
			rs.connectionByID[id] = &RoomConn{id: id}
			rs.appendRawCommand(&tmp.Command{
				Command: &tmp.Command_IdCommand{
					IdCommand: &tmp.IdCommand{
						NewId: id,
					},
				},
			})
			return &tmp.AddWriterResponse{Id: id}, nil
		}
	}
	return &tmp.AddWriterResponse{}, nil
}

/*
GetHandler ...
*/
func (rs *RoomServer) GetHandler() http.Handler {
	return roomServerHandler{rs: rs}
}

/*
Close ...
*/
func (rs *RoomServer) Close() {
	rs.alive = false
}

func (rs *RoomServer) appendRawCommand(command *tmp.Command) {
	rawCommand, _ := proto.Marshal(command)
	rs.history.AppendCommand(rawCommand)
}

/*
NewRoomServer ...
*/
func NewRoomServer(setting *tmp.RoomSetting) (rs *RoomServer) {
	rs = &RoomServer{
		connectionByID: make(map[string]*RoomConn),
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
				rs.appendRawCommand(&tmp.Command{
					Command: &tmp.Command_TickCommand{
						TickCommand: &tmp.TickCommand{
							RandomSeed: randomBuffer,
						},
					},
				})
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
	if rc, ok := rsh.rs.connectionByID[id]; ok {
		rc.Connect(rsh.rs, w, r)
	} else {
		http.Error(w, "Id not found", http.StatusBadRequest)
	}
}
