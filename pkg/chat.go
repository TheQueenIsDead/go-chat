package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"github.com/starfederation/datastar-go/datastar"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type Message struct {
	Id      uint32
	Room    string
	Message string
	User    string
	Sent    time.Time
}

var (
	rooms    = []string{"Daily", "Go", "Nats", "Templ", "Datastar", "Basecoat"}
	messages = []Message{
		{1, "Daily", "Welcome to the Daily chatroom 👋", "TheQueenIsDead", time.Now()},
		{1, "Daily", "This chatroom is refreshed daily with no persistence.", "TheQueenIsDead", time.Now()},
		{1, "Go", "Welcome to the Go chatroom 👋", "TheQueenIsDead", time.Now()},
		{1, "Go", "What a delightful language!", "TheQueenIsDead", time.Now()},
		{1, "NATS", "Welcome to the NATS chatroom 👋", "TheQueenIsDead", time.Now()},
		{1, "NATS", "For all your event driven needs.", "TheQueenIsDead", time.Now()},
		{1, "Templ", "Welcome to the Templ chatroom 👋", "TheQueenIsDead", time.Now()},
		{1, "Templ", "Still not sold on this framework, but strong types are good?.", "TheQueenIsDead", time.Now()},
		{1, "Datastar", "Welcome to the Datastar chatroom 👋", "TheQueenIsDead", time.Now()},
		{1, "Datastar", "SSE driven UI is quite an experience.", "TheQueenIsDead", time.Now()},
		{1, "Basecoat", "Welcome to the Basecoat chatroom 👋", "TheQueenIsDead", time.Now()},
		{1, "Basecoat", "Please share your tips for using as little Tailwind as possible!.", "TheQueenIsDead", time.Now()},
	}
)

func (s *Server) Index(w http.ResponseWriter, r *http.Request) {
	err := Index().Render(r.Context(), w)
	if err != nil {
		s.log.Error(err)
	}
}

func (s *Server) GetRooms(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	err := sse.PatchElementTempl(RoomsComponent(rooms))
	if err != nil {
		s.log.Error(err)
	}
}

func (s *Server) GetMessages(w http.ResponseWriter, r *http.Request) {

	room := r.PathValue("room")

	// Filter messages for the current room
	var roomMessages []Message
	for _, m := range messages {
		if strings.ToLower(m.Room) == strings.ToLower(room) {
			roomMessages = append(roomMessages, m)
		}
	}

	sse := datastar.NewSSE(w, r)

	// Patch the signal for the active room
	signal, err := json.Marshal(map[string]string{"active": room})
	if err != nil {
		s.log.Error(err)
		return
	}
	err = sse.PatchSignals(signal)
	if err != nil {
		s.log.Error(err)
		return
	}

	// Update the messages screen with the messages above
	err = sse.PatchElementTempl(MessagesComponent(room, roomMessages))
	if err != nil {
		s.log.Error(err)
		return
	}

	// Signal the client to scroll to the newest message
	err = sse.ExecuteScript(fmt.Sprintf("document.getElementById(\"%d\").scrollIntoView()", roomMessages[len(roomMessages)-1].Id))
	if err != nil {
		s.log.Error(err)
		return
	}
}

func (s *Server) PutMessage(w http.ResponseWriter, r *http.Request) {
	room := r.PathValue("room")
	message := r.FormValue("message")

	sse := datastar.NewSSE(w, r)
	signal, err := json.Marshal(map[string]string{"messageinput": ""})
	if err != nil {
		s.log.Error(err)
		return
	}
	err = sse.PatchSignals(signal)
	if err != nil {
		s.log.Error(err)
		return
	}

	msg := Message{
		Id:      rand.Uint32(),
		Room:    room,
		Message: message,
		User:    r.Context().Value("user").(string),
		Sent:    time.Now(),
	}
	messages = append(messages, msg)

	// Publish the message to NATS
	data, err := json.Marshal(msg)
	if err != nil {
		s.log.Error(err)
		return
	}
	_, err = s.js.PublishAsync("messages", data)
	if err != nil {
		s.log.Error(err)
	}
	s.log.WithFields(logrus.Fields{"message": msg.Message, "user": msg.User}).Debug("published to nats")
}

func (s *Server) SSE(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	s.log.WithField("user", r.Context().Value("user")).Debug("client subscribing")
	sub, err := s.js.Subscribe("messages", func(msg *nats.Msg) {
		var message Message
		err := json.Unmarshal(msg.Data, &message)
		if err != nil {
			s.log.Error(err)
			return
		}
		err = sse.PatchElementTempl(MessageComponent(message), datastar.WithSelector(fmt.Sprintf("#%s_messages", message.Room)), datastar.WithModeAppend())
		if err != nil {
			s.log.Error(err)
			return
		}
		err = sse.ExecuteScript(fmt.Sprintf("document.getElementById(\"%d\").scrollIntoView()", message.Id))
		if err != nil {
			s.log.Error(err)
			return
		}
		return
	})
	if err != nil {
		s.log.Error(err)
		return
	}

	select {
	case <-r.Context().Done():
		logger := s.log.WithFields(logrus.Fields{
			"user": r.Context().Value("user"),
			"err":  r.Context().Err(),
		})
		logger.Debug("context done")
		if err := sub.Unsubscribe(); err != nil {
			logger.Error(err)
		}
		break
	}
}
