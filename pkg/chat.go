package pkg

import (
	"go-chat/web"
	"log"
	"net/http"
)

func Chat(w http.ResponseWriter, r *http.Request) {
	err := web.Chat().Render(r.Context(), w)
	if err != nil {
		log.Panic(err)
	}
}
