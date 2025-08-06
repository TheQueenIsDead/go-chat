package pkg

import (
	"go-chat/web"
	"log"
	"net/http"
)

func Rooms(w http.ResponseWriter, r *http.Request) {
	err := web.Rooms().Render(r.Context(), w)
	if err != nil {
		log.Panic(err)
	}
}
