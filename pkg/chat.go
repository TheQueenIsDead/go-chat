package pkg

import (
	"go-chat/web"
	"log"
	"net/http"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	err := web.Chat().Render(r.Context(), w)
	if err != nil {
		log.Panic(err)
	}
}
