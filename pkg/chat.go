package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/starfederation/datastar-go/datastar"
	"log"
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

func (s *Server) GetRooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := RoomsComponent(rooms).Render(r.Context(), w)
	if err != nil {
		log.Panic(err)
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
	signal, _ := json.Marshal(map[string]string{"active": room})
	err := sse.PatchSignals(signal)
	if err != nil {
		log.Println(err)
	}

	// Update the messages screen with the messages above
	err = sse.PatchElementTempl(MessagesComponent(room, roomMessages))
	if err != nil {
		log.Fatal(err)
	}

	// Signal the client to scroll to the newest message
	sse.ExecuteScript(fmt.Sprintf("document.getElementById(\"%d\").scrollIntoView()", roomMessages[len(roomMessages)-1].Id))
}

func (s *Server) PutMessage(w http.ResponseWriter, r *http.Request) {
	room := r.PathValue("room")
	message := r.FormValue("message")

	sse := datastar.NewSSE(w, r)
	signal, _ := json.Marshal(map[string]string{"messageinput": ""})
	sse.PatchSignals(signal)

	fmt.Println(room, message)
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
		fmt.Println(err)
	}
	_, err = s.js.PublishAsync("messages", data)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Published to NATS!!!", message)

}

func (s *Server) SSE(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	ctx := r.Context()
	fmt.Println("Subscwibin")
	_, _ = s.js.Subscribe("messages", func(msg *nats.Msg) {
		var message Message
		err := json.Unmarshal(msg.Data, &message)
		if err != nil {
			fmt.Println("Failed to unmarshal")
		}
		fmt.Println("Listened and got da NATS messag", msg)
		sse.PatchElementTempl(MessageComponent(message), datastar.WithSelector(fmt.Sprintf("#%s_messages", message.Room)), datastar.WithModeAppend())
		err = sse.ExecuteScript(fmt.Sprintf("document.getElementById(\"%d\").scrollIntoView()", message.Id))
		if err != nil {
			log.Println(err)
		}
		return
	})

	select {
	case <-ctx.Done():
		fmt.Println(ctx.Err())
		break
	}

	fmt.Println("Getting da heck outta here")
}
