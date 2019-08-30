package webserver

import (
	"fmt"
	"net/http"
)

func SlashCommand(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	default:
		fmt.Fprintf(w, "Sorry, only POST method is supported.")
	}
}
