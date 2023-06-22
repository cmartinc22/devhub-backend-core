package handlers

import (
	"net/http"

	"github.com/pedidosya/peya-go/logs"
	"github.com/pedidosya/peya-go/server"
)

func HandleReadGreeting(w http.ResponseWriter, r *http.Request) {
	logs.Debug("[handlers] handling read /greeting request")
	server.OK(w, r, map[string]interface{}{
		"id":      1,
		"content": "Hello, World!",
	})
}

func HandleWriteGreeting(w http.ResponseWriter, r *http.Request) {
	logs.Debug("[handlers] handling write /greeting request")

	server.OK(w, r, map[string]interface{}{})
}
