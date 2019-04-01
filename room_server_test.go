package main

import (
	proto "github.com/golang/protobuf/proto"
	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wsrpcgo/protobuf"
)

type context struct {
	rs     *RoomServer
	server *httptest.Server
	t      *testing.T
	dialer *websocket.Dialer
}

func makeWsProto(s string) string {
	return "ws" + strings.TrimPrefix(s, "http")
}

func readerToString(r io.Reader) string {
	bytes, _ := ioutil.ReadAll(r)
	return string(bytes)
}

type roomServerSuite struct {
	suite.Suite
	rs     *RoomServer
	server *httptest.Server
	dialer *websocket.Dialer
}

func (s *roomServerSuite) SetupTest() {
	s.rs = NewRoomServer()
	s.server = httptest.NewServer(s.rs.GetHandler())
	s.dialer = &websocket.Dialer{}
}

func TestRoomServerSuite(t *testing.T) {
	suite.Run(t, new(roomServerSuite))
}

func (s *roomServerSuite) TestGet400WhenMissingId() {
	url := makeWsProto(s.server.URL)
	_, response, _ := s.dialer.Dial(url, nil)
	assert.Equal(s.T(), response.StatusCode, http.StatusBadRequest)
}

func (s *roomServerSuite) TestGet400WhenIdNotFound() {
	url := makeWsProto(s.server.URL + "?id=not_found")
	_, response, _ := s.dialer.Dial(url, nil)
	assert.Equal(s.T(), response.StatusCode, http.StatusBadRequest)
}

func (s *roomServerSuite) TestForwardsCommandToEveryone() {
	id := "test"
	s.rs.AddConnection(id)
	url := makeWsProto(s.server.URL + "?id=" + id)
	ws, _, _ := s.dialer.Dial(url, nil)
	message := []byte("hello")

	ws.WriteMessage(websocket.BinaryMessage, message)
	_, wsMessage, _ := ws.ReadMessage()
	var actual tmp.Command
	proto.Unmarshal(wsMessage, &actual)
	expected := tmp.Command{
		Command: &tmp.Command_WriterCommand{
			WriterCommand: &tmp.WriterCommand{
				Id:      id,
				Command: message,
			},
		},
	}

	if !proto.Equal(&actual, &expected) {
		s.T().Errorf("actual=%s expect=%s", actual.String(), expected.String())
	}
}
