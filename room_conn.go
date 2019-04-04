package main

import (
	"fmt"
	"log"
	"net/http"

	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	tmp "github.com/lijiaqigreat/wsrpcgo/protobuf"
)

type RoomConnState struct {
	ch <-chan []byte
	ws *websocket.Conn
}

type RoomConn struct {
	id    string
	state RoomConnState
}

func (rc *RoomConn) Connect(rs *RoomServer, w http.ResponseWriter, r *http.Request) {
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}

	commandChan := rs.history.CreateChan(0)

	cleanUp := func() {
		fmt.Println("cleanUp")
		ws.Close()
		rc.state.ws = nil
	}

	// write ws to history
	go func() {
		defer cleanUp()
		for {
			_, message, err := ws.ReadMessage()
			if err != nil {
				fmt.Println(err)
				break
			}
			commandRaw, err := proto.Marshal(&tmp.Command{
				Command: &tmp.Command_WriterCommand{
					WriterCommand: &tmp.WriterCommand{
						Id:      rc.id,
						Command: message,
					},
				},
			})
			if err != nil {
				fmt.Println(err)
			} else {
				rs.history.AppendCommand(commandRaw)
			}
		}
	}()

	// write history to ws
	go func() {
		defer cleanUp()
		for command := range commandChan {
			err := ws.WriteMessage(websocket.BinaryMessage, command)
			if err != nil {
				break
			}
		}
	}()
}
