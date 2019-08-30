package webserver

import "net/http"

func Server() {
	http.HandleFunc("/slack/isp", SlashCommand)

	http.ListenAndServe(":8080", nil)
}
