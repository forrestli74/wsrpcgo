package main

import (
	proto "github.com/golang/protobuf/proto"
	//"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
	"wsrpcgo/protobuf"
)

func TestRoomServerSuite(t *testing.T) {
	suite.Run(t, new(RoomServerSuite))
}

type RoomServerSuite struct {
	suite.Suite
}

func (s *RoomServerClientSuite) TestSendsTick() {
	size := 2
	setting := &tmp.RoomSetting{
		TickSetting: &tmp.TickSetting{
			Size:            uint32(size),
			FrequencyMillis: 10,
		},
	}
	rs := NewRoomServer(setting)
	defer rs.Close()
	ch := rs.history.CreateChan(0)
	rawCommand := <-ch
	var actual tmp.Command
	err := proto.Unmarshal(rawCommand, &actual)
	assert.Nil(s.T(), err)
	seed := actual.GetTickCommand().GetRandomSeed()
	assert.Equal(s.T(), len(seed), size)
}
