//go:build !test
// +build !test

package routes

import (
	"fmt"
	"net/http"

	"github.com/cmartinc22/devhub-backend-core/handlers"
	"github.com/pedidosya/peya-go/server"
)

func AddSchemasRoutes(s *server.Server, path string) {
	// Expose schemas
	s.AddRouteWithOptions(
		fmt.Sprintf("/%s/{apiVersion}/schemas/{schema}", path),
		handlers.HandleGetSchemas(path),
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
