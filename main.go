package main

import (
	"context"
	"github.com/moby/moby/pkg/namesgenerator"
	"go-chat/pkg"
	"go-chat/web"
	"log"
	"net/http"
)

// User middleware
func userMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("user")
		if err != nil {
			username := namesgenerator.GetRandomName(0)
			cookie = &http.Cookie{Name: "user", Value: username}
			http.SetCookie(w, cookie)
		}
		ctx := context.WithValue(r.Context(), "user", cookie.Value)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func main() {

	_, _, err := pkg.InitNats(true, false)
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
	h.HandleFunc("/rooms/{room}/messages", pkg.Messages)
	h.HandleFunc("POST /rooms/{room}/messages", pkg.NewMessage)
	h.HandleFunc("/sse", pkg.SSE)

	err = http.ListenAndServe(":8080", userMiddleware(h))
	if err != nil {
		log.Panic(err)
	}
}
