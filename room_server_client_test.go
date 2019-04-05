package main

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	tmp "github.com/lijiaqigreat/wsrpcgo/protobuf"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}

func assertProtoEqual(t *testing.T, actual, expected proto.Message) {
	if !proto.Equal(actual, expected) {
		t.Errorf("actual=%v expect=%s", actual, expected)
	}
}

func TestRoomServerClientSuite(t *testing.T) {
	suite.Run(t, new(RoomServerClientSuite))
}

type RoomServerClientSuite struct {
	suite.Suite
	rs     *RoomServer
	server *httptest.Server
	dialer *websocket.Dialer
}

func (s *RoomServerClientSuite) AddAndConnectID(id string) (*websocket.Conn, *http.Response, error) {
	s.rs.AddWriter(nil, &tmp.AddWriterRequest{
		ProposedIds: []string{id},
	})
	url := makeWsProto(s.server.URL + "?id=" + id)
	return s.dialer.Dial(url, nil)
}

func (s *RoomServerClientSuite) SetupTest() {
	s.rs = NewRoomServer(nil)
	s.server = httptest.NewServer(s.rs.GetHandler())
	s.dialer = &websocket.Dialer{}
}

func (s *RoomServerClientSuite) TearDownTest() {
	s.server.Close()
	s.rs.Close()
}

func (s *RoomServerClientSuite) TestGet400WhenMissingId() {
	url := makeWsProto(s.server.URL)
	_, response, _ := s.dialer.Dial(url, nil)
	assert.Equal(s.T(), response.StatusCode, http.StatusBadRequest)
}

func (s *RoomServerClientSuite) TestGet400WhenIdNotFound() {
	url := makeWsProto(s.server.URL + "?id=not_found")
	_, response, _ := s.dialer.Dial(url, nil)
	assert.Equal(s.T(), response.StatusCode, http.StatusBadRequest)
}

func (s *RoomServerClientSuite) TestSendsIdCommandOnJoin() {
	id := "test"
	ws, _, _ := s.AddAndConnectID(id)
	_, wsMessage, _ := ws.ReadMessage()
	actual := new(tmp.Command)
	proto.Unmarshal(wsMessage, actual)

	assertProtoEqual(s.T(), actual, &tmp.Command{
		Command: &tmp.Command_IdCommand{
			IdCommand: &tmp.IdCommand{
				NewId: id,
			},
		},
	})
}

func (s *RoomServerClientSuite) TestForwardsCommandToEveryone() {
	id1 := "test1"
	id2 := "test2"
	ws1, _, _ := s.AddAndConnectID(id1)
	ws2, _, _ := s.AddAndConnectID(id2)
	message := []byte("hello")

	ws1.WriteMessage(websocket.BinaryMessage, message)

	// skip id command
	ws1.ReadMessage()
	ws2.ReadMessage()
	ws1.ReadMessage()
	ws2.ReadMessage()

	_, wsMessage, _ := ws1.ReadMessage()
	_, wsMessage2, _ := ws2.ReadMessage()

	assert.Equal(s.T(), wsMessage, wsMessage2)
	actual := new(tmp.Command)
	proto.Unmarshal(wsMessage, actual)
	assertProtoEqual(s.T(), actual, &tmp.Command{
		Command: &tmp.Command_WriterCommand{
			WriterCommand: &tmp.WriterCommand{
				Id:      id1,
				Command: message,
			},
		},
	})
}
