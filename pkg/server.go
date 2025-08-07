package pkg

import (
	"context"
	"fmt"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/nats-io/nats.go"
	"log"
	"net/http"
)

type Server struct {
	http *http.Server
	js   nats.JetStreamContext
}

func NewServer() *Server {

	s := new(Server)

	// Configure NATS
	_, js, _, err := InitNats(true, true)
	if err != nil {
		fmt.Println(err)
		return s
	}
	s.js = js

	// Configure HTTP
	h := http.NewServeMux()
	// The main entrypoint to the app. Content will be lazy loaded on document render.
	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		err := Index().Render(r.Context(), w)
		if err != nil {
			log.Panic(err)
		}
	})

	h.HandleFunc("/rooms", s.GetRooms)
	h.HandleFunc("/rooms/{room}/messages", s.GetMessages)
	h.HandleFunc("POST /rooms/{room}/messages", s.PutMessage)
	h.HandleFunc("/sse", s.SSE)
	s.http = &http.Server{
		Handler: userMiddleware(h),
		Addr:    ":8080",
	}

	return s
}
func (s *Server) Run() error {
	return s.http.ListenAndServe()
}

// userMiddleware attempts to retrieve the randomly generated username from a cookie provided by the client,
// and if not present, will generate a new username and set a cookie on response.
// The username is placed into the request context for access in the handlers.
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
