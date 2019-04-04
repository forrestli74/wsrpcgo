package main

import (
	"testing"

	proto "github.com/golang/protobuf/proto"
	tmp "github.com/lijiaqigreat/wsrpcgo/protobuf"
)

// func TestRoomServer_ConnectMs_Connects(t *testing.T) {
// 	rc := NewRoomServer()
// 	rc.SetConnection(RoomConnection{
// 		id: "id",
// 		secret: "secret",
// 	})
// 	err := rc.ConnectWs("id", "secret", nil)
// 	if(err != nil) {
// 		t.Errorf(err.Error())
// 	}
// }

func TestBasicProto(t *testing.T) {
	command := &tmp.Command{}
	data, _ := proto.Marshal(command)
	command2 := &tmp.Command{}
	proto.Unmarshal(data, command2)
	if !proto.Equal(command, command2) {
		t.Errorf("fail")
	}
}

func TestBasicProtoOneof(t *testing.T) {
	debugRequest := &tmp.DebugRequest{
		B:   []byte{1, 2, 3},
		S:   "test",
		I32: 32,
		I64: 64,
		Oo: &tmp.DebugRequest_Oi32{
			32,
		},
	}
	data, _ := proto.Marshal(debugRequest)
	debugRequest2 := &tmp.DebugRequest{}
	proto.Unmarshal(data, debugRequest2)
	if !proto.Equal(debugRequest, debugRequest2) {
		t.Errorf("fail")
	}
}

func TestBasicChan(t *testing.T) {
	ch := make(chan string)
	go func() { ch <- "c1" }()
	str := <-ch
	if str != "c1" {
		t.Errorf("fail")
	}
}

func TestDebug(t *testing.T) {
	// var impl tmp.RoomServiceServer
	// impl = RoomServiceImpl{}
	// impl.Debug(nil, nil)

	// var impl2 interface{}
	// impl2 = impl1

	// impl3 := reflection.ValueOf(impl2)

}
