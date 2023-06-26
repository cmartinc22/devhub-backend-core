package handlers

import (
	"net/http"

	"github.com/cmartinc22/devhub-backend-core/models"
	"github.com/pedidosya/peya-go/server"
)

func HandleAlive(info models.ServiceInfo) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		server.OK(w, r, info)
	}
}
