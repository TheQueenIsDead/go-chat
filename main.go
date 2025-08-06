package main

import (
	"go-chat/pkg"
	"go-chat/web"
	"log"
	"net/http"
)

func main() {
  //TIP <p>Press <shortcut actionId="ShowIntentionActions"/> when your caret is at the underlined text
  // to see how GoLand suggests fixing the warning.</p><p>Alternatively, if available, click the lightbulb to view possible fixes.</p>
  s := "gopher"
  fmt.Printf("Hello and welcome, %s!\n", s)

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
