package main

import (
	"go-chat/pkg"
	"go-chat/web"
	"log"
	"net/http"
)

func main() {

	_, _, err := pkg.InitNats(true, true)
	if err != nil {
		return
	}

	h := http.NewServeMux()
	// The main entrypoint to the app. Content will be lazy loaded on document render.
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := web.Index().Render(r.Context(), w)
		if err != nil {
			log.Panic(err)
		}
	})
	h.HandleFunc("/rooms", pkg.Rooms)
	h.HandleFunc("/chat", pkg.Chat)

	err = http.ListenAndServe(":8080", h)
	if err != nil {
		log.Panic(err)
	}
}
