package pkg

import (
	"context"
	"github.com/moby/moby/pkg/namesgenerator"
	"github.com/nats-io/nats.go"
	"github.com/sirupsen/logrus"
	"net/http"
)

type Server struct {
	http *http.Server
	js   nats.JetStreamContext
	log  *logrus.Logger
}

func NewServer() *Server {

	s := new(Server)

	// Logging
	s.log = logrus.New()
	s.log.SetLevel(logrus.DebugLevel)

	// Configure NATS
	_, js, _, err := InitNats(true, false)
	if err != nil {
		s.log.Panic(err)
	}
	s.js = js

	// Configure HTTP
	h := http.NewServeMux()
	h.HandleFunc("/", s.Index)
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
	s.log.Infof("Listening on port %s", s.http.Addr)
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
