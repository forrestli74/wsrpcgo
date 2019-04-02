package main

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wsrpcgo/protobuf"
)

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
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

func (s *RoomServerClientSuite) AddAndConnectId(id string) (*websocket.Conn, *http.Response, error) {
	s.rs.AddConnection(id)
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

func (s *RoomServerClientSuite) TestForwardsCommandToEveryone() {
	id1 := "test1"
	id2 := "test2"
	ws1, _, _ := s.AddAndConnectId(id1)
	ws2, _, _ := s.AddAndConnectId(id2)
	message := []byte("hello")

	ws1.WriteMessage(websocket.BinaryMessage, message)
	_, wsMessage, _ := ws1.ReadMessage()
	_, wsMessage2, _ := ws2.ReadMessage()
	var actual tmp.Command
	proto.Unmarshal(wsMessage, &actual)
	expected := tmp.Command{
		Command: &tmp.Command_WriterCommand{
			WriterCommand: &tmp.WriterCommand{
				Id:      id1,
				Command: message,
			},
		},
	}

	assert.Equal(s.T(), wsMessage, wsMessage2)
	if !proto.Equal(&actual, &expected) {
		s.T().Errorf("actual=%s expect=%s", actual.String(), expected.String())
	}
}
