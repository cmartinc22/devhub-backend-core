//go:build !test
// +build !test

package routes

import (
	"net/http"

	"github.com/cmartinc22/devhub-backend-core/handlers"
	"github.com/cmartinc22/devhub-backend-core/models"
	"github.com/pedidosya/peya-go/server"
)

func AddHealthRoutes(s *server.Server, info models.ServiceInfo) {
	// Expose health route
	s.AddRouteWithOptions(
		"/",
		handlers.HandleAlive(info),
		&server.RouteOptions{
			CORSEnabled: true,
			CORSOptions: []server.CORSOption{
				server.AllowedMethods([]string{http.MethodGet}),
			},
			TimeOutSeconds: 0,
		},
		http.MethodGet)

	s.AddRouteWithOptions(
		"/alive",
		handlers.HandleAlive(info),
		&server.RouteOptions{
			CORSEnabled: true,
			CORSOptions: []server.CORSOption{
				server.AllowedMethods([]string{http.MethodGet}),
			},
			TimeOutSeconds: 0,
		},
		http.MethodGet)
}
