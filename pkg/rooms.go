package pkg

import (
	"go-chat/web"
	"log"
	"net/http"
)

var rooms = []string{"Daily", "Go", "Nats", "Templ", "Datastar", "Basecoat"}

func Rooms(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "text/html")
	err := web.Rooms(rooms).Render(r.Context(), w)
	if err != nil {
		log.Panic(err)
	}
}
