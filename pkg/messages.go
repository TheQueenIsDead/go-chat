package pkg

import (
	"fmt"
	"go-chat/pkg/models"
	"go-chat/web"
	"log"
	"net/http"
	"slices"
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
	fmt.Println("Checking for", room)

	m := slices.DeleteFunc(messages, func(m models.Message) bool {
		return strings.ToLower(m.Room) != strings.ToLower(room)
	})

	w.Header().Set("Content-Type", "text/html")
	err := web.Messages(m).Render(r.Context(), w)
	if err != nil {
		log.Panic(err)
	}
}
