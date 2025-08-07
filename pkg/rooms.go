package pkg

import (
	"encoding/json"
	"fmt"
	"github.com/starfederation/datastar-go/datastar"
	"go-chat/pkg/models"
	"go-chat/web"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

var rooms = []string{"Daily", "Go", "Nats", "Templ", "Datastar", "Basecoat"}

func Rooms(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	err := web.Rooms(rooms).Render(r.Context(), w)
	if err != nil {
		log.Panic(err)
	}
}

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

var msgChan = make(chan models.Message, 100)

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

func NewMessage(w http.ResponseWriter, r *http.Request) {
	room := r.PathValue("room")
	message := r.FormValue("message")

	fmt.Println(room, message)
	msg := models.Message{
		Room:    room,
		Message: message,
		User:    "",
		Sent:    time.Now(),
	}
	messages = append(messages, msg)

	// Place the message on the global channel
	msgChan <- msg

}

func SSE(w http.ResponseWriter, r *http.Request) {
	sse := datastar.NewSSE(w, r)
	ctx := r.Context()
	for ctx.Err() == nil {
		select {
		case <-ctx.Done():
			fmt.Println(ctx.Err())
			break
		case msg := <-msgChan:
			fmt.Println("Listened and got da messag", msg)
			id := rand.Uint32()
			html := fmt.Sprintf(`<div id="%d" class="bg-white p-3 rounded shadow">
											<p class="text-sm text-gray-500">%s</p>
											<p>%s</p>
										</div>`, id, msg.User, msg.Message)
			sse.PatchElements(html, datastar.WithSelector(fmt.Sprintf("#%s_messages", msg.Room)), datastar.WithModeAppend())
			err := sse.ExecuteScript(fmt.Sprintf("document.getElementById(\"%d\").scrollIntoView()", id))
			if err != nil {
				log.Println(err)
			}
		}
	}
	fmt.Println("Getting da heck outta here")
}
