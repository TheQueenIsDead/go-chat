package pkg

import (
	"go-chat/web"
	"log"
	"net/http"
)

func Rooms(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := web.Rooms().Render(r.Context(), w)
	if err != nil {
		log.Panic(err)
	}
}
