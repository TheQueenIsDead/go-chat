package pkg

import (
	"encoding/json"
	"github.com/starfederation/datastar-go/datastar"
	"go-chat/pkg/models"
	"go-chat/web"
	"log"
	"net/http"
	"strings"
	"time"
)

var messages = []models.Message{
	{"Daily", "Welcome to the Daily chatroom 👋", "TheQueenIsDead", time.Now()},
	{"Daily", "This chatroom is refreshed daily with no persistence.", "TheQueenIsDead", time.Now()},
	{"Go", "Welcome to the Go chatroom 👋", "TheQueenIsDead", time.Now()},
	{"Go", "What a delightful language!", "TheQueenIsDead", time.Now()},
	{"NATS", "Welcome to the NATS chatroom 👋", "TheQueenIsDead", time.Now()},
	{"NATS", "For all your event driven needs.", "TheQueenIsDead", time.Now()},
	{"Templ", "Welcome to the Templ chatroom 👋", "TheQueenIsDead", time.Now()},
	{"Templ", "Still not sold on this framework, but strong types are good?.", "TheQueenIsDead", time.Now()},
	{"Datastar", "Welcome to the Datastar chatroom 👋", "TheQueenIsDead", time.Now()},
	{"Datastar", "SSE driven UI is quite an experience.", "TheQueenIsDead", time.Now()},
	{"Basecoat", "Welcome to the Basecoat chatroom 👋", "TheQueenIsDead", time.Now()},
	{"Basecoat", "Please share your tips for using as little Tailwind as possible!.", "TheQueenIsDead", time.Now()},
}

func Messages(w http.ResponseWriter, r *http.Request) {

	room := r.PathValue("room")

	// Filter messages for the current room
	var roomMessages []models.Message
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
	err = sse.PatchElementTempl(web.Messages(room, roomMessages))
	if err != nil {
		log.Fatal(err)
	}
}
