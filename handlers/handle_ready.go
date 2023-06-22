package handlers

import (
	"context"
	"net/http"

	"github.com/pedidosya/peya-go/logs"
	"github.com/pedidosya/peya-go/server"
)

func HandleReady(w http.ResponseWriter, r *http.Request, readyStatus func(context.Context) string) {
	logs.Debug("[handlers] handling ready request")

	status := readyStatus(r.Context())
	sCode := 200
	if status != "ok" {
		sCode = 503
		status = "not yet ready"
	}
	server.Render(w, r, status, sCode)
}
