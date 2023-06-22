//go:build !test
// +build !test

package main

import (
	"net/http"

	"github.com/cmartinc22/devhub-backend-core/handlers"
	"github.com/pedidosya/peya-go/server"
)

func addASchemasRoutes(s *server.Server) {
	// Expose schemas
	s.AddRouteWithOptions(
		"/api/{apiVersion}/schemas/{schema}",
		handlers.HandleGetSchemas("api"),
		&server.RouteOptions{
			CORSEnabled: true,
			CORSOptions: []server.CORSOption{
				server.AllowedMethods([]string{http.MethodGet, http.MethodOptions}),
			},
			TimeOutSeconds: 0,
		},
		http.MethodGet,
		http.MethodOptions,
	)
}
